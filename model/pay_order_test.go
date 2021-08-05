package model

import (
	"testing"
	"time"
	"tpay_backend/test"
)

func TestPayOrderModel_FindNotifyBreakOrderNo(t *testing.T) {
	m := NewPayOrderModel(test.DbEngine)
	orderStatus := []int64{PayOrderStatusPending, PayOrderStatusPaid, PayOrderStatusFail}
	nextTime := time.Now().Unix()

	orderNos, err := m.FindNotifyBreakOrderNo(orderStatus, nextTime)
	if err != nil {
		t.Errorf("FindNotifyBreakOrderNo() error = %v", err)
		return
	}

	t.Logf("结果: %v", orderNos)
}

func TestPayOrderModel_FindOrderByMerchantId(t *testing.T) {
	m := NewPayOrderModel(test.DbEngine)

	var merchantId int64 = 28
	orderNo := "7161839048746535761091"

	order, err := m.FindOrderByMerchantId(merchantId, orderNo)
	if err != nil {
		t.Errorf("FindOrderByMerchantId() error = %v", err)
		return
	}

	t.Logf("结果: %+v", order)
}

func TestPayOrderModel_FindOneDetailByOrderNo(t *testing.T) {
	m := NewPayOrderModel(test.DbEngine)
	orderNo := "7161847798870707435297"
	order, err := m.FindOneDetailByOrderNo(orderNo)
	if err != nil {
		t.Errorf("FindOneDetailByOrderNo() error = %v", err)
		return
	}
	t.Logf("结果：%+v", order)
}

func TestPayOrderModel_FindExportData(t *testing.T) {
	m := NewPayOrderModel(test.DbEngine)

	f := FindExportData{
		MerchantId:          28,
		MerchantName:        "",
		Currency:            "USD",
		PlatformChannelId:   0,
		StartCreateTime:     0,
		EndCreateTime:       0,
		OrderStatus:         0,
		PlatformChannelName: "",
		OrderType:           "pro",
	}

	got, err := m.FindExportData(f)
	if err != nil {
		t.Errorf("FindExportData() error = %v", err)
		return
	}
	t.Logf("got: %+v", got)
}

func TestPayOrderModel_Insert(t *testing.T) {
	m := NewPayOrderModel(test.DbEngine)

	for i := 1; i < 100000; i++ {
		o := &PayOrder{
			MerchantNo:        "16191491717906143664",
			OrderNo:           "大量数据测试",
			MerchantOrderNo:   "",
			UpstreamOrderNo:   "",
			ReqAmount:         0,
			IncreaseAmount:    0,
			MerchantFee:       0,
			UpstreamAmount:    0,
			UpstreamFee:       0,
			OrderStatus:       0,
			Currency:          "",
			CreateTime:        0,
			UpdateTime:        0,
			NotifyUrl:         "",
			ReturnUrl:         "",
			PlatformChannelId: 0,
			UpstreamChannelId: 0,
			NotifyStatus:      0,
			NotifyFailTimes:   0,
			NextNotifyTime:    0,
			Subject:           "",
			AreaId:            0,
			Mode:              "",
			MerchantRate:      0,
			MerchantSingleFee: 0,
			PaymentAmount:     0,
		}
		if err := m.Insert(o); err != nil {
			t.Errorf("Insert() error = %v", err)
			return
		}
	}
}
