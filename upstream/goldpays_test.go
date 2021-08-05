package upstream

import (
	"log"
	"strings"
	"testing"
	"tpay_backend/utils"
)

func GetGoldPaysObj() Upstream {
	upMerchantNo := "C1618479683363"
	jsonStr := `{
    "Host": "https://www.goldpays.in",
    "SecretKey": "FmoBBrtFMOeceHOzzFXOdgdLCskXQEMpD8IiMn",
    "PayNotifyPath": "/notify/goldpays/pay",
    "TransferNotifyPath": "/notify/goldpays/transfer"
}`
	obj, err := NewThreeSevenPay(upMerchantNo, jsonStr)
	if err != nil {
		log.Fatalf("获取对象失败:%v", err)
	}

	return obj
}

//{"@timestamp":"2021-04-19T15:21:32.687+08","level":"info","content":"tpay_backend/upstream.(*GoldPays).Pay:response:{\"code\":200,\"success\":true,\"data\":{\"merchant\":\"C1618479683363\",\"orderId\":\"T2021041915213197082580\",\"platOrderId\":\"304275327377604608\",\"url\":\"https://www.goldpays.in/paymentPage?mid=C1618479683363\u0026oid=304275327377604608\u0026s=037907b8c9f0c91d109e6fb08bfc47b8\",\"sign\":\"a68ecbbd4655b49899434c0a33a0e57a\"}}"}
func TestGoldPays_Pay(t *testing.T) {
	up := GetGoldPaysObj()

	req := &PayRequest{
		Amount:       10056,
		Currency:     "",
		OrderNo:      "P" + utils.GetDailyId(),
		NotifyUrl:    "https://tpay-api.mangopay-test.com/notify/goldpays/pay",
		ReturnUrl:    "/xxx/x",
		ProductType:  "",
		CustomName:   "李明",
		CustomMobile: "138690108",
		CustomEmail:  "h132@193.com",
	}
	resp, err := up.Pay(req)
	if err != nil {
		t.Errorf("请求上游失败, err=%v", err)
		return
	}

	t.Logf("结果：%+v", resp)
}

func TestGoldPays_PayOrderQuery(t *testing.T) {
	up := GetGoldPaysObj()

	//T2021041915213197082580
	req := &PayOrderQueryRequest{
		OrderNo: "T2021041915410615962430",
	}
	resp, err := up.PayOrderQuery(req)
	if err != nil {
		t.Errorf("请求上游失败, err=%v", err)
		return
	}

	t.Logf("结果：%+v", resp)
}

func TestGoldPays_QueryBalance(t *testing.T) {
	up := GetGoldPaysObj()

	resp, err := up.QueryBalance()

	if err != nil {
		t.Errorf("请求上游失败, err=%v", err)
		return
	}

	t.Logf("结果：%+v", resp)
}

func TestGoldPays_Transfer(t *testing.T) {
	up := GetGoldPaysObj()

	req := &TransferRequest{
		Amount:      5000,
		OrderNo:     "T" + utils.GetDailyId(),
		NotifyUrl:   "https://tpay-api.mangopay-test.comnotify/goldpays/transfer",
		ReturnUrl:   "/xxx/x",
		ProductType: GoldPaysPayOutModeIMPS,
		BankName:    "中国建设银行",
		BankCardNo:  "601700721234598",
		//BankCode:           "454555",
		BankCardHolderName: "李四",
		CardHolderMobile:   "12312323",
		CardHolderEmail:    "rtt@gg.com",
	}

	resp, err := up.Transfer(req)
	if err != nil {
		t.Errorf("请求上游失败, err=%v", err)
		return
	}

	t.Logf("结果：%+v", resp)
}

func TestGoldPays_TransferOrderQuery(t *testing.T) {
	up := GetGoldPaysObj()

	req := &TransferOrderQueryRequest{
		OrderNo: "T2021041915410615962430",
	}
	resp, err := up.TransferOrderQuery(req)
	if err != nil {
		t.Errorf("请求上游失败, err=%v", err)
		return
	}

	t.Logf("结果：%+v", resp)
}

func TestGoldPays_TransferOrderQuery2(t *testing.T) {
	a := "/notify/zf777pay/pay"
	a = strings.TrimPrefix(a, "/")
	t.Logf("结果：%+v", a)
}
