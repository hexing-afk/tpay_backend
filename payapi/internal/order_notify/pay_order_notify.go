package order_notify

import (
	"context"
	"encoding/json"
	"time"
	"tpay_backend/model"
	"tpay_backend/payapi/internal/svc"
	"tpay_backend/utils"

	"github.com/tal-tech/go-zero/core/logx"
)

type PayOrderNotify struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewPayOrderNotify(ctx context.Context, svcCtx *svc.ServiceContext) *PayOrderNotify {
	return &PayOrderNotify{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *PayOrderNotify) OrderNotify(orderNo string) {
	// 1.查询订单
	order, err := model.NewPayOrderModel(l.svcCtx.DbEngine).FindOneByOrderNo(orderNo)
	if err != nil {
		l.Errorf("查询代收订单[%v]失败, err=%v", orderNo, err)
		return
	}

	l.Infof("代收订单信息：%+v", order)

	//通知状态-超时
	if order.NotifyStatus == model.PayNotifyStatusTimeOut {
		l.Errorf("代收订单[%v]已经超出了通知时间", orderNo)
		return
	}

	//通知状态-已通知
	if order.NotifyStatus == model.PayNotifyStatusSuccess {
		l.Errorf("代收订单[%v]已成功通知，不需要再通知", orderNo)
		return
	}

	// 没有通知地址就不再通知
	if order.NotifyUrl == "" {
		l.Infof("代收订单[%v]缺少回调地址，不再进行通知", orderNo)
		data := &model.OrderNotifyLog{
			OrderNo:         order.OrderNo,
			MerchantOrderNo: order.MerchantOrderNo,
			Status:          model.OrderNotifyStatusFail,
			OrderType:       model.NotifyLogOrderTypePay,
			Result:          "订单缺少回调地址",
		}
		if err := model.NewOrderNotifyLogModel(l.svcCtx.DbEngine).Insert(data); err != nil {
			l.Errorf("代收订单异步通知日志记录失败, data：%+v, err：%v", data, err)
			return
		}

		payNotify := model.PayNotifyInfo{
			NotifyStatus:    model.PayNotifyStatusTimeOut,
			NotifyFailTimes: order.NotifyFailTimes + 1,
			NextNotifyTime:  0,
		}
		if err := model.NewPayOrderModel(l.svcCtx.DbEngine).UpdateNotify(order.Id, payNotify); err != nil {
			l.Errorf("修改订单-通知超时 失败, order.Id:%v, err:%v", order.Id, err)
		}

		return
	}

	// 2.查询商户
	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneByMerchantNo(order.MerchantNo)
	if err != nil {
		l.Errorf("查询商户信息失败, MerchantNo:%v, err:%v", order.MerchantNo, err)
		return
	}

	// 3.打包参数
	packReq := &utils.PackPayNotifyParamsRequest{
		MerchantNo:      order.MerchantNo,
		Timestamp:       time.Now().Unix(),
		NotifyType:      utils.PayNotifyType,
		OrderNo:         order.OrderNo,
		MerchantOrderNo: order.MerchantOrderNo,
		ReqAmount:       order.ReqAmount,
		Currency:        order.Currency,
		OrderStatus:     order.OrderStatus,
		PayTime:         order.UpdateTime,
		Subject:         order.Subject,
		PaymentAmount:   order.PaymentAmount,
	}
	dataByte, err := utils.PackPayNotifyParams(packReq, merchant.Md5Key)
	if err != nil {
		l.Errorf("打包参数失败, packReq: %+v, err=%v", packReq, err)
		return
	}

	l.Infof("代收订单[%v]第%v次通知, 失败次数：[%v]", orderNo, order.NotifyFailTimes+1, order.NotifyFailTimes)

	// 4.发送通知
	result, err := utils.PostForm(order.NotifyUrl, dataByte)
	if err != nil {
		l.Errorf("发送失败, utils.Post() error: %v", err)
		result = []byte(err.Error())
	}

	resultStr := string(result)
	resultMap := map[string]string{
		"url":    order.NotifyUrl,
		"result": resultStr,
	}

	resultJ, jErr := json.Marshal(resultMap)
	if jErr != nil {
		l.Errorf("通知结果转json失败, resultMap: %+v, err=%v", resultMap, jErr)
		return
	}

	// 5.根据通知结果修改订单通知结果
	switch resultStr {
	case ResultNotifySuccess:
		// 插入通知成功日志
		log := &model.OrderNotifyLog{
			OrderNo:         order.OrderNo,
			MerchantOrderNo: order.MerchantOrderNo,
			Status:          model.OrderNotifyStatusSuccess,
			OrderType:       model.NotifyLogOrderTypePay,
		}
		if err := model.NewOrderNotifyLogModel(l.svcCtx.DbEngine).Insert(log); err != nil {
			l.Errorf("代收订单异步通知日志记录失败, data：%+v, err：%v", log, err)
		}

		// 修改代收订单异步通知信息-通知成功
		payNotify := model.PayNotifyInfo{
			NotifyStatus:    model.PayNotifyStatusSuccess,
			NotifyFailTimes: order.NotifyFailTimes,
			NextNotifyTime:  0,
		}
		if err := model.NewPayOrderModel(l.svcCtx.DbEngine).UpdateNotify(order.Id, payNotify); err != nil {
			l.Errorf("修改代收订单[%v]通知状态失败, err=[%v]", orderNo, err)
		}

		l.Infof("代收订单第%v次通知结束", order.NotifyFailTimes+1)
		return
	default:
		// 插入通知失败日志
		log := &model.OrderNotifyLog{
			OrderNo:         order.OrderNo,
			MerchantOrderNo: order.MerchantOrderNo,
			Status:          model.OrderNotifyStatusFail,
			OrderType:       model.NotifyLogOrderTypePay,
			Result:          string(resultJ),
		}
		if err := model.NewOrderNotifyLogModel(l.svcCtx.DbEngine).Insert(log); err != nil {
			l.Errorf("代收订单异步通知日志记录失败, data：%+v, err：%v", log, err)
		}

		// 下一次通知的等待时间
		waitTime := GetNotifyWaitTimeById(order.NotifyFailTimes + 1)
		if waitTime == -1 {
			// 修改代收订单异步通知信息-通知超时
			payNotify := model.PayNotifyInfo{
				NotifyStatus:    model.PayNotifyStatusTimeOut,
				NotifyFailTimes: order.NotifyFailTimes + 1,
				NextNotifyTime:  0,
			}
			if err := model.NewPayOrderModel(l.svcCtx.DbEngine).UpdateNotify(order.Id, payNotify); err != nil {
				l.Errorf("修改代收订单[%v]通知状态失败, err=[%v]", orderNo, err)
			}
		} else {
			// redis记录数据
			redisKey := GetPayNotifyExpireKey(orderNo)
			redisRet := l.svcCtx.Redis.SetNX(context.TODO(), redisKey, order.NotifyFailTimes, time.Duration(waitTime)*time.Second)
			if redisRet.Err() != nil {
				l.Errorf("redis插入数据失败, key=%v, err=%v", redisKey, err)
			}

			// 修改代收订单异步通知信息-异步通知进行中
			payNotify := model.PayNotifyInfo{
				NotifyStatus:    model.PayNotifyStatusNotifying,
				NotifyFailTimes: order.NotifyFailTimes + 1,
				NextNotifyTime:  time.Now().Add(time.Duration(waitTime) * time.Second).Unix(),
			}
			if err := model.NewPayOrderModel(l.svcCtx.DbEngine).UpdateNotify(order.Id, payNotify); err != nil {
				l.Errorf("修改代收订单[%v]通知状态失败, err=[%v]", orderNo, err)
			}
		}

		l.Infof("代收订单第%v次通知结束", order.NotifyFailTimes+1)
		return
	}
}
