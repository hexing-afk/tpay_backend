package model

import (
	"bytes"
	"fmt"
	"time"

	"gorm.io/gorm"
)

const MerchantWalletLogTableName = "merchant_wallet_log"

const (
	//金额变动类型
	OpTypeAddBalance   = 1 // 加余额
	OpTypeMinusBalance = 2 // 减余额
	OpTypeAddFrozen    = 3 // 加冻结金额
	OpTypeMinusFrozen  = 4 // 减冻结金额

	//变动金额来源
	AmountSourcePlatform   = 1 // 手动调账
	AmountSourceRecharge   = 2 // 充值
	AmountSourceWithdraw   = 3 // 提现
	AmountSourceTransfer   = 4 // 代付
	AmountSourceCollection = 5 // 收款
)

type MerchantWalletLog struct {
	Id           int64  `gorm:"id"`            // 日志号
	MerchantId   int64  `gorm:"merchant_id"`   // 商户id
	CreateTime   int64  `gorm:"create_time"`   // 创建时间
	OpType       int64  `gorm:"op_type"`       // 变动类型：1+，2-
	ChangeAmount int64  `gorm:"change_amount"` // 变动金额
	AfterBalance int64  `gorm:"after_balance"` // 变动后余额
	BusinessNo   string `gorm:"business_no"`   // 业务单号
	Source       int64  `gorm:"source"`        // 变动来源：1-手动调账
	Remark       string `gorm:"remark"`        // 备注
}

func (t *MerchantWalletLog) TableName() string {
	return MerchantWalletLogTableName
}

func NewMerchantWalletLogModel(db *gorm.DB) *MerchantWalletLogModel {
	return &MerchantWalletLogModel{db: db}
}

type MerchantWalletLogModel struct {
	db *gorm.DB
}

// 插入一条记录
func (m *MerchantWalletLogModel) Insert(o *MerchantWalletLog) error {
	o.CreateTime = time.Now().Unix()
	result := m.db.Create(o)
	return result.Error
}

type FindMerchantWalletLogList struct {
	Page            int64
	PageSize        int64
	Id              int64
	BusinessNo      string
	Username        string
	OpType          int64
	Source          int64
	Currency        string
	StartCreateTime int64
	EndCreateTime   int64
	IdOrBusinessNo  string
	MerchantId      int64
	OpTypeList      []int64
}

type MerchantWalletLogList struct {
	Username string `gorm:"username"` // 商户账号
	Currency string `gorm:"currency"` // 商户币种
	MerchantWalletLog
}

// 查询列表
func (m *MerchantWalletLogModel) FindList(f FindMerchantWalletLogList) ([]MerchantWalletLogList, int64, error) {
	var (
		merchantWalletLog = fmt.Sprintf("%s log ", MerchantWalletLogTableName)
		merchant          = fmt.Sprintf("left join %s m on m.id = log.merchant_id ", MerchantTableName)

		selectField = "log.id, log.business_no, log.change_amount, log.after_balance, log.op_type, " +
			"log.source, log.remark, log.create_time, m.username, m.currency "

		whereBuffer = bytes.NewBufferString(" 1=1 ")
		args        []interface{}
	)

	if f.Id != 0 {
		whereBuffer.WriteString("and log.id = ? ")
		args = append(args, f.Id)
	}
	if f.MerchantId != 0 {
		whereBuffer.WriteString("and log.merchant_id = ? ")
		args = append(args, f.MerchantId)
	}

	if f.IdOrBusinessNo != "" {
		whereBuffer.WriteString("and ( log.id like ? or log.business_no like ? ) ")
		args = append(args, "%"+f.IdOrBusinessNo+"%", "%"+f.IdOrBusinessNo+"%")
	}

	if f.BusinessNo != "" {
		whereBuffer.WriteString("and log.business_no = ? ")
		args = append(args, f.BusinessNo)
	}

	if f.Username != "" {
		whereBuffer.WriteString("and m.username like ? ")
		args = append(args, "%"+f.Username+"%")
	}

	if len(f.OpTypeList) > 0 {
		whereBuffer.WriteString("and log.op_type in ? ")
		args = append(args, f.OpTypeList)
	}

	if f.OpType != 0 {
		whereBuffer.WriteString("and log.op_type = ? ")
		args = append(args, f.OpType)
	}

	if f.Source != 0 {
		whereBuffer.WriteString("and log.source = ? ")
		args = append(args, f.Source)
	}

	if f.Currency != "" {
		whereBuffer.WriteString("and m.currency = ? ")
		args = append(args, f.Currency)
	}

	if f.StartCreateTime != 0 {
		whereBuffer.WriteString("and log.create_time >= ? ")
		args = append(args, f.StartCreateTime)
	}

	if f.EndCreateTime != 0 {
		whereBuffer.WriteString("and log.create_time <= ? ")
		args = append(args, f.EndCreateTime)
	}

	if f.Page == 0 {
		f.Page = 1
	}
	if f.PageSize == 0 {
		f.PageSize = 10
	}

	var total int64
	result := m.db.Table(merchantWalletLog).Joins(merchant).Where(whereBuffer.String(), args...).Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	if total == 0 {
		return nil, 0, nil
	}

	whereBuffer.WriteString("order by log.create_time desc limit ? offset ? ")
	args = append(args, f.PageSize, (f.Page-1)*f.PageSize)

	var list []MerchantWalletLogList
	result = m.db.Table(merchantWalletLog).Select(selectField).Joins(merchant).Where(whereBuffer.String(), args...).Scan(&list)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return list, total, nil
}

