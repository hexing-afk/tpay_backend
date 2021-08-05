package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

const PlatformBankCardTableName = "platform_bank_card"

const (
	//状态：1-启用， 2-禁用
	PlatformBankCardEnable  = 1
	PlatformBankCardDisable = 2
)

type PlatformBankCard struct {
	Id            int64  `gorm:"id"`             // id
	BankName      string `gorm:"bank_name"`      // 银行名称
	AccountName   string `gorm:"account_name"`   // 开户名
	CreateTime    int64  `gorm:"create_time"`    // 创建时间
	CardNumber    string `gorm:"card_number"`    // 银行卡号
	BranchName    string `gorm:"branch_name"`    // 支行名称
	Currency      string `gorm:"currency"`       // 币种
	MaxAmount     int64  `gorm:"max_amount"`     // 最大收款额度
	QrCode        string `gorm:"qr_code"`        // 收款二维码
	Remark        string `gorm:"remark"`         // 备注
	Status        int64  `gorm:"status"`         // 状态：1-启用， 2-禁用
	TodayReceived int64  `gorm:"today_received"` // 今日已收金额
	BankCode      string `gorm:"bank_code"`
}

func (t *PlatformBankCard) TableName() string {
	return PlatformBankCardTableName
}

func NewPlatformBankCardModel(db *gorm.DB) *PlatformBankCardModel {
	return &PlatformBankCardModel{db: db}
}

type PlatformBankCardModel struct {
	db *gorm.DB
}

// 插入一条记录
func (m *PlatformBankCardModel) Insert(data *PlatformBankCard) error {
	data.CreateTime = time.Now().Unix()
	result := m.db.Create(data)
	return result.Error
}

// 修改
func (m *PlatformBankCardModel) Update(id int64, data PlatformBankCard) error {
	setMap := map[string]interface{}{
		"branch_name": data.BranchName,
		"qr_code":     data.QrCode,
		"remark":      data.Remark,
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

	if data.Currency != "" {
		setMap["currency"] = data.Currency
	}

	if data.MaxAmount != 0 {
		setMap["max_amount"] = data.MaxAmount
	}

	result := m.db.Model(&PlatformBankCard{Id: id}).Updates(&setMap)
	return result.Error
}

// 删除
func (m *PlatformBankCardModel) Delete(id int64) error {
	result := m.db.Delete(&PlatformBankCard{Id: id})
	return result.Error
}

// 检查
func (m *PlatformBankCardModel) CheckById(id int64) (bool, error) {
	var cnt int64
	if err := m.db.Model(&PlatformBankCard{}).Where("id = ? ", id).Count(&cnt).Error; err != nil {
		return false, err
	}

	return cnt > 0, nil
}

// 检查银行，卡号，币种是否重复
func (m *PlatformBankCardModel) CheckBankCard(bankName, cardNumber, currency string) (bool, error) {
	var cnt int64
	result := m.db.Model(&PlatformBankCard{}).Where("bank_name = ? and card_number= ? and currency= ?", bankName, cardNumber, currency).Count(&cnt)
	if result.Error != nil {
		return false, result.Error
	}

	return cnt > 0, nil
}

// 根据id查询一条记录
func (m *PlatformBankCardModel) FindOneById(id int64) (*PlatformBankCard, error) {
	var o PlatformBankCard
	m.db.Model(&PlatformBankCard{}).Where("id = ? ", id).First(&o)
	return &o, nil
}

type FindPlatformBankCardList struct {
	Search   string
	Currency string
	Page     int64
	PageSize int64
	Status   int64
}

// 查询列表
func (m *PlatformBankCardModel) FindList(f FindPlatformBankCardList) ([]*PlatformBankCard, int64, error) {
	var (
		//card          = fmt.Sprintf("%s p", PlatformBankCardTableName)
		//rechargeOrder = fmt.Sprintf("left join %s m on m.platform_bank_card_id = p.id and m.order_status = ? ", MerchantRechargeOrderTableName)

		//selectField = "p.id, p.bank_name, p.account_name, p.create_time, p.card_number, " +
		//	"p.branch_name, p.currency, p.max_amount, p.remark, p.status "

		whereStr = " 1=1 "
		args     []interface{}
	)

	if f.Search != "" {
		whereStr += "and (account_name like ? or card_number like ? )"
		args = append(args, "%"+f.Search+"%", "%"+f.Search+"%")
	}
	if f.Currency != "" {
		whereStr += "and currency = ? "
		args = append(args, f.Currency)
	}

	if f.Status != 0 {
		whereStr += "and status = ? "
		args = append(args, f.Status)
	}

	if f.Page == 0 {
		f.Page = 1
	}
	if f.PageSize == 0 {
		f.PageSize = 10
	}

	var cnt int64
	if err := m.db.Model(&PlatformBankCard{}).Where(whereStr, args...).Count(&cnt).Error; err != nil {
		return nil, 0, err
	}

	//whereStr += "group by p.id order by p.create_time desc "
	whereStr += "order by create_time desc limit ? offset ? "
	args = append(args, f.PageSize, (f.Page-1)*f.PageSize)

	var resp []*PlatformBankCard
	if err := m.db.Model(&PlatformBankCard{}).Where(whereStr, args...).Scan(&resp).Error; err != nil {
		return nil, 0, err
	}

	return resp, cnt, nil
}

// 启用
func (m *PlatformBankCardModel) EnableCard(id int64) error {
	result := m.db.Model(&PlatformBankCard{Id: id}).Update("status", PlatformBankCardEnable)
	return result.Error
}

// 禁用
func (m *PlatformBankCardModel) DisableCard(id int64) error {
	result := m.db.Model(&PlatformBankCard{Id: id}).Update("status", PlatformBankCardDisable)
	return result.Error
}

// 累加银行卡今日已收款
func (m *PlatformBankCardModel) PlusTodayReceived(id, amount int64) error {
	sqlStr := fmt.Sprintf("UPDATE %s SET today_received=today_received+? WHERE id = ?", PlatformBankCardTableName)
	result := m.db.Exec(sqlStr, amount, id)
	return result.Error
}

// 今日已收金额清零
func (m *PlatformBankCardModel) TodayReceivedClear() error {
	sqlStr := fmt.Sprintf("UPDATE %s SET today_received=0 WHERE today_received > 0", PlatformBankCardTableName)
	result := m.db.Exec(sqlStr)
	return result.Error
}
