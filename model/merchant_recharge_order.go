package model

import (
	"bytes"
	"fmt"
	"gorm.io/gorm"
	"time"
)

const MerchantRechargeOrderTableName = "merchant_recharge_order"

const (
	// 订单状态
	RechargeOrderStatusPending = 1 // 待处理
	RechargeOrderStatusPass    = 2 // 通过
	RechargeOrderStatusReject  = 3 // 驳回
)

type MerchantRechargeOrder struct {
	Id                 int64  `gorm:"id"`                    // 商户充值记录id
	OrderNo            string `gorm:"order_no"`              // 订单号
	OrderAmount        int64  `gorm:"order_amount"`          // 订单金额
	MerchantId         int64  `gorm:"merchant_id"`           // 商户id
	OrderStatus        int64  `gorm:"order_status"`          // 订单状态:1-待处理；2-处理中; 3-通过； 4-驳回
	CreateTime         int64  `gorm:"create_time"`           // 创建时间
	RechargeRemark     string `gorm:"recharge_remark"`       // 充值备注
	PlatformBankCardId int64  `gorm:"platform_bank_card_id"` // 平台收款卡id
	FinishTime         int64  `gorm:"finish_time"`           // 处理完成时间
	AuditRemark        string `gorm:"audit_remark"`          // 审核备注
	BankName           string `gorm:"bank_name"`             // 收款银行
	CardNumber         string `gorm:"card_number"`           // 收款卡号
	PayeeName          string `gorm:"payee_name"`            // 收款人姓名
	BranchName         string `gorm:"branch_name"`           // 支行名称
	Currency           string `gorm:"currency"`              // 币种
}

func (t *MerchantRechargeOrder) TableName() string {
	return MerchantRechargeOrderTableName
}

func NewMerchantRechargeOrderModel(db *gorm.DB) *MerchantRechargeOrderModel {
	return &MerchantRechargeOrderModel{db: db}
}

type MerchantRechargeOrderModel struct {
	db *gorm.DB
}

// 插入一条记录
func (m *MerchantRechargeOrderModel) Insert(data *MerchantRechargeOrder) error {
	data.CreateTime = time.Now().Unix()
	result := m.db.Create(data)
	return result.Error
}

// 根据订单号查询一条记录
func (m *MerchantRechargeOrderModel) FindByOrderNo(orderNo string) (*MerchantRechargeOrder, error) {
	var order = &MerchantRechargeOrder{}
	result := m.db.Model(order).Where("order_no = ?", orderNo).First(order)
	return order, result.Error
}

type FindRechargeOrderList struct {
	Page            int64
	PageSize        int64
	MerchantId      int64
	MerchantName    string
	OrderNo         string
	OrderStatus     int64
	StartCreateTime int64
	EndCreateTime   int64
}

type RechargeOrderList struct {
	MerchantRechargeOrder
	MerchantName string `gorm:"merchant_name"`
}

// 查询列表
func (m *MerchantRechargeOrderModel) FindList(f FindRechargeOrderList) ([]*RechargeOrderList, int64, error) {
	var (
		order    = fmt.Sprintf("%s o", MerchantRechargeOrderTableName)
		merchant = fmt.Sprintf("left join %s m on m.id = o.merchant_id ", MerchantTableName)

		selectField = "o.order_no, o.order_amount, o.recharge_remark, o.bank_name, o.payee_name, " +
			"o.card_number, o.order_status, o.audit_remark, o.create_time, o.finish_time, " +
			"o.branch_name, o.currency, m.username AS merchant_name "

		whereBuffer = bytes.NewBufferString(" 1=1 ")
		args        []interface{}
	)

	if f.MerchantId != 0 {
		whereBuffer.WriteString("and m.id = ? ")
		args = append(args, f.MerchantId)
	}

	if f.MerchantName != "" {
		whereBuffer.WriteString("and m.username like ? ")
		args = append(args, "%"+f.MerchantName+"%")
	}

	if f.OrderNo != "" {
		whereBuffer.WriteString("and o.order_no = ? ")
		args = append(args, "%"+f.OrderNo+"%")
	}

	if f.OrderStatus != 0 {
		whereBuffer.WriteString("and o.order_status = ? ")
		args = append(args, f.OrderStatus)
	}

	if f.StartCreateTime != 0 {
		whereBuffer.WriteString("and o.create_time >= ? ")
		args = append(args, f.StartCreateTime)
	}

	if f.EndCreateTime != 0 {
		whereBuffer.WriteString("and o.create_time <= ? ")
		args = append(args, f.EndCreateTime)
	}

	if f.Page == 0 {
		f.Page = 1
	}
	if f.PageSize == 0 {
		f.PageSize = 10
	}

	var total int64
	if ret := m.db.Table(order).Joins(merchant).Where(whereBuffer.String(), args...).Count(&total); ret.Error != nil {
		return nil, 0, ret.Error
	}

	if total == 0 {
		return nil, 0, nil
	}

	whereBuffer.WriteString("order by o.create_time desc limit ? offset ? ")
	args = append(args, f.PageSize, (f.Page-1)*f.PageSize)

	var list []*RechargeOrderList
	if ret := m.db.Table(order).Select(selectField).Joins(merchant).Where(whereBuffer.String(), args...).Scan(&list); ret.Error != nil {
		return nil, 0, ret.Error
	}

	return list, total, nil
}

// 修改订单状态-通过
func (m *MerchantRechargeOrderModel) UpdateStatusToPass(orderNo, remark string) error {
	setMap := map[string]interface{}{
		"finish_time":  time.Now().Unix(),
		"order_status": RechargeOrderStatusPass,
	}

	if remark != "" {
		setMap["audit_remark"] = remark
	}

	result := m.db.Model(&MerchantRechargeOrder{}).Where("order_no = ? and order_status= ? ", orderNo, RechargeOrderStatusPending).Updates(setMap)
	return result.Error
}

// 修改订单状态-驳回
func (m *MerchantRechargeOrderModel) UpdateStatusToReject(orderNo, remark string) error {
	setMap := map[string]interface{}{
		"finish_time":  time.Now().Unix(),
		"order_status": RechargeOrderStatusReject,
	}

	if remark != "" {
		setMap["audit_remark"] = remark
	}

	result := m.db.Model(&MerchantRechargeOrder{}).
		Where("order_no = ? and order_status= ? ", orderNo, RechargeOrderStatusPending).Updates(setMap)
	return result.Error
}
