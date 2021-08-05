package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

const MerchantBankCardTableName = "merchant_bank_card"

type MerchantBankCard struct {
	Id          int64  `gorm:"id"`
	MerchantId  int64  `gorm:"merchant_id"`  // 商户id
	CreateTime  int64  `gorm:"create_time"`  // 创建时间
	BankName    string `gorm:"bank_name"`    // 银行名称
	BankCode    string `gorm:"bank_code"`    // 银行代码
	BranchName  string `gorm:"branch_name"`  // 支行名称
	AccountName string `gorm:"account_name"` // 开户名
	CardNumber  string `gorm:"card_number"`  // 银行卡号
	Currency    string `gorm:"currency"`     // 币种
	Remark      string `gorm:"remark"`       // 备注
	UpdateTime  int64  `gorm:"update_time"`  // 更新时间
}

func (t *MerchantBankCard) TableName() string {
	return MerchantBankCardTableName
}

func NewMerchantBankCardModel(db *gorm.DB) *MerchantBankCardModel {
	return &MerchantBankCardModel{db: db}
}

type MerchantBankCardModel struct {
	db *gorm.DB
}

func (m *MerchantBankCardModel) Insert(data *MerchantBankCard) error {
	data.CreateTime = time.Now().Unix()
	return m.db.Create(data).Error
}

func (m *MerchantBankCardModel) Delete(id, merchantId int64) error {
	return m.db.Where("id = ? and merchant_id = ?", id, merchantId).Delete(&MerchantBankCard{}).Error
}

func (m *MerchantBankCardModel) Update(id, merchantId int64, data MerchantBankCard) error {
	setMap := map[string]interface{}{
		"remark":      data.Remark,
		"branch_name": data.BranchName,
		"update_time": time.Now().Unix(),
	}

	if data.BankName != "" {
		setMap["bank_name"] = data.BankName
	}

	if data.AccountName != "" {
		setMap["account_name"] = data.AccountName
	}

	if data.CardNumber != "" {
		setMap["card_number"] = data.CardNumber
	}

	return m.db.Table(MerchantBankCardTableName).Where("id = ? and merchant_id = ?", id, merchantId).Updates(setMap).Error
}

func (m *MerchantBankCardModel) CheckMerchantCard(merchantId int64, bankName, cardNumber string) (bool, error) {
	var cnt int64
	result := m.db.Model(&MerchantBankCard{}).
		Where("merchant_id = ? and bank_name = ? and card_number = ?", merchantId, bankName, cardNumber).Count(&cnt)
	if result.Error != nil {
		return false, result.Error
	}

	return cnt > 0, nil
}

func (m *MerchantBankCardModel) FindByMerchantCardId(cardId, merchantId int64) (*MerchantBankCard, error) {
	var data = &MerchantBankCard{}
	result := m.db.Model(&MerchantBankCard{}).Where("id = ? and merchant_id = ?", cardId, merchantId).First(data)
	return data, result.Error
}

type FindMerchantBankCardList struct {
	Search     string
	MerchantId int64
	Page       int64
	PageSize   int64
}

type MerchantBankCardList struct {
	MerchantUsername string `gorm:"merchant_username"`
	MerchantBankCard
}

func (m *MerchantBankCardModel) FindList(f FindMerchantBankCardList) ([]*MerchantBankCardList, int64, error) {
	var (
		cardTable     = fmt.Sprintf("%s card", MerchantBankCardTableName)
		merchantTable = fmt.Sprintf("left join %s m on m.id = card.merchant_id ", MerchantTableName)

		selectField = "m.username AS merchant_username, card.id, card.bank_name, card.account_name, card.card_number, " +
			"card.branch_name, card.currency, card.remark, card.create_time, card.update_time "

		whereStr = " 1=1 "
		args     []interface{}
	)

	if f.Search != "" {
		whereStr += "and (card.account_name like ? or card.card_number like ? or m.username like ? )"
		args = append(args, "%"+f.Search+"%", "%"+f.Search+"%", "%"+f.Search+"%")
	}

	if f.MerchantId != 0 {
		whereStr += "and card.merchant_id = ? "
		args = append(args, f.MerchantId)
	}

	if f.Page == 0 {
		f.Page = 1
	}
	if f.PageSize == 0 {
		f.PageSize = 10
	}

	var total int64
	err := m.db.Table(cardTable).Joins(merchantTable).Where(whereStr, args...).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return nil, 0, nil
	}

	whereStr += "order by card.create_time desc "
	whereStr += "limit ? offset ? "
	args = append(args, f.PageSize, (f.Page-1)*f.PageSize)

	var resp []*MerchantBankCardList
	err = m.db.Table(cardTable).Select(selectField).Joins(merchantTable).Where(whereStr, args...).Scan(&resp).Error
	if err != nil {
		return nil, 0, err
	}

	return resp, total, nil
}
