package upstream_notify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"tpay_backend/model"
	"tpay_backend/payapi/internal/logic"
	"tpay_backend/upstream"

	"github.com/tal-tech/go-zero/core/logx"
	"tpay_backend/payapi/internal/svc"
)

type ThreeSevenPayTransferLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewThreeSevenPayTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) ThreeSevenPayTransferLogic {
	return ThreeSevenPayTransferLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThreeSevenPayTransferLogic) ThreeSevenPayTransfer(body []byte) error {
	var reqData struct {
		Success   int64   `json:"success"`   // 请求是否成功 1、成功；0、失败
		Message   string  `json:"message"`   // 出错消息，请求处理失败才会出现
		Ticket    string  `json:"ticket"`    // 访问票据
		IsPay     int64   `json:"ispay"`     // 是否支付，0 没有支付 1 已经支付
		PayCode   string  `json:"paycode"`   // 支付代码	支付网关返回编码
		PayAmount float64 `json:"payamount"` // 支付金额 支付网关返回的实际金额，业务逻辑中应使用此金额作为入金金额而非定单金额
		PayTime   string  `json:"msg"`       // 支付时间	字符串类型格式为： 2000-01-01 23:34:56
		PayUser   string  `json:"status"`    // 支付用户
		Sign      string  `json:"sign"`      // 签名
		Amount    int64   `json:"amount"`    // 创建订单时的金额，原样返回
		Note      string  `json:"note"`      // 创建订单时的备注，原样返回
		UserId    string  `json:"userid"`    // 商户编号
		OrderId   string  `json:"orderid"`   // 商户订单号
		PayType   string  `json:"type"`      // 支付类型
		SerialNo  string  `json:"serialno"`  // 支付备注
		BMount    string  `json:"bmount"`    // 尾部金额

		IsCancel  int64  `json:"iscancel"`  // 是否被取消 0 没有取消 1 已经取消
		OrderType int64  `json:"ordertype"` // 订单类型 1=支付充值订单 2=代付提现订单
		Mark      string `json:"mark"`      // 订单取消原因
	}

	// 1.解析接口数据
	if err := json.Unmarshal(body, &reqData); err != nil {
		return errors.New(fmt.Sprintf("解析json参数失败:%v, body:%v", err, string(body)))
	}

	// 2.验证参数
	if reqData.UserId == "" || reqData.Sign == "" || reqData.OrderId == "" || reqData.Ticket == "" {
		return errors.New(fmt.Sprintf("缺少必须参数,reqData:%+v", reqData))
	}

	// 3.获取上游
	up, err := model.NewUpstreamModel(l.svcCtx.DbEngine).FindOneByUpstreamMerchantNo(reqData.UserId)
	if err != nil {
		if err == model.ErrRecordNotFound {
			return errors.New(fmt.Sprintf("未找到对应的上游:UpstreamMerchantNo:%v", reqData.UserId))
		} else {
			return errors.New(fmt.Sprintf("查询上游信息失败:err:%v,UpstreamMerchantNo:%v", err, reqData.UserId))
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
	if reqData.IsCancel == upstream.ThreeSevenPayCancelled {
		dataMap["orderid"] = reqData.OrderId
		dataMap["amount"] = strconv.FormatInt(reqData.Amount, 10)
		dataMap["sign"] = reqData.Sign
	} else {
		dataMap["orderid"] = reqData.OrderId
		dataMap["amount"] = strconv.FormatFloat(reqData.PayAmount, 'f', -1, 64)
		dataMap["sign"] = reqData.Sign
	}
	if err := upObj.CheckSign(dataMap); err != nil {
		logx.Errorf("校验签名失败err:%v,dataMap:%+v", err, dataMap)
		return errors.New("校验签名失败")
	}

	// 5.查询订单
	order, err := model.NewTransferOrderModel(l.svcCtx.DbEngine).FindByOrderNo(reqData.OrderId)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("订单[%v]不存在", reqData.OrderId)
			return errors.New(fmt.Sprintf("找不到订单[%v]", reqData.OrderId))
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

	if order.ReqAmount != reqData.Amount {
		l.Errorf("订单[%v]金额不对, order.reqAmount:%v, reqData.Amount:%v", reqData.OrderId, order.ReqAmount, reqData.Amount)
		return errors.New("订单金额不对")
	}

	// 6.同步订单信息
	var orderStatus int64
	var failReason string
	if reqData.IsPay == upstream.ThreeSevenPayPaid {
		orderStatus = model.PayOrderStatusPaid
	} else if reqData.IsCancel == upstream.ThreeSevenPayCancelled {
		orderStatus = model.PayOrderStatusFail
		failReason = reqData.Mark
	} else {
		l.Errorf("上游通知的是一个未知的订单状态, reqData.IsPay:%v, reqData.IsCancel", reqData.IsPay, reqData.IsCancel)
		return errors.New("订单状态不对")
	}

	if err := NewSyncOrder(context.TODO(), l.svcCtx).SyncTransferOrder(order, orderStatus, failReason); err != nil {
		l.Errorf("同步订单信息, orderNo:%v, MerchantNo:%v, err:%v", order.OrderNo, order.MerchantNo, err)
		return err
	}

	return nil
}
