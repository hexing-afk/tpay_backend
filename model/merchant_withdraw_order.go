package model

import (
	"bytes"
	"fmt"
	"gorm.io/gorm"
	"time"
)

const MerchantWithdrawOrderTableName = "merchant_withdraw_order"

const (
	WithdrawOrderStatusPending      = 1 // 待处理
	WithdrawOrderStatusReject       = 2 // 驳回
	WithdrawOrderStatusPass         = 3 // 通过审核
	WithdrawOrderStatusAllot        = 4 // 派单中
	WithdrawOrderStatusAllotSuccess = 5 // 派单成功
	WithdrawOrderStatusAllotFail    = 6 // 派单失败
	WithdrawOrderStatusSuccess      = 7 // 成功
)

type MerchantWithdrawOrder struct {
	Id              int64  `gorm:"id"`                // 商户提现订单id
	OrderNo         string `gorm:"order_no"`          // 订单号
	MerchantId      int64  `gorm:"merchant_id"`       // 商户Id
	OrderAmount     int64  `gorm:"order_amount"`      // 订单金额
	MerchantFee     int64  `gorm:"merchant_fee"`      // 商户手续费
	RealAmount      int64  `gorm:"real_amount"`       // 实际到账金额
	DecreaseAmount  int64  `gorm:"decrease_amount"`   // 扣减商户账户金额
	Remark          string `gorm:"remark"`            // 备注
	BankName        string `gorm:"bank_name"`         // 收款银行
	BranchName      string `gorm:"branch_name"`       // 支行名称
	PayeeName       string `gorm:"payee_name"`        // 收款人
	CardNumber      string `gorm:"card_number"`       // 收款卡号
	CreateTime      int64  `gorm:"create_time"`       // 创建时间
	AuditTime       int64  `gorm:"audit_time"`        // 审核时间
	AuditRemark     string `gorm:"audit_remark"`      // 审核备注
	OrderStatus     int64  `gorm:"order_status"`      // 订单状态: 订单状态: 1-待处理; 2-驳回; 3-通过审核； 4-派单中； 5-派单成功； 6-派单失败；7-成功
	DeductionMethod int64  `gorm:"deduction_method"`  // 商户提现手续费扣款方式
	SuccessRemark   string `gorm:"success_remark"`    // 提现成功备注
	Currency        string `gorm:"currency"`          // 币种
	TransferOrderNo string `gorm:"transfer_order_no"` // 代付单号(派单单号)
	BankCode        string `gorm:"bank_code"`         // 银行代码
	AreaId          int64  `gorm:"area_id"`
}

func (t *MerchantWithdrawOrder) TableName() string {
	return MerchantWithdrawOrderTableName
}

func NewMerchantWithdrawOrderModel(db *gorm.DB) *MerchantWithdrawOrderModel {
	return &MerchantWithdrawOrderModel{db: db}
}

type MerchantWithdrawOrderModel struct {
	db *gorm.DB
}

func (m *MerchantWithdrawOrderModel) Insert(data *MerchantWithdrawOrder) error {
	data.CreateTime = time.Now().Unix()
	result := m.db.Create(data)
	return result.Error
}

func (m *MerchantWithdrawOrderModel) FindByOrderNo(orderNo string) (*MerchantWithdrawOrder, error) {
	var o = &MerchantWithdrawOrder{}
	result := m.db.Model(o).Where("order_no = ?", orderNo).First(o)

	return o, result.Error
}

func (m *MerchantWithdrawOrderModel) FindByTransferOrderNo(transferOrderNo string) (*MerchantWithdrawOrder, error) {
	var o = &MerchantWithdrawOrder{}
	result := m.db.Model(o).Where("transfer_order_no = ?", transferOrderNo).First(o)

	return o, result.Error
}

type FindWithdrawOrderList struct {
	Page            int64
	PageSize        int64
	StartCreateTime int64
	EndCreateTime   int64
	MerchantId      int64
	MerchantName    string
	OrderNo         string
	OrderStatus     int64
	OrderStatuses   []int64
}

type WithdrawOrderList struct {
	MerchantWithdrawOrder
	MerchantName string `gorm:"merchant_name"` // 商户名称
	ChannelName  string `gorm:"channel_name"`  // 上游通道名称
}

