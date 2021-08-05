package model

import "gorm.io/gorm"

const BankTableName = "bank"

type Bank struct {
	Id           int64  `gorm:"id"`
	BankName     string `gorm:"Bank_name"`     // 银行名
	SingleAmount string `gorm:"single_amount"` // 单笔金额
}

func (t *Bank) TableName() string {
	return BankTableName
}

func NewBankModel(db *gorm.DB) *BankModel {
	return &BankModel{db: db}
}

type BankModel struct {
	db *gorm.DB
}

func (m *BankModel) FindMany() ([]*Bank, error) {
	var list []*Bank
	err := m.db.Table(BankTableName).Scan(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (m *BankModel) FindManyLimit(amount int64) ([]*Bank, error) {
	var list []*Bank
	err := m.db.Table(BankTableName).
		Where(`single_amount >= ?`, amount).
		Scan(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (m *BankModel) FindByBank(bank string) (*Bank, error) {
	var data Bank
	err := m.db.Table(BankTableName).Where("bank = ? ", bank).Limit(1).Scan(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}
