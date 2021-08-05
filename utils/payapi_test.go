package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestSendTransfer(t *testing.T) {
	url := "https://tpay-api.mangopay-test.com/system/transfer"
	req := TransferApiRequest{
		MerchantNo:         "",
		Timestamp:          time.Now().Unix(),
		Amount:             6000,
		Currency:           "VND",
		MchOrderNo:         fmt.Sprintf("T%v", time.Now().Unix()),
		TradeType:          "TRANSFER_BANK",
		OrderSource:        0,
		BankName:           "NGAN HANG TMCP A CHAU (ACB)",
		BankCardHolderName: "张三",
		BankCardNo:         "7946566648111",
		BankBranchName:     "",
		BankCode:           "118",
		Remark:             "",
	}
	secretKey := ""

	resp, err := SendTransfer(url, req, secretKey)
	if err != nil {
		t.Errorf("SendTransfer() error = %v", err)
		return
	}
	t.Logf("resp：%+v", resp)
}
