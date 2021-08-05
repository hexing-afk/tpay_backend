package model

import (
	"fmt"
	"testing"
	"tpay_backend/test"
	"tpay_backend/utils"
)

func TestTransferOrderModel_FindByOrderNo(t *testing.T) {
	transferOrderModel := NewTransferOrderModel(test.DbEngine)

	orderNo := "12344"
	order, err := transferOrderModel.FindByOrderNo(orderNo)
	if err != nil {
		t.Errorf("FindByOrderNo() error = %v", err)
		return
	}

	t.Logf("订单：%+v", order)
}

func TestTransferOrderModel_FindByMerchantId(t *testing.T) {
	m := NewTransferOrderModel(test.DbEngine)

	var merchantId int64 = 47
	orderNo := "8161889113632589723526"

	order, err := m.FindByMerchantId(merchantId, orderNo)
	if err != nil {
		t.Errorf("FindByMerchantId() error = %v", err)
		return
	}

	t.Logf("结果：%+v", order)
}

func TestTransferOrderModel_FindExportData(t *testing.T) {
	m := NewTransferOrderModel(test.DbEngine)
	f := FindTransferExportData{
		OrderNo:           "",
		MerchantOrderNo:   "",
		UpstreamOrderNo:   "",
		MerchantName:      "",
		Currency:          "",
		PlatformChannelId: 0,
		StartCreateTime:   0,
		EndCreateTime:     0,
		OrderStatus:       0,
		OrderSourceList:   []int64{1},
		MerchantNo:        "16191491717906143664",
		OrderType:         "test",
	}

	data, err := m.FindExportData(f)
	if err != nil {
		t.Errorf("FindExportData() error = %v", err)
		return
	}
	t.Logf("data：%+v", data)
}

func TestTransferOrderModel_Insert(t *testing.T) {
	for i := 0; i < 10000; i++ {
		xxxx()
	}
}

func xxxx() {
	order := new(TransferOrder)

	order.OrderNo = order.GenerateOrderNo(TransferModeTest)
	order.MerchantNo = "16190890594704021927"
	order.MerchantOrderNo = "MM" + utils.GetUniqueId()
	order.ReqAmount = utils.RandInt64(100, 1000000)
	order.Currency = "VND"

	order.OrderStatus = TransferOrderStatusPending
	order.NotifyStatus = TransferNotifyStatusNot
	order.OrderSource = TransferOrderSourceInterface
	order.BankName = "中国银行"
	order.AccountName = "张三"
	order.CardNumber = "6342340234923888"
	order.BranchName = "上海分行"
	order.Remark = "测试"
	order.AreaId = 3
	order.Mode = TransferModePro

	order.MerchantRate = 3.5
	order.MerchantSingleFee = 10

	if err := NewTransferOrderModel(test.DbEngine).Insert(order); err != nil {

	}

	fmt.Println("插入成功id=", order.Id)
}
