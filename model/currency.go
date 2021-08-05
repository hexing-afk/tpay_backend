package model

import "gorm.io/gorm"

const CurrencyTableName = "currency"

const (
	CurrencyUSD = "USD"
	CurrencyKHR = "KHR"
	CurrencyCNY = "CNY"
	CurrencyTHB = "THB"
	CurrencyINR = "INR"

	// 币种对应的金额是否需要除以100显示
	DivideHundred    = 1 // 是
	NotDivideHundred = 2 // 否
)

type Currency struct {
	Id              int64  `gorm:"id"`
	Currency        string `gorm:"currency"`          // 币种
	Symbol          string `gorm:"symbol"`            // 符号
	Country         string `gorm:"country"`           // 国家
	IsDivideHundred int64  `gorm:"is_divide_hundred"` // 前端是否需要除以100
}

func (t *Currency) TableName() string {
	return CurrencyTableName
}

func NewCurrencyModel(db *gorm.DB) *CurrencyModel {
	return &CurrencyModel{db: db}
}

type CurrencyModel struct {
	db *gorm.DB
}

func (m *CurrencyModel) FindMany() ([]*Currency, error) {
	var list []*Currency
	err := m.db.Table(CurrencyTableName).Scan(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (m *CurrencyModel) FindByCurrency(currency string) (*Currency, error) {
	var data Currency
	err := m.db.Table(CurrencyTableName).Where("currency = ? ", currency).Limit(1).Scan(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}