// 查询列表
func (m *MerchantWithdrawOrderModel) FindList(f FindWithdrawOrderList) ([]*WithdrawOrderList, int64, error) {
	var (
		order     = fmt.Sprintf("%s o", MerchantWithdrawOrderTableName)
		merchant  = fmt.Sprintf("left join %s m on m.id = o.merchant_id", MerchantTableName)
		transfer  = fmt.Sprintf("left join %s t on t.order_no = o.transfer_order_no", TransferOrderTableName)
		upChannel = fmt.Sprintf("left join %s u on u.id = t.upstream_channel_id", UpstreamChannelTableName)

		selectField = "o.order_no, o.order_amount, o.merchant_fee, o.real_amount, o.remark, " +
			"o.bank_name, o.payee_name, o.card_number, o.branch_name, o.audit_remark, " +
			"o.order_status, o.transfer_order_no, o.create_time, o.audit_time, o.currency, " +
			"o.deduction_method, m.username as merchant_name, u.channel_name"

		whereBuffer = bytes.NewBufferString(" 1=1 ")
		args        []interface{}
	)

	if f.MerchantName != "" {
		whereBuffer.WriteString("and m.username like ? ")
		args = append(args, "%"+f.MerchantName+"%")
	}

	if f.MerchantId != 0 {
		whereBuffer.WriteString("and o.merchant_id = ? ")
		args = append(args, f.MerchantId)
	}

	if f.OrderNo != "" {
		whereBuffer.WriteString("and o.order_no like ? ")
		args = append(args, "%"+f.OrderNo+"%")
	}

	if f.OrderStatus != 0 {
		whereBuffer.WriteString("and o.order_status = ? ")
		args = append(args, f.OrderStatus)
	}

	if f.OrderStatuses != nil {
		whereBuffer.WriteString("and o.order_status in (")
		for i := 0; i < len(f.OrderStatuses); i++ {
			if i == len(f.OrderStatuses)-1 {
				whereBuffer.WriteString(fmt.Sprintf("%v) ", f.OrderStatuses[i]))
			} else {
				whereBuffer.WriteString(fmt.Sprintf("%v, ", f.OrderStatuses[i]))
			}
		}
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
	result := m.db.Table(order).Joins(merchant).Joins(transfer).Joins(upChannel).Where(whereBuffer.String(), args...).Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	if total == 0 {
		return nil, 0, nil
	}

	whereBuffer.WriteString("order by o.create_time desc limit ? offset ? ")
	args = append(args, f.PageSize, (f.Page-1)*f.PageSize)

	var list []*WithdrawOrderList
	result = m.db.Table(order).Select(selectField).Joins(transfer).Joins(merchant).Joins(upChannel).Where(whereBuffer.String(), args...).Scan(&list)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return list, total, nil
}

// 修改状态-审核通过
func (m *MerchantWithdrawOrderModel) UpdateStatusToPass(id int64, remark string) error {
	setMap := map[string]interface{}{
		"audit_time":   time.Now().Unix(),
		"order_status": WithdrawOrderStatusPass,
	}
	if remark != "" {
		setMap["audit_remark"] = remark
	}

	result := m.db.Model(&MerchantWithdrawOrder{}).Where("id = ? and order_status = ?", id, WithdrawOrderStatusPending).Updates(setMap)
	return result.Error
}

// 修改状态-驳回
func (m *MerchantWithdrawOrderModel) UpdateStatusToReject(id int64, remark string) error {
	setMap := map[string]interface{}{
		"audit_time":   time.Now().Unix(),
		"order_status": WithdrawOrderStatusReject,
	}
	if remark != "" {
		setMap["audit_remark"] = remark
	}

	result := m.db.Model(&MerchantWithdrawOrder{}).Where("id = ? and order_status = ?", id, WithdrawOrderStatusPending).Updates(setMap)
	return result.Error
}

// 修改状态-提现成功
func (m *MerchantWithdrawOrderModel) UpdateStatusToSuccess(id int64, remark string) error {
	setMap := map[string]interface{}{
		"order_status": WithdrawOrderStatusSuccess,
	}
	if remark != "" {
		setMap["success_remark"] = remark
	}

	result := m.db.Model(&MerchantWithdrawOrder{}).Where("id = ? ", id).Updates(setMap)
	return result.Error
}

// 修改状态-派单中
func (m *MerchantWithdrawOrderModel) UpdateStatusToAllot(id int64, transferOrderNo string) error {
	setMap := map[string]interface{}{
		"order_status":      WithdrawOrderStatusAllot,
		"transfer_order_no": transferOrderNo,
	}
	result := m.db.Model(&MerchantWithdrawOrder{}).Where("id = ? ", id).Updates(&setMap)
	return result.Error
}

// 修改状态-派单成功
func (m *MerchantWithdrawOrderModel) UpdateStatusToAllotSuccess(id int64) error {
	result := m.db.Model(&MerchantWithdrawOrder{}).Where("id = ? ", id).Update("order_status", WithdrawOrderStatusAllotSuccess)
	return result.Error
}

// 修改状态-派单失败
func (m *MerchantWithdrawOrderModel) UpdateStatusToAllotFail(id int64) error {
	result := m.db.Model(&MerchantWithdrawOrder{}).Where("id = ? ", id).Update("order_status", WithdrawOrderStatusAllotFail)
	return result.Error
}

// 修改状态-派单失败
func (m *MerchantWithdrawOrderModel) UpdateTransferNo(id int64, transferOrderNo string) error {
	result := m.db.Model(&MerchantWithdrawOrder{}).Where("id = ? ", id).Update("transfer_order_no", transferOrderNo)
	return result.Error
}
