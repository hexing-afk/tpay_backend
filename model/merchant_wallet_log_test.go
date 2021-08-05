package model

import (
	"testing"
	"tpay_backend/test"
)

func TestMerchantWalletLogModel_FindList(t *testing.T) {
	m := NewMerchantWalletLogModel(test.DbEngine)

	f := FindMerchantWalletLogList{
		Page:       1,
		PageSize:   10,
		MerchantId: 28,
		OpType:     OpTypeAddBalance,
		OpTypeList: []int64{OpTypeAddBalance, OpTypeMinusBalance},
	}
	list, total, err := m.FindList(f)
	if err != nil {
		t.Errorf("FindList() error = %v", err)
		return
	}

	t.Logf("total:%v, data: %+v", total, list)
}

func TestMerchantWalletLogModel_FindExportData(t *testing.T) {
	m := NewMerchantWalletLogModel(test.DbEngine)
	f := FindWalletExportData{
		Username:        "",
		OpType:          0,
		Source:          0,
		Currency:        "",
		StartCreateTime: 0,
		EndCreateTime:   0,
		IdOrBusinessNo:  "45",
		MerchantId:      0,
		OpTypeList:      []int64{OpTypeAddBalance, OpTypeMinusBalance},
	}

	log, err := m.FindExportData(f)
	if err != nil {
		t.Errorf("FindExportData() error = %v", err)
		return
	}
	t.Logf("结果：%+v", log)
}
