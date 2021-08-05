package model

import (
	"gorm.io/gorm"
)

const MerchantWithdrawConfigTableName = "merchant_withdraw_config"

const (
	// 手续费扣费方式
	MerchantWithdrawDeductionInner = 1 // 内扣
	MerchantWithdrawDeductionOut   = 2 // 外扣
)

type MerchantWithdrawConfig struct {
	Id              int64   `gorm:"id"`
	MerchantId      int64   `gorm:"merchant_id"`       // 商户id
	SingleMinAmount int64   `gorm:"single_min_amount"` // 单笔提现最小金额
	SingleMaxAmount int64   `gorm:"single_max_amount"` // 单笔提现最大金额
	DeductionMethod int64   `gorm:"deduction_method"`  // 手续费扣费方式: 1-内扣， 2-外扣
	Rate            float64 `gorm:"rate"`              // 提现费率
	SingleFee       int64   `gorm:"single_fee"`        // 单笔提现手续费
}

func (t *MerchantWithdrawConfig) TableName() string {
	return MerchantWithdrawConfigTableName
}

func NewMerchantWithdrawConfigModel(db *gorm.DB) *MerchantWithdrawConfigModel {
	return &MerchantWithdrawConfigModel{db: db}
}

type MerchantWithdrawConfigModel struct {
	db *gorm.DB
}

// 插入一条记录
func (m *MerchantWithdrawConfigModel) Insert(data *MerchantWithdrawConfig) error {
	result := m.db.Create(data)
	return result.Error
}

// 忽略0值|空值, 更新数据
func (m *MerchantWithdrawConfigModel) UpdatePortion(data MerchantWithdrawConfig) error {
	result := m.db.Model(&MerchantWithdrawConfig{}).Updates(data)
	return result.Error
}

// 更新数据, 包括0值和空值
func (m *MerchantWithdrawConfigModel) Update(id int64, data MerchantWithdrawConfig) error {
	setMap := map[string]interface{}{
		"single_min_amount": data.SingleMinAmount,
		"single_max_amount": data.SingleMaxAmount,
		"rate":              data.Rate,
		"single_fee":        data.SingleFee,
	}

	if data.DeductionMethod != 0 {
		setMap["deduction_method"] = data.DeductionMethod
	}

	result := m.db.Model(&MerchantWithdrawConfig{}).Where("id = ?", id).Updates(setMap)
	return result.Error
}

func (m *MerchantWithdrawConfigModel) CheckByMerchantId(merchantId int64) (bool, error) {
	var cnt int64
	result := m.db.Model(&MerchantWithdrawConfig{}).Where("merchant_id = ?", merchantId).Count(&cnt)
	if result.Error != nil {
		return false, result.Error
	}
	return cnt > 0, nil
}

func (m *MerchantWithdrawConfigModel) FindOneById(id int64) (*MerchantWithdrawConfig, error) {
	var data = &MerchantWithdrawConfig{}
	result := m.db.Model(&MerchantWithdrawConfig{}).Where("id = ? ", id).First(data)
	return data, result.Error
}

func (m *MerchantWithdrawConfigModel) FindOneByMerchantId(merchantId int64) (*MerchantWithdrawConfig, error) {
	var data = &MerchantWithdrawConfig{}
	result := m.db.Model(&MerchantWithdrawConfig{}).Where("merchant_id = ?", merchantId).First(data)
	return data, result.Error
}
