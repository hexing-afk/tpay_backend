package upstream

import (
	"log"
	"testing"
	"tpay_backend/utils"
)

func GetXPayConfig() Upstream {
	upMerchantNo := "16196815202336265781"
	jsonConfig := `{
    "Host": "http://10.41.1.242:6000",
    "SecretKey": "qizvkzmg92trkqpwi0c8bzw5y6r31bsv",
    "PayNotifyPath": "/notify/xpay/pay",
    "TransferNotifyPath": "/notify/xpay/transfer"
}`
	up, err := NewXPay(upMerchantNo, jsonConfig)
	if err != nil {
		log.Fatalf("获取上游失败, err=%v", err)
		return nil
	}

	return up
}

func TestXPay_Pay(t *testing.T) {
	up := GetXPayConfig()
	req := &PayRequest{
		Amount:       1000,
		Currency:     "CNY",
		OrderNo:      "TPay" + utils.GetDailyId(),
		NotifyUrl:    "http://127.0.0.1:8000/notify/xpay/pay",
		ReturnUrl:    "",
		ProductType:  "UNION",
		Subject:      "TPay测试",
		CustomName:   "",
		CustomMobile: "",
		CustomEmail:  "",
		Attach:       "",
	}
	resp, err := up.Pay(req)
	if err != nil {
		t.Errorf("Pay() error = %v", err)
		return
	}

	t.Logf("resp: %+v", resp)
}

// TPay2021043010420175885521  P161975052138276660980
func TestXPay_PayOrderQuery(t *testing.T) {
	up := GetXPayConfig()
	req := &PayOrderQueryRequest{
		OrderNo:         "",
		UpstreamOrderNo: "P161975052138276660980",
	}
	resp, err := up.PayOrderQuery(req)
	if err != nil {
		t.Errorf("PayOrderQuery() error = %v", err)
		return
	}
	t.Logf("resp: %+v", resp)
}

func TestXPay_Transfer(t *testing.T) {
	up := GetXPayConfig()
	req := &TransferRequest{
		Amount:             1000,
		Currency:           "CNY",
		OrderNo:            "TPay" + utils.GetDailyId(),
		NotifyUrl:          "http://127.0.0.1:8000/notify/xpay/transfer",
		ReturnUrl:          "",
		ProductType:        "UNION",
		Attach:             "",
		Remark:             "",
		BankName:           "TpayTest",
		BankBranchName:     "TpayTest",
		BankCardNo:         "TpayTest",
		BankCode:           "TpayTest",
		BankCardHolderName: "TpayTest",
		CardHolderMobile:   "",
		CardHolderEmail:    "",
	}
	resp, err := up.Transfer(req)
	if err != nil {
		t.Errorf("Transfer() error = %v", err)
		return
	}
	t.Logf("resp: %+v", resp)
}

// TPay2021043010524350180252  T161975116372150515972
func TestXPay_TransferOrderQuery(t *testing.T) {
	up := GetXPayConfig()
	req := &TransferOrderQueryRequest{
		OrderNo:         "",
		UpstreamOrderNo: "T161975116372150515972",
	}
	resp, err := up.TransferOrderQuery(req)
	if err != nil {
		t.Errorf("TransferOrderQuery() error = %v", err)
		return
	}
	t.Logf("resp: %+v", resp)
}

func TestXPay_QueryBalance(t *testing.T) {
	up := GetXPayConfig()
	resp, err := up.QueryBalance()
	if err != nil {
		t.Errorf("QueryBalance() error = %v", err)
		return
	}
	t.Logf("resp: %+v", resp)
}
