package model

import (
	"gorm.io/gorm"
	"testing"
	"time"
	"tpay_backend/test"
)

func TestMerchantModel_Save(t *testing.T) {
	model := NewMerchantModel(test.DbEngine)

	merchantId := int64(39)
	merchant, err := model.FindOneById(merchantId)

	if err != nil {
		t.Logf("findErr:%v", err)
		return
	}

	merchant.IpWhiteList = "aaa,bbb "

	saveErr := model.Save(merchant)

	t.Logf("saveErr:%v", saveErr)
}

func TestMerchantModel_MinusBalanceFreezeTx(t *testing.T) {
	var merchantId int64 = 41
	var amount int64 = 100
	walletLogExt := WalletLogExt{
		BusinessNo: "test001",
		Source:     99,
		Remark:     "事务测试",
	}

	for i := 1; i <= 10; i++ {
		go func() {
			txErr := test.DbEngine.Transaction(func(tx *gorm.DB) error {
				if err := NewMerchantModel(tx).MinusBalanceFreezeTx(merchantId, amount, walletLogExt); err != nil {
					t.Errorf("MinusBalanceFreezeTx() error = %v", err)
					return err
				}

				return nil
			})

			if txErr != nil {
				t.Errorf("error: %v", txErr)
				return
			}

		}()
		time.Sleep(5 * time.Second)
	}

	t.Log("成功")
}
