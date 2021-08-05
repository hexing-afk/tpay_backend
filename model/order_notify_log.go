package model

import (
	"gorm.io/gorm"
	"time"
)

const OrderNotifyLogTableName = "order_notify_log"

const (
	// 通知状态
	OrderNotifyStatusNotifying = 0 // 通知中
	OrderNotifyStatusSuccess   = 1 // 成功
	OrderNotifyStatusFail      = 2 // 失败

	// 通知日志订单类型
	NotifyLogOrderTypePay      = 1 // 代收订单
	NotifyLogOrderTypeTransfer = 2 // 代付订单
)

type OrderNotifyLog struct {
	Id              int64  `db:"id"`
	OrderNo         string `db:"order_no"`          // 内部订单号
	MerchantOrderNo string `db:"merchant_order_no"` // 商户订单号
	Status          int64  `db:"status"`            // 通知状态(0通知中,1成功,2失败)
	Result          string `db:"result"`            // 通知返回的结果
	CreateTime      int64  `db:"create_time"`       // 通知时间
	OrderType       int64  `db:"order_type"`        // 订单类型(1代收订单, 2代付订单)
}

func (t *OrderNotifyLog) TableName() string {
	return OrderNotifyLogTableName
}

func NewOrderNotifyLogModel(db *gorm.DB) *OrderNotifyLogModel {
	return &OrderNotifyLogModel{db: db}
}

type OrderNotifyLogModel struct {
	db *gorm.DB
}

func (m *OrderNotifyLogModel) Insert(data *OrderNotifyLog) error {
	data.CreateTime = time.Now().Unix()
	result := m.db.Create(data)
	return result.Error
}

func (m *OrderNotifyLogModel) Update(id int64, data OrderNotifyLog) error {
	result := m.db.Model(&OrderNotifyLog{}).Where("id=?", id).Updates(&data)
	return result.Error
}
