package upstream_notify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"tpay_backend/model"
	"tpay_backend/payapi/internal/logic"
	"tpay_backend/payapi/internal/svc"
	"tpay_backend/upstream"

	"github.com/tal-tech/go-zero/core/logx"
)

type XPayPayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type XPayNotifyRequest struct {
	Data               string `form:"data"`
	Sign               string `form:"sign"`
	UpstreamMerchantNo string `form:"merchant_no"`
}

func NewXPayPayLogic(ctx context.Context, svcCtx *svc.ServiceContext) XPayPayLogic {
	return XPayPayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *XPayPayLogic) XPayPay(req XPayNotifyRequest) error {
	// 1.验证参数
	if req.UpstreamMerchantNo == "" || req.Sign == "" || req.Data == "" {
		return errors.New(fmt.Sprintf("缺少必须参数,reqData:%+v", req))
	}

	var data struct {
		MerchantNo      string `json:"merchant_no"`    // 商户编号
		Timestamp       int64  `json:"timestamp"`      // 时间戳
		NotifyType      string `json:"notify_type"`    // 通知类型
		OrderNo         string `json:"order_no"`       // 平台订单号
		MerchantOrderNo string `json:"mch_order_no"`   // 外部订单号
		ReqAmount       int64  `json:"req_amount"`     // 订单请求金额
		Currency        string `json:"currency"`       // 币种
		OrderStatus     int64  `json:"order_status"`   // 订单状态
		PayTime         int64  `json:"pay_time"`       // 支付时间(时间戳)
		Subject         string `json:"subject"`        // 商品名称
		PaymentAmount   int64  `json:"payment_amount"` // 实际支付金额
	}

	if err := json.Unmarshal([]byte(req.Data), &data); err != nil {
		return errors.New(fmt.Sprintf("解析json参数失败:%v, data:%v", err, req.Data))
	}

	// 2.1.验证业务数据参数
	if data.MerchantOrderNo == "" || data.OrderNo == "" || data.ReqAmount == 0 || data.OrderStatus == 0 {
		return errors.New(fmt.Sprintf("缺少必须参数,reqData:%v", req.Data))
	}

	// 3.获取上游
	up, err := model.NewUpstreamModel(l.svcCtx.DbEngine).FindOneByUpstreamMerchantNo(req.UpstreamMerchantNo)
	if err != nil {
		if err == model.ErrRecordNotFound {
			return errors.New(fmt.Sprintf("未找到对应的上游:UpstreamMerchantNo:%v", req.UpstreamMerchantNo))
		} else {
			return errors.New(fmt.Sprintf("查询上游信息失败:err:%v,UpstreamMerchantNo:%v", err, req.UpstreamMerchantNo))
		}
	}

	logx.Infof("上游信息:%+v", up)

	upObj, err := logic.NewFuncLogic(l.svcCtx).GetUpstreamObject(up)
	if err != nil {
		logx.Errorf("获取上游对象失败err:%v,upstream:%+v", err, up)
		return errors.New("获取上游对象失败")
	}

	// 4.校验签名
	dataMap := make(map[string]interface{})
	dataMap["data"] = req.Data
	dataMap["sign"] = req.Sign
	if err := upObj.CheckSign(dataMap); err != nil {
		logx.Errorf("校验签名失败err:%v,dataMap:%+v", err, dataMap)
		return errors.New("校验签名失败")
	}

	// 5.查询订单
	order, err := model.NewPayOrderModel(l.svcCtx.DbEngine).FindOneByOrderNo(data.MerchantOrderNo)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("订单[%v]不存在", data.MerchantOrderNo)
		} else {
			l.Errorf("查询代收订单[%v]失败, err=%v", data.MerchantOrderNo, err)
		}
		return err
	}

	logx.Infof("订单信息:%+v", order)

	if order.OrderStatus == model.PayOrderStatusPaid {
		l.Errorf("代收订单已支付，重复通知, order.OrderNo:%v", order.OrderNo)
		return nil
	}

	//if order.OrderStatus != model.PayOrderStatusPending {
	//	l.Errorf("代收订单不是待支付订单, order.OrderNo:%v, order.OrderStatus:%v", order.OrderNo, order.OrderStatus)
	//	return errors.New("订单状态不允许")
	//}

	// 订单币种
	if order.Currency != data.Currency {
		l.Errorf("订单[%v]金额币种不对, order.Currency:%v, reqData.Currency:%v", data.MerchantOrderNo, order.Currency, data.Currency)
		return errors.New("订单金额币种不对")
	}

	// 上游金额单位:分
	if order.ReqAmount != data.ReqAmount {
		l.Errorf("订单[%v]金额不对, order.reqAmount:%v, reqData.Amount:%v", data.MerchantOrderNo, order.ReqAmount, data.ReqAmount)
		return errors.New("订单金额不对")
	}

	order.PaymentAmount = data.PaymentAmount

	// 6.同步订单信息
	var orderStatus int64
	switch data.OrderStatus {
	case upstream.XPayOrderStatusPaySuccess:
		orderStatus = model.PayOrderStatusPaid
	case upstream.XPayOrderStatusPayFail:
		orderStatus = model.PayOrderStatusFail
	default:
		l.Errorf("上游通知的是一个未知的订单状态, order_status:%v", data.OrderStatus)
		return errors.New("订单状态不对")
	}
	if err := NewSyncOrder(context.TODO(), l.svcCtx).SyncPayOrder(order, orderStatus, ""); err != nil {
		l.Errorf("同步订单信息, orderNo:%v, MerchantNo:%v, err:%v", order.OrderNo, order.MerchantNo, err)
		return err
	}

	return nil
}
