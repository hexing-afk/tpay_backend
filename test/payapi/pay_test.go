package payapi

import (
	"fmt"
	"testing"
	"time"
)

func GetTpayObj() *Tpay {
	return NewTpay(
		//"http://10.41.1.242:8000",
		"http://127.0.0.1:8000",
		"16233081479017453324",
		"1mter8ykmouhzsdrm2yyb9nbuaf9ot1j",
	)
}

func TestTpay_Pay(t *testing.T) {
	tpay := GetTpayObj()

	req := PayReq{}
	req.Subject = "停车费"
	//req.Amount = utils.RandInt64(1000, 10000)
	req.Amount = 150
	req.Currency = "CNY"
	req.MchOrderNo = fmt.Sprintf("PA%d", time.Now().Unix())
	req.TradeType = "toppay"
	req.NotifyUrl = "https://tpay-api.mangopay-test.com/notify/xpay/pay"
	req.ReturnUrl = "https://test-sell-mobile.wesam.vip/#/jumping?order_no"
	req.Attach = "attach_123"

	res, err := tpay.Pay(req)

	t.Logf("err:%v", err)
	t.Logf("res:%+v", res)
}

func TestTpay_PayOrderQuery(t *testing.T) {
	tpay := GetTpayObj()

	req := PayOrderQueryReq{}
	req.MchOrderNo = "PA1618391306"
	//req.OrderNo = "7161839130676685289242"

	res, err := tpay.PayOrderQuery(req)

	t.Logf("err:%v", err)
	t.Logf("res:%+v", res)
}

func TestTpay_Transfer(t *testing.T) {
	tpay := GetTpayObj()

	req := TransferReq{}
	req.Amount = 600                                        // 订单金额
	req.Currency = "USD"                                    // 币种
	req.MchOrderNo = fmt.Sprintf("TR%d", time.Now().Unix()) // 外部订单号(商户系统内部的订单号)
	req.TradeType = "transfer_unionpay"                     // 交易类型
	req.NotifyUrl = "http://xxx.com/notify/transfer"        // 异步通知地址
	req.ReturnUrl = "http://xxx.com/return/transfer"        // 同步跳转地址
	req.Attach = "attach_456"                               // 原样返回字段
	req.BankName = "中国银行"                                   // 收款银行名称
	req.BankCardHolderName = "李白"                           // 银行卡持卡人姓名
	req.BankCardNo = "6216616105001489359"                  // 银行卡号
	req.BankBranchName = "广州分行"                             // 收款银行支行名称

	res, err := tpay.Transfer(req)

	t.Logf("err:%v", err)
	t.Logf("res:%+v", res)
}

func TestTpay_TransferOrderQuery(t *testing.T) {
	tpay := GetTpayObj()

	req := TransferOrderQueryReq{}
	req.MchOrderNo = "TR1618393105"
	req.OrderNo = "8161839310558589810826"

	res, err := tpay.TransferOrderQuery(req)

	t.Logf("err:%v", err)
	t.Logf("res:%+v", res)
}

func TestTpay_QueryBalance(t *testing.T) {
	tpay := GetTpayObj()

	res, err := tpay.QueryBalance()

	t.Logf("err:%v", err)
	t.Logf("res:%+v", res)
}
