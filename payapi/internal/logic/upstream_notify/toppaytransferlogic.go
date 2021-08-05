package upstream_notify

import (
	"context"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/tal-tech/go-zero/core/logx"
	"net/url"
	"strings"
	"tpay_backend/model"
	"tpay_backend/payapi/internal/svc"
	"tpay_backend/upstream"
)

type ToppayTransferLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type ToppayNotifyRequest struct {
	Data               string `form:"data"`
	Sign               string `form:"sign"`
	UpstreamMerchantNo string `form:"merchant_no"`
}

func NewToppayTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) ToppayTransferLogic {
	return ToppayTransferLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ToppayTransferLogic) ToppayTransfer(body []byte) error {
	bodyStr := string(body)
	bodyStr = strings.ReplaceAll(bodyStr, "\\u003c", "<")
	bodyStr = strings.ReplaceAll(bodyStr, "\\u003e", ">")
	bodyStr = strings.ReplaceAll(bodyStr, "\\u0026", "&")

	fmt.Printf("bodyStr:%v\n", bodyStr)

	//解析url字符串
	u, err := url.ParseQuery(bodyStr)
	if err != nil {
		panic(err)
	}

	var reqData struct {
		Appid  string `json:"appid"`   //请求应用appid	商户注册时申请的应用编号，由运营申请分配
		TfSign string `json:"tf_sign"` //	签名字符串	签名数据

		//业务参数
		Subject              string `json:"subject"`              //商品名称	电子产品(笔记本)
		Businessnumber       string `json:"businessnumber"`       //业务订单号	20110110165455000001
		Businessrecordnumber string `json:"businessrecordnumber"` //支付订单号	20110110165455000001
		Status               string `json:"status"`               //交易状态： 成功/失败/已退款	成功
		Billamount           string `json:"billamount"`           //订单金额	单位：元，99.99
		Transactionamount    string `json:"transactionamount"`    //交易金额	单位：元，99.99
		Transactiontype      string `json:"transactiontype"`      //交易类型
		Inputdate            string `json:"inputdate"`            //交易创建时间，格式为yyyy-MM-dd HH:mm:ss
		Remark               string `json:"remark"`               //结果说明
		Sss                  string `json:"sss"`                  //结果说明
		Transactiondate      string `json:"transactiondate"`      //结果说明
	}

	for k, v := range u {
		fmt.Printf("u[%v]=%v\n", k, v)
	}

	fmt.Println("appid-------------", u.Get("appid"))

	reqData.Appid = u.Get("appid")
	reqData.TfSign = u.Get("tfSign")
	reqData.Subject = u.Get("subject")
	reqData.Businessnumber = u.Get("businessnumber")
	reqData.Businessrecordnumber = u.Get("businessrecordnumber")
	reqData.Status = u.Get("status")
	reqData.Billamount = u.Get("billamount")
	reqData.Transactionamount = u.Get("transactionamount")
	reqData.Transactiontype = u.Get("transactiontype")
	reqData.Inputdate = u.Get("inputdate")
	reqData.Remark = u.Get("remark")
	reqData.Sss = u.Get("sss")

	if reqData.Sss != "123456a87g2GSG&*^Ihgqasrg" {
		return errors.New(fmt.Sprintf("签名错误,reqData.Sss:%+v", reqData.Sss))
	}

	// 1.解析接口数据
	//if err := json.Unmarshal(body, &reqData); err != nil {
	//	return errors.New(fmt.Sprintf("解析json参数失败:%v, body:%v", err, string(body)))
	//}

	// 2.验证参数
	if reqData.Businessnumber == "" || reqData.Businessrecordnumber == "" {
		return errors.New(fmt.Sprintf("缺少必须参数,reqData:%+v", reqData))
	}

	// 5.查询订单
	order, err := model.NewTransferOrderModel(l.svcCtx.DbEngine).FindByOrderNo(reqData.Businessnumber)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("订单[%v]不存在", reqData.Businessnumber)
		} else {
			l.Errorf("查询代付订单[%v]失败, err=%v", reqData.Businessnumber, err)
		}
		return err
	}

	logx.Infof("订单信息:%+v", order)

	if order.OrderStatus == model.PayOrderStatusPaid {
		l.Errorf("代付订单已支付，重复通知, order.OrderNo:%v", order.OrderNo)
		return nil
	}

	if order.OrderStatus != model.PayOrderStatusPending {
		l.Errorf("代付订单不是待支付订单, order.OrderNo:%v, order.OrderStatus:%v", order.OrderNo, order.OrderStatus)
		return errors.New("订单状态不允许")
	}

	currency, err := model.NewCurrencyModel(l.svcCtx.DbEngine).FindByCurrency(order.Currency)
	if err != nil {
		l.Errorf("查询币种失败, err:%v", err)
		return errors.New("币种不存在")
	}

	billamountDe, err := decimal.NewFromString(reqData.Billamount)
	if err != nil {
		l.Errorf("金额有误, err:%v", err)
		return errors.New("金额有误")
	}
	transactionamountDe, err := decimal.NewFromString(reqData.Transactionamount)
	if err != nil {
		l.Errorf("金额有误, err:%v", err)
		return errors.New("金额有误")
	}

	var billamount int64 = 0
	var transactionamount int64 = 0

	if currency.IsDivideHundred == model.DivideHundred {
		billamount = billamountDe.Mul(decimal.NewFromInt(100)).IntPart()
		transactionamount = transactionamountDe.Mul(decimal.NewFromInt(100)).IntPart()
	}

	if order.ReqAmount != billamount {
		l.Errorf("订单[%v]金额不对, order.reqAmount:%v, reqData.Billamount:%v", order.OrderNo, order.ReqAmount, billamount)
		return errors.New("订单金额不对")
	}

	order.PayeeRealAmount = transactionamount

	// 6.同步订单信息
	var orderStatus int64
	var failReason string

	switch reqData.Status {
	case upstream.TopPayUnpaid:
		orderStatus = model.PayOrderStatusPaid
	case upstream.TopPayRefund:
		orderStatus = model.PayOrderStatusFail
		failReason = reqData.Remark
	default:
		l.Errorf("上游通知的是一个未知的订单状态, reqData.Status:%v", reqData.Status)
		return errors.New("订单状态不对")
	}

	if err := NewSyncOrder(context.TODO(), l.svcCtx).SyncTransferOrder(order, orderStatus, failReason); err != nil {
		l.Errorf("同步订单信息, orderNo:%v, MerchantNo:%v, err:%v", order.OrderNo, order.MerchantNo, err)
		return err
	}

	return nil
}
