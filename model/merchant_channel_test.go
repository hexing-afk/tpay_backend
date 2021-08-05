package model

import (
	"testing"
	"tpay_backend/test"
)

func TestMerchantChannelModel_FindPlfMchTransferChannel(t *testing.T) {
	m := NewMerchantChannelModel(test.DbEngine)
	var merchantId int64 = 36
	currency := "USD"

	ret, err := m.FindPlfMchTransferChannel(merchantId, 0, currency)
	if err != nil {
		t.Errorf("FindPlfMchTransferChannel() error = %v", err)
		return
	}

	t.Logf("ret: %+v", ret)
}

func TestMerchantChannelModel_FindOneByPlatformId(t *testing.T) {
	m := NewMerchantChannelModel(test.DbEngine)
	var platformId int64 = 11
	ids, err := m.FindOneByPlatformId(platformId)
	if err != nil {
		t.Errorf("FindOneByPlatformId() error = %v", err)
		return
	}
	t.Logf("ids: %+v", ids)
}
