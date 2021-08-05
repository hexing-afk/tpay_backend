package order_notify

import (
	"strings"
	"testing"
	"tpay_backend/utils"
)

func TestGetPayNotifyExpireKey(t *testing.T) {
	orderNo := "7161839048746535761091"
	t.Logf("key: %s", GetPayNotifyExpireKey(orderNo))
}

func TestGetTransferNotifyExpireKey(t *testing.T) {
	orderNo := "8161839310558589810826"
	t.Logf("key: %s", GetTransferNotifyExpireKey(orderNo))
}

func TestSing(t *testing.T) {
	data := `{"amount":5000000,"currency":"VND","mch_order_no":"20210609000000002243","merchant_no":"16231367335980509487","notify_url":"http://192.168.50.33:10960/paymentresponse/notify/hxpay/merchant19","return_url":"http://testh5.mc900.com/MCkey/mine?fromPayer=true","subject":"1","timestamp":1623222465,"trade_type":"ACBBANK"}AITYNWOMEM6XLX335FGWLQIFKX6STIWW`
	sign := strings.ToLower(utils.Md5(data))
	t.Logf("sign=%v", sign)
}
