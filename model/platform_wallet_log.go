package model

import (
	"bytes"
	"gorm.io/gorm"
	"time"
)

const PlatformWalletLogTableName = "platform_wallet_log"

const (
	PlatformIncomeSourceWithdraw = 3 // 商户提现手续费
	PlatformIncomeSourceTransfer = 4 // 商户代付手续费
	PlatformIncomeSourcePay      = 5 // 商户代收手续费
)

type PlatformWalletLog struct {
	Id          int64  `gorm:"id"`           // 平台钱包日志id
	BusinessNo  string `gorm:"business_no"`  // 业务号
	Source      int64  `gorm:"source"`       // 收益来源:3-商户提现；4-商户代付；5-商户代收；
	MerchantFee int64  `gorm:"merchant_fee"` // 商户手续费
	UpstreamFee int64  `gorm:"upstream_fee"` // 上游手续费
	Income      int64  `gorm:"income"`       // 收益
	CreateTime  int64  `gorm:"create_time"`  // 创建时间
	Currency    string `gorm:"currency"`     // 币种
}

func (t *PlatformWalletLog) TableName() string {
	return PlatformWalletLogTableName
}

type PlatformWalletLogModel struct {
	db *gorm.DB
}

func NewPlatformWalletLogModel(db *gorm.DB) *PlatformWalletLogModel {
	return &PlatformWalletLogModel{db: db}
}

func (m *PlatformWalletLogModel) Insert(data *PlatformWalletLog) error {
	data.CreateTime = time.Now().Unix()
	result := m.db.Create(data)
	return result.Error
}

type PlatformWalletLogListReq struct {
	BusinessNo string
	Source     int64
	Currency   string
	Page       int64
	PageSize   int64
}

func (m *PlatformWalletLogModel) FindList(req PlatformWalletLogListReq) ([]PlatformWalletLog, int64, error) {

	whereBuffer := bytes.NewBufferString(" 1=1 ")
	var args []interface{}

	if req.BusinessNo != "" {
		whereBuffer.WriteString("and business_no like ? ")
		args = append(args, "%"+req.BusinessNo+"%")
	}

	if req.Source != 0 {
		whereBuffer.WriteString("and source = ? ")
		args = append(args, req.Source)
	}

	if req.Currency != "" {
		whereBuffer.WriteString("and currency = ? ")
		args = append(args, req.Currency)
	}

	var total int64
	if err := m.db.Table(PlatformWalletLogTableName).Where(whereBuffer.String(), args...).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return nil, 0, nil
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	whereBuffer.WriteString(" ORDER BY create_time DESC limit ? offset ?")
	args = append(args, req.PageSize, (req.Page-1)*req.PageSize)

	var list []PlatformWalletLog
	if err := m.db.Table(PlatformWalletLogTableName).Where(whereBuffer.String(), args...).Find(&list).Error; err != nil {
		return nil, total, err
	}
	return list, total, nil
}
