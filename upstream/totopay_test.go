package upstream

import (
	"fmt"
	"log"
	"testing"
	"tpay_backend/utils"
)

func GetTotopayObj() Upstream {
	upMerchantNo := "3d2f9bb9-f370-412e-91fb-d011f76706f3"
	jsonStr := `{"Host":"http://10.41.1.242:8888/","SecretKey":"107efadfbaca01b2c8d577df51534bf5","PayNotifyPath":"/notify/totopay/pay","TransferNotifyPath":"/notify/totopay/transfer"}`

	obj, err := NewTotopay(upMerchantNo, jsonStr)
	if err != nil {
		log.Fatalf("获取对象失败:%v", err)
	}

	return obj
}

func TestTotopay_Pay(t *testing.T) {
	obj := GetTotopayObj()

	req := &PayRequest{
		Amount:      300,
		Currency:    "USD",
		OrderNo:     fmt.Sprintf("P%d", utils.RandInt64(100000, 999999)),
		NotifyUrl:   "http://www.abc.com/notify",
		ReturnUrl:   "http://www.abc.com/return",
		ProductType: "xx_unionpay",
	}

	res, err := obj.Pay(req)

	t.Logf("err:%v", err)
	t.Logf("res:%v", res)
}

func TestTotopay_PayOrderQuery(t *testing.T) {
	obj := GetTotopayObj()

	req := &PayOrderQueryRequest{
		OrderNo: "TEST263133",
	}

	res, err := obj.PayOrderQuery(req)

	t.Logf("err:%v", err)
	t.Logf("res:%v", res)
}

func TestTotopay_Transfer(t *testing.T) {
	obj := GetTotopayObj()

	req := &TransferRequest{
		Amount:             100,
		Currency:           "USD",
		OrderNo:            fmt.Sprintf("T%d", utils.RandInt64(100000, 999999)),
		NotifyUrl:          "http://www.xxx.com/totopay/notify",
		ProductType:        "withdraw",
		BankName:           "中国建设银行",
		BankCardNo:         "62222444555888",
		BankCardHolderName: "李四",
	}

	res, err := obj.Transfer(req)

	t.Logf("err:%v", err)
	t.Logf("res:%v", res)
}

func TestTotopay_TransferOrderQuery(t *testing.T) {
	obj := GetTotopayObj()

	req := &TransferOrderQueryRequest{
		OrderNo: "T143211",
	}

	res, err := obj.TransferOrderQuery(req)

	t.Logf("err:%v", err)
	t.Logf("res:%v", res)
}

func TestTotopay_QueryBalance(t *testing.T) {
	obj := GetTotopayObj()

	res, err := obj.QueryBalance()

	t.Logf("err:%v", err)
	t.Logf("res:%v", res)
}
