package upstream

import (
	"encoding/json"
	"github.com/tal-tech/go-zero/core/logx"
	"log"
	"testing"
	"tpay_backend/utils"
)

func Get777PayObj() Upstream {
	upMerchantNo := "vn1067"
	jsonStr := `{
    "Host": "https://api.zf77777.org",
    "SecretKey": "8f581647-7112-4aa3-82bb-e6171ea732aa",
    "PayNotifyPath": "/notify/777pay/pay",
    "TransferNotifyPath": "/notify/777pay/transfer"
}`
	obj, err := NewThreeSevenPay(upMerchantNo, jsonStr)
	if err != nil {
		log.Fatalf("获取对象失败:%v", err)
	}

	return obj
}

func TestThreeSevenPay_Pay(t *testing.T) {
	up := Get777PayObj()

	req := &PayRequest{
		Amount:       6000,
		Currency:     "",
		OrderNo:      utils.GetDailyId(),
		NotifyUrl:    "https://tpay-api.mangopay-test.com/notify/zf777pay/pay",
		ReturnUrl:    "",
		ProductType:  "vietinbankipay",
		CustomName:   "",
		CustomMobile: "",
		CustomEmail:  "",
		Attach:       "",
	}
	resp, err := up.Pay(req)
	if err != nil {
		t.Errorf("err=%v", err)
		return
	}

	t.Logf("结果：%+v", resp)
}

func TestThreeSevenPay_PayOrderQuery(t *testing.T) {
	up := Get777PayObj()

	req := &PayOrderQueryRequest{
		OrderNo:         "",
		UpstreamOrderNo: "67450843-3b94-480d-b74b-27ca88a12490",
	}
	resp, err := up.PayOrderQuery(req)
	if err != nil {
		t.Errorf("err=%v", err)
		return
	}

	t.Logf("结果：%+v", resp)
}

func TestThreeSevenPay_Transfer(t *testing.T) {
	up := Get777PayObj()

	req := &TransferRequest{
		Amount:             6000,
		Currency:           "",
		OrderNo:            utils.GetDailyId(),
		NotifyUrl:          "https://tpay-api.mangopay-test.com/notify/zf777pay/transfer",
		ReturnUrl:          "",
		ProductType:        "bank",
		Attach:             "测试",
		BankName:           "Bank Central Asia",
		BankCardNo:         "789465156266",
		BankCode:           "",
		BankCardHolderName: "张三",
		CardHolderMobile:   "",
		CardHolderEmail:    "",
	}
	resp, err := up.Transfer(req)
	if err != nil {
		t.Errorf("err=%v", err)
		return
	}

	t.Logf("结果：%+v", resp)

}

func TestThreeSevenPay_TransferOrderQuery(t *testing.T) {
	up := Get777PayObj()

	req := &TransferOrderQueryRequest{
		OrderNo:         "",
		UpstreamOrderNo: "02f526be-fb34-4765-a787-19daa9ba3496",
	}
	resp, err := up.TransferOrderQuery(req)
	if err != nil {
		t.Errorf("TransferOrderQuery() error = %v", err)
		return
	}
	t.Logf("结果：%+v", resp)
}

func TestThreeSevenPay_TransferOrderQuery2(t *testing.T) {

	body := `{"payamount":500,"mark":"test","ordertype":2,"iscancel":1,"ticket":"b6bfca5a-68e8-45c6-a1f4-17002f51a4ff","userid":"vn1067","orderid":"8161915873361994545274","type":"bank","sign":"0d2b66d7378ed39b809942a20566b615","pageurl":"https://form.zf77777.org/api/paypage?ticket=b6bfca5a-68e8-45c6-a1f4-17002f51a4ff","amount":7000,"bmount":"7000.00","serialno":null,"upi":null,"qrcode":null,"note":"","success":1,"message":null}`
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
	if err := json.Unmarshal([]byte(body), &reqData); err != nil {
		t.Errorf("err=%v", err)
		return
	}
	t.Logf("结果：%+v", reqData)
}

func TestThreeSevenPay_Sign(t *testing.T) {
	data := make(map[string]interface{})
	data["orderid"] = "qqq"
	data["amount"] = 100.00
	up := Get777PayObj()
	sign := up.GenerateSign(data)
	t.Logf("结果：%+v", sign)
}

func TestThreeSevenPay_CheckSign(t *testing.T) {
	body := `{"paycode":"test","payamount":16900,"payamountd":16900.0,"paytime":"2021-04-23 14:01:30","payuser":"yeya","ispay":1,"ticket":"48daf267-443d-4276-a171-0a8758c3276c","userid":"vn1067","orderid":"7161915758863294213334","type":"vietcombank","sign":"1ecaa1c900d0f6a2e0b154f99507b28b","pageurl":"https://form.zf77777.org/api/paypage?ticket=48daf267-443d-4276-a171-0a8758c3276c","amount":16900,"bmount":"16900.00","serialno":"9460","upi":null,"qrcode":null,"note":"","success":1,"message":null}`
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
	if err := json.Unmarshal([]byte(body), &reqData); err != nil {
		t.Errorf("err=%v", err)
		return
	}

	up := Get777PayObj()

	// 4.校验签名
	dataMap := make(map[string]interface{})
	if reqData.IsCancel == ThreeSevenPayCancelled {
		dataMap["orderid"] = reqData.OrderId
		dataMap["amount"] = reqData.Amount
		dataMap["sign"] = reqData.Sign
	} else {
		dataMap["orderid"] = reqData.OrderId
		dataMap["amount"] = reqData.PayAmount
		dataMap["sign"] = reqData.Sign
	}
	if err := up.CheckSign(dataMap); err != nil {
		logx.Errorf("校验签名失败err:%v,dataMap:%+v", err, dataMap)
		return
	}

	t.Logf("结果")
}
