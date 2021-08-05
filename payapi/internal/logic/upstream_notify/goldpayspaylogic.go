package upstream_notify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
	"tpay_backend/model"
	"tpay_backend/payapi/internal/logic"
	"tpay_backend/upstream"

	"github.com/tal-tech/go-zero/core/logx"
	"tpay_backend/payapi/internal/svc"
)

type GoldPaysPayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGoldPaysPayLogic(ctx context.Context, svcCtx *svc.ServiceContext) GoldPaysPayLogic {
	return GoldPaysPayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GoldPaysPayLogic) GoldPaysPay(body []byte) error {
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
		return errors.New(fmt.Sprintf("解析json参数失败:%v, body:%v", err, string(body)))
	}

	// 2.验证参数
	if reqData.AppId == "" || reqData.Sign == "" || reqData.OrderNo == "" || reqData.ReqNo == "" || reqData.Amount == "" || reqData.OrderStatus == "" {
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

	logx.Infof("上游信息:%+v", up)

	upObj, err := logic.NewFuncLogic(l.svcCtx).GetUpstreamObject(up)
	if err != nil {
		logx.Errorf("获取上游对象失败err:%v,upstream:%+v", err, up)
		return errors.New("获取上游对象失败")
	}

	// 4.校验签名
	dataMap := make(map[string]interface{})
	if err := json.Unmarshal(body, &dataMap); err != nil {
		logx.Errorf("解析body到map失败err:%v", err)
		return errors.New("解析body到map失败")
	}

	if err := upObj.CheckSign(dataMap); err != nil {
		logx.Errorf("校验签名失败err:%v,dataMap:%+v", err, dataMap)
		return errors.New("校验签名失败")
	}

	// 5.查询订单
	order, err := model.NewPayOrderModel(l.svcCtx.DbEngine).FindOneByOrderNo(reqData.ReqNo)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("订单[%v]不存在", reqData.ReqNo)
		} else {
			l.Errorf("查询代收订单[%v]失败, err=%v", reqData.ReqNo, err)
		}
		return err
	}

	logx.Infof("订单信息:%+v", order)

	if order.OrderStatus == model.PayOrderStatusPaid {
		l.Errorf("代收订单已支付，重复通知, order.OrderNo:%v", order.OrderNo)
		return nil
	}

	if order.OrderStatus != model.PayOrderStatusPending {
		l.Errorf("代收订单不是待支付订单, order.OrderNo:%v, order.OrderStatus:%v", order.OrderNo, order.OrderStatus)
		return errors.New("订单状态不允许")
	}

	// 上游金额单位是卢比 平台金额单位是分 1卢比=100分
	orderAmount := decimal.NewFromInt(order.ReqAmount).Div(decimal.NewFromInt(100)).Round(2).String()
	if orderAmount != reqData.Amount {
		l.Errorf("订单[%v]金额不对, order.reqAmount:%v, reqData.Amount:%v", reqData.ReqNo, order.ReqAmount, reqData.Amount)
		return errors.New("订单金额不对")
	}

	paymentAmount, err := strconv.ParseInt(reqData.Amount, 10, 64)
	if err != nil {
		l.Errorf("字符串转int64失败, err=%v", err)
		return err
	}

	order.PaymentAmount = paymentAmount

	// 6.同步订单信息
	var orderStatus int64
	switch reqData.OrderStatus {
	case upstream.GoldPaysOrderStatusPaySuccess:
		orderStatus = model.PayOrderStatusPaid
	case upstream.GoldPaysOrderStatusPayFail:
		orderStatus = model.PayOrderStatusFail
	default:
		l.Errorf("上游通知的是一个未知的订单状态, order_status:%v", reqData.OrderStatus)
		return errors.New("订单状态不对")
	}
	if err := NewSyncOrder(context.TODO(), l.svcCtx).SyncPayOrder(order, orderStatus, reqData.Msg); err != nil {
		l.Errorf("同步订单信息, orderNo:%v, MerchantNo:%v, err:%v", order.OrderNo, order.MerchantNo, err)
		return err
	}

	return nil
}
