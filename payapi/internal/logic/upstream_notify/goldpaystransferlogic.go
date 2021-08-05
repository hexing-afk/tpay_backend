package upstream_notify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"tpay_backend/model"
	"tpay_backend/payapi/internal/logic"
	"tpay_backend/upstream"

	"github.com/tal-tech/go-zero/core/logx"
	"tpay_backend/payapi/internal/svc"
)

type GoldPaysTransferLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGoldPaysTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) GoldPaysTransferLogic {
	return GoldPaysTransferLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GoldPaysTransferLogic) GoldPaysTransfer(body []byte) error {
	var reqData struct {
		AppId       string `json:"merchant"`    // 商户号
		ReqNo       string `json:"orderId"`     // 商户订单号
		OrderNo     string `json:"platOrderId"` // 平台订单号
		Amount      string `json:"amount"`      // 金额
		Msg         string `json:"msg"`         // 处理消息
		OrderStatus string `json:"status"`      // 交易状态，值见数据字典
		Sign        string `json:"sign"`        // 签名
	}

	// 1.解析接口数据
	if err := json.Unmarshal(body, &reqData); err != nil {
		return errors.New(fmt.Sprintf("解析json参数失败:%v", err))
	}

	// 2.参数验证
	if reqData.AppId == "" || reqData.ReqNo == "" || reqData.OrderNo == "" || reqData.OrderStatus == "" || reqData.Sign == "" || reqData.Amount == "" {
		return errors.New(fmt.Sprintf("缺少必须参数,reqData:%+v", reqData))
	}

	// 3.获取上游
	up, err := model.NewUpstreamModel(l.svcCtx.DbEngine).FindOneByUpstreamMerchantNo(reqData.AppId)
	if err != nil {
		if err == model.ErrRecordNotFound {
			return errors.New(fmt.Sprintf("未找到对应的上游:UpstreamMerchantNo:%v", reqData.AppId))
		} else {
			return errors.New(fmt.Sprintf("查询上游信息失败:err:%v,UpstreamMerchantNo:%v", err, reqData.AppId))
		}
	}

	upObj, err := logic.NewFuncLogic(l.svcCtx).GetUpstreamObject(up)
	if err != nil {
		logx.Errorf("获取上游对象失败err:%v,upstream:%+v", err, up)
		return errors.New("获取上游对象失败")
	}

	// 4.验签
	dataMap := make(map[string]interface{})
	if err := json.Unmarshal(body, &dataMap); err != nil {
		logx.Errorf("解析body到map失败err:%v", err)
		return errors.New("解析body到map失败")
	}

	if err := upObj.CheckSign(dataMap); err != nil {
		l.Errorf("验签失败, data:%+v, err:%v", dataMap, err)
		return err
	}

	// 5.查询订单
	order, err := model.NewTransferOrderModel(l.svcCtx.DbEngine).FindByOrderNo(reqData.ReqNo)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("订单[%v]不存在", reqData.ReqNo)
			return errors.New(fmt.Sprintf("找不到订单[%v]", reqData.ReqNo))
		} else {
			return errors.New("查询订单失败")
		}
	}

	l.Infof("订单信息：%+v", order)

	if order.OrderStatus == model.TransferOrderStatusPaid {
		l.Errorf("代付订单已支付，重复通知, order.OrderNo:%v", order.OrderNo)
		return nil
	}

	if order.OrderStatus != model.TransferOrderStatusPending {
		l.Errorf("代付订单不是待支付订单, order.OrderNo:%v, order.OrderStatus:%v", order.OrderNo, order.OrderStatus)
		return errors.New("订单状态不允许")
	}

	// 上游金额单位是卢比 平台金额单位是分 1卢比=100分
	orderAmount := decimal.NewFromInt(order.ReqAmount).Div(decimal.NewFromInt(100)).Round(2).String()
	if orderAmount != reqData.Amount {
		l.Errorf("订单[%v]金额不对, order.reqAmount:%v, reqData.Amount:%v", reqData.ReqNo, order.ReqAmount, reqData.Amount)
		return errors.New("订单金额不对")
	}

	// 6.同步订单信息
	var orderStatus int64
	switch reqData.OrderStatus {
	case upstream.GoldPaysOrderStatusPaySuccess:
		orderStatus = model.TransferOrderStatusPaid
	case upstream.GoldPaysOrderStatusPayFail:
		orderStatus = model.TransferOrderStatusFail
	default:
		l.Errorf("上游通知的是一个未知的订单状态, order_status:%v", reqData.OrderStatus)
		return errors.New("订单状态不对")
	}
	if err := NewSyncOrder(context.TODO(), l.svcCtx).SyncTransferOrder(order, orderStatus, reqData.Msg); err != nil {
		l.Errorf("同步订单信息, orderNo:%v, MerchantNo:%v, err:%v", order.OrderNo, order.MerchantNo, err)
		return err
	}

	return nil
}
