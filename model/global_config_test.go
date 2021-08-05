package model

import (
	"testing"
	"tpay_backend/test"
)

func TestGlobalConfigModel_FindValueByKey(t *testing.T) {
	configModel := NewGlobalConfigModel(test.DbEngine)
	key := ConfigTransferTradeTypeSlice

	value, err := configModel.FindValueByKey(key)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	t.Logf("value:%v", value)
}

func TestGlobalConfigModel_TotpIsClose(t *testing.T) {
	m := NewGlobalConfigModel(test.DbEngine)

	isClose, err := m.TotpIsClose()
	if err != nil {
		t.Errorf("TotpIsClose-error: %v", err)
		return
	}
	t.Logf("isClose: %v", isClose)
}
