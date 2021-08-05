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

type TransferOrderNotify struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewTransferOrderNotify(ctx context.Context, svcCtx *svc.ServiceContext) *TransferOrderNotify {
	return &TransferOrderNotify{
		Logger: logx.WithContext(ctx),
		svcCtx: svcCtx,
	}
}

func (l *TransferOrderNotify) OrderNotify(orderNo string) {
	// 1.查询订单
	order, err := model.NewTransferOrderModel(l.svcCtx.DbEngine).FindByOrderNo(orderNo)
	if err != nil {
		l.Errorf("查询代付订单[%v]失败, err=%v", orderNo, err)
		return
	}

	l.Infof("代付订单信息：%+v", order)

	//订单-通知成功
	if order.NotifyStatus == model.TransferNotifyStatusSuccess {
		l.Errorf("代付订单[%v]已经成功通知", orderNo)
		return
	}

	//订单-通知超时
	if order.NotifyStatus == model.TransferNotifyStatusTimeOut {
		l.Errorf("代付订单[%v]已超出通知时间", orderNo)
		return
	}

	if order.NotifyUrl == "" {
		l.Errorf("代付订单[%v]缺少回调地址，不再进行通知", orderNo)
		//插入通知日志
		data := &model.OrderNotifyLog{
			OrderNo:         order.OrderNo,
			MerchantOrderNo: order.MerchantOrderNo,
			Status:          model.OrderNotifyStatusFail,
			OrderType:       model.NotifyLogOrderTypeTransfer,
			Result:          "订单缺少回调地址",
		}
		if err := model.NewOrderNotifyLogModel(l.svcCtx.DbEngine).Insert(data); err != nil {
			l.Errorf("代付订单异步通知日志记录失败, data：%+v, err：%v", data, err)
			return
		}

		transferNotify := model.TransferNotifyInfo{
			NotifyStatus:    model.TransferNotifyStatusTimeOut,
			NotifyFailTimes: order.NotifyFailTimes + 1,
			NextNotifyTime:  0,
		}

		if err := model.NewTransferOrderModel(l.svcCtx.DbEngine).UpdateNotify(order.Id, transferNotify); err != nil {
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
	packReq := &utils.PackTransferNotifyParamsRequest{
		MerchantNo:      merchant.MerchantNo,
		Timestamp:       time.Now().Unix(),
		NotifyType:      utils.TransferNotifyType,
		OrderNo:         order.OrderNo,
		MerchantOrderNo: order.MerchantOrderNo,
		ReqAmount:       order.ReqAmount,
		Currency:        order.Currency,
		OrderStatus:     order.OrderStatus,
		PayTime:         order.UpdateTime,
	}
	postBytes, err := utils.PackTransferNotifyParams(packReq, merchant.Md5Key)
	if err != nil {
		l.Errorf("打包数据失败, err=[%v]", err)
		return
	}

	l.Infof("代付订单[%v]第%v次通知, 失败次数=[%v]", order.OrderNo, order.NotifyFailTimes+1, order.NotifyFailTimes)

	// 4.发送通知,接收通知结果
	result, err := utils.PostForm(order.NotifyUrl, postBytes)
	if err != nil {
		l.Errorf("发送失败, utils.Post() error: %v", err)
		result = []byte(err.Error())
	}

	resultStr := string(result)
	resultMap := map[string]string{
		"url": order.NotifyUrl,
		"res": resultStr,
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
			OrderType:       model.NotifyLogOrderTypeTransfer,
		}
		if err := model.NewOrderNotifyLogModel(l.svcCtx.DbEngine).Insert(log); err != nil {
			l.Errorf("代付订单异步通知日志记录失败, data：%+v, err：%v", log, err)
		}

		// 修改代收订单异步通知信息-通知成功
		transferNotify := model.TransferNotifyInfo{
			NotifyStatus:    model.PayNotifyStatusSuccess,
			NotifyFailTimes: order.NotifyFailTimes,
			NextNotifyTime:  0,
		}
		if err := model.NewTransferOrderModel(l.svcCtx.DbEngine).UpdateNotify(order.Id, transferNotify); err != nil {
			l.Errorf("修改代付订单[%v]通知状态失败, err=[%v]", orderNo, err)
		}

		l.Infof("代付订单第%v次通知结束", order.NotifyFailTimes+1)
		return
	default:
		// 插入通知失败日志
		log := &model.OrderNotifyLog{
			OrderNo:         order.OrderNo,
			MerchantOrderNo: order.MerchantOrderNo,
			Status:          model.OrderNotifyStatusFail,
			OrderType:       model.NotifyLogOrderTypeTransfer,
			Result:          string(resultJ),
		}
		if err := model.NewOrderNotifyLogModel(l.svcCtx.DbEngine).Insert(log); err != nil {
			l.Errorf("代付订单异步通知日志记录失败, data：%+v, err：%v", log, err)
		}

		// 下一次通知的等待时间
		waitTime := GetNotifyWaitTimeById(order.NotifyFailTimes + 1)
		if waitTime == -1 {
			// 修改代收订单异步通知信息-通知超时
			transferNotify := model.TransferNotifyInfo{
				NotifyStatus:    model.TransferNotifyStatusTimeOut,
				NotifyFailTimes: order.NotifyFailTimes + 1,
				NextNotifyTime:  0,
			}
			if err := model.NewTransferOrderModel(l.svcCtx.DbEngine).UpdateNotify(order.Id, transferNotify); err != nil {
				l.Errorf("修改代付订单[%v]通知状态失败, err=[%v]", orderNo, err)
			}
		} else {
			// redis记录数据
			redisKey := GetTransferNotifyExpireKey(orderNo)
			redisRet := l.svcCtx.Redis.SetNX(context.TODO(), redisKey, order.NotifyFailTimes, time.Duration(waitTime)*time.Second)
			if redisRet.Err() != nil {
				l.Errorf("redis插入数据失败, key=%v, err=%v", redisKey, err)
			}

			// 修改代收订单异步通知信息-通知进行中
			transferNotify := model.TransferNotifyInfo{
				NotifyStatus:    model.TransferNotifyStatusNotifying,
				NotifyFailTimes: order.NotifyFailTimes + 1,
				NextNotifyTime:  time.Now().Add(time.Duration(waitTime) * time.Second).Unix(),
			}
			if err := model.NewTransferOrderModel(l.svcCtx.DbEngine).UpdateNotify(order.Id, transferNotify); err != nil {
				l.Errorf("修改代付订单[%v]通知状态失败, err=[%v]", orderNo, err)
			}
		}

		l.Infof("代付订单第%v次通知结束", order.NotifyFailTimes+1)
		return
	}
}