type FindWalletExportData struct {
	Id              int64
	BusinessNo      string
	Username        string
	OpType          int64
	Source          int64
	Currency        string
	StartCreateTime int64
	EndCreateTime   int64
	IdOrBusinessNo  string
	MerchantId      int64
	OpTypeList      []int64
}

type WalletExportData struct {
	Total int64
	List  []MerchantWalletLogList
}

// 查询导出数据
func (m *MerchantWalletLogModel) FindExportData(f FindWalletExportData) (*WalletExportData, error) {
	var (
		merchantWalletLog = fmt.Sprintf("%s log ", MerchantWalletLogTableName)
		merchant          = fmt.Sprintf("left join %s m on m.id = log.merchant_id ", MerchantTableName)

		selectField = "log.id, log.business_no, log.change_amount, log.after_balance, log.op_type, " +
			"log.source, log.remark, log.create_time, m.currency "

		whereBuffer = bytes.NewBufferString(" 1=1 ")
		args        []interface{}
	)

	if f.Id != 0 {
		whereBuffer.WriteString("and log.id = ? ")
		args = append(args, f.Id)
	}
	if f.MerchantId != 0 {
		whereBuffer.WriteString("and log.merchant_id = ? ")
		args = append(args, f.MerchantId)
	}

	if f.IdOrBusinessNo != "" {
		whereBuffer.WriteString("and ( log.id like ? or log.business_no like ? ) ")
		args = append(args, "%"+f.IdOrBusinessNo+"%", "%"+f.IdOrBusinessNo+"%")
	}

	if f.BusinessNo != "" {
		whereBuffer.WriteString("and log.business_no = ? ")
		args = append(args, f.BusinessNo)
	}

	if f.Username != "" {
		whereBuffer.WriteString("and m.username like ? ")
		args = append(args, "%"+f.Username+"%")
	}

	if f.OpType != 0 {
		whereBuffer.WriteString("and log.op_type = ? ")
		args = append(args, f.OpType)
	}

	if len(f.OpTypeList) > 0 {
		whereBuffer.WriteString("and log.op_type in ? ")
		args = append(args, f.OpTypeList)
	}

	if f.Source != 0 {
		whereBuffer.WriteString("and log.source = ? ")
		args = append(args, f.Source)
	}

	if f.Currency != "" {
		whereBuffer.WriteString("and m.currency = ? ")
		args = append(args, f.Currency)
	}

	if f.StartCreateTime != 0 {
		whereBuffer.WriteString("and log.create_time >= ? ")
		args = append(args, f.StartCreateTime)
	}

	if f.EndCreateTime != 0 {
		whereBuffer.WriteString("and log.create_time <= ? ")
		args = append(args, f.EndCreateTime)
	}

	var total int64
	result := m.db.Table(merchantWalletLog).Joins(merchant).Where(whereBuffer.String(), args...).Count(&total)
	if result.Error != nil {
		return nil, result.Error
	}

	if total == 0 {
		return nil, nil
	}

	whereBuffer.WriteString("order by log.create_time asc ")

	var list []MerchantWalletLogList
	result = m.db.Table(merchantWalletLog).Select(selectField).Joins(merchant).Where(whereBuffer.String(), args...).Scan(&list)
	if result.Error != nil {
		return nil, result.Error
	}

	return &WalletExportData{
		Total: total,
		List:  list,
	}, nil
}
