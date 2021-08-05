package model

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"tpay_backend/utils"

	"gorm.io/gorm"
)

const (
	MerchantStatusEnable  = 1 // 启用
	MerchantStatusDisable = 2 // 禁用

	// 是否是平台商户
	IsBelongPlatformTrue  = 1 // 是
	IsBelongPlatformFalse = 2 // 否
)

const MerchantTableName = "merchant"

type Merchant struct {
	Id               int64    `gorm:"id"`
	MerchantNo       string   `gorm:"merchant_no"`   // 商户编号
	Username         string   `gorm:"username"`      // 用户名
	Password         string   `gorm:"password"`      // 密码
	CreateTime       int64    `gorm:"create_time"`   // 创建时间
	UpdateTime       int64    `gorm:"update_time"`   // 更新时间
	Phone            string   `gorm:"phone"`         // 手机号
	Email            string   `gorm:"email"`         // 邮箱
	Currency         string   `gorm:"currency"`      // 币种(来源于currency表)
	Status           int64    `gorm:"status"`        // 账号状态：1-启用, 2-禁用
	Md5Key           string   `gorm:"md5_key"`       // md5通信秘钥(32位)
	Balance          int64    `gorm:"balance"`       // 余额
	FrozenAmount     int64    `gorm:"frozen_amount"` // 冻结金额
	PayPassword      string   `gorm:"pay_password"`  // 支付密码
	IpWhiteList      string   `gorm:"ip_white_list"` // ip白名单(英文逗号分隔)
	IpWhiteListSlice []string `gorm:"-"`
	TotpSecret       string   `gorm:"totp_secret"` // TOTP认证秘钥
	AreaId           int64    `gorm:"area_id"`     // 地区id
}

func (t *Merchant) TableName() string {
	return MerchantTableName
}

func (t *Merchant) AfterFind(tx *gorm.DB) (err error) {
	t.IpWhiteListSlice = strings.Split(utils.TrimCommaStr(t.IpWhiteList), ",")
	return
}

func NewMerchantModel(db *gorm.DB) *MerchantModel {
	return &MerchantModel{db: db}
}

type MerchantModel struct {
	db *gorm.DB
}

// 插入一条记录
func (m *MerchantModel) Insert(data *Merchant) error {
	data.CreateTime = time.Now().Unix()
	result := m.db.Create(data)
	return result.Error
}

func (m *MerchantModel) Save(data *Merchant) error {
	result := m.db.Save(data)
	return result.Error
}

// 检查账号是否存在
func (m *MerchantModel) CheckByName(username string) (bool, error) {
	var cnt int64
	err := m.db.Model(&Merchant{}).Where("username=?", username).Count(&cnt).Error
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}

// 修改商户信息-包括0值和空值
func (m *MerchantModel) Update(id int64, data Merchant) error {
	setMap := map[string]interface{}{
		"update_time": time.Now().Unix(),
		"phone":       data.Phone,
		"email":       data.Email,
	}
	result := m.db.Model(&Merchant{Id: id}).Updates(setMap)
	return result.Error
}

func (m *MerchantModel) UpdateTotpSecret(id int64, data Merchant) error {
	setMap := map[string]interface{}{
		"update_time": time.Now().Unix(),
		"totp_secret": data.TotpSecret,
	}
	result := m.db.Model(&Merchant{Id: id}).Updates(setMap)
	return result.Error
}

func (m *MerchantModel) FindOneById(id int64) (*Merchant, error) {
	a := &Merchant{}
	result := m.db.Model(a).Where("id=?", id).First(a)
	return a, result.Error
}

func (m *MerchantModel) FindOneByMerchantNo(merchantNo string) (*Merchant, error) {
	a := &Merchant{}
	result := m.db.Model(a).Where("merchant_no=?", merchantNo).Limit(1).First(a)
	return a, result.Error
}

func (m *MerchantModel) FindOneByUsername(username string) (*Merchant, error) {
	a := &Merchant{}
	result := m.db.Model(a).Where("username=?", username).Limit(1).First(a)
	return a, result.Error
}

// 查询币种对应的平台商户
func (m *MerchantModel) FindPlatformMerchant(currency string, areaId int64) (*Merchant, error) {
	a := &Merchant{}
	result := m.db.Model(a).Where("is_belong_platform=? and currency=? and area_id=?", IsBelongPlatformTrue, currency, areaId).First(a)
	return a, result.Error
}

// 修改登录密码
func (m *MerchantModel) UpdatePassword(id int64, newPassword string) error {
	return m.db.Model(Merchant{Id: id}).Update("password", newPassword).Error
}

// 启用商户
func (m *MerchantModel) EnableMerchant(id int64) error {
	result := m.db.Model(&Merchant{Id: id}).Update("status", MerchantStatusEnable)
	return result.Error
}

// 禁用商户
func (m *MerchantModel) DisableMerchant(id int64) error {
	result := m.db.Model(&Merchant{Id: id}).Update("status", MerchantStatusDisable)
	return result.Error
}

type FindMerchantList struct {
	Username        string
	ContactDetails  string
	Currency        string
	StartCreateTime int64
	EndCreateTime   int64
	Page            int64
	PageSize        int64
}

type MerchantList struct {
	Merchant
	AreaName string `gorm:"area_name"` // 地区名称
}

// 查询商户列表
func (m *MerchantModel) FindList(f FindMerchantList) ([]*MerchantList, int64, error) {
	var (
		whereStr = " 1=1 "
		args     []interface{}
	)
	if f.Username != "" {
		whereStr += "and m.username like ? "
		args = append(args, "%"+f.Username+"%")
	}
	if f.Currency != "" {
		whereStr += "and m.currency = ? "
		args = append(args, f.Currency)
	}
	if f.ContactDetails != "" {
		whereStr += "and (m.phone like ? or m.email like ?) "
		args = append(args, "%"+f.ContactDetails+"%", "%"+f.ContactDetails+"%")
	}
	if f.StartCreateTime > 0 {
		whereStr += "and m.create_time >= ? "
		args = append(args, f.StartCreateTime)
	}
	if f.EndCreateTime > 0 {
		whereStr += "and m.create_time <= ? "
		args = append(args, f.EndCreateTime)
	}

	if f.Page == 0 {
		f.Page = 1
	}
	if f.PageSize == 0 {
		f.PageSize = 10
	}

	merchantTable := fmt.Sprintf("%s m", MerchantTableName)
	areaTable := fmt.Sprintf("left join %s ar on ar.id = m.area_id ", AreaTableName)
	selectStr := "m.*, ar.area_name "

	var cnt int64
	err := m.db.Table(merchantTable).Joins(areaTable).Where(whereStr, args...).Count(&cnt).Error
	if err != nil {
		return nil, 0, err
	}

	whereStr += "order by m.create_time desc "
	whereStr += "limit ? offset ? "
	args = append(args, f.PageSize, (f.Page-1)*f.PageSize)

	var resp []*MerchantList
	err = m.db.Table(merchantTable).Joins(areaTable).Select(selectStr).Where(whereStr, args...).Scan(&resp).Error
	if err != nil {
		return nil, 0, err
	}

	return resp, cnt, nil
}

type WalletLogExt struct {
	BusinessNo string
	Source     int64
	Remark     string
}

// 解结(加余额，减冻结金额)
func (m *MerchantModel) PlusBalanceUnfreezeTx(merchantId int64, amount int64, walletLogExt WalletLogExt) error {
	// 加余额
	if err := m.PlusBalance(merchantId, amount, walletLogExt); err != nil {
		return err
	}

	// 减冻结金额
	if err := m.MinusFrozenAmount(merchantId, amount, walletLogExt); err != nil {
		return err
	}

	return nil
}

// 冻冻(减余额，加冻结金额)
func (m *MerchantModel) MinusBalanceFreezeTx(merchantId int64, amount int64, walletLogExt WalletLogExt) error {
	// 减余额
	if err := m.MinusBalance(merchantId, amount, walletLogExt); err != nil {
		return err
	}

	// 加冻结金额
	if err := m.PlusFrozenAmount(merchantId, amount, walletLogExt); err != nil {
		return err
	}

	return nil
}

// 加余额
func (m *MerchantModel) PlusBalance(merchantId int64, amount int64, walletLogExt WalletLogExt) error {
	sqlStr := fmt.Sprintf("UPDATE %s SET balance=balance+? WHERE id=?", MerchantTableName)
	result := m.db.Exec(sqlStr, amount, merchantId)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("The number of affected rows is 0 ")
	}

	// 查询余额
	balance, err := m.QueryBalance(merchantId)
	if err != nil {
		return err
	}

	// 2.插入余额变动记录
	merchantWalletLog := new(MerchantWalletLog)
	merchantWalletLog.MerchantId = merchantId              // 商户id
	merchantWalletLog.OpType = OpTypeAddBalance            // 变动类型：1+，2-
	merchantWalletLog.ChangeAmount = amount                // 变动金额
	merchantWalletLog.AfterBalance = balance               // 变动后余额
	merchantWalletLog.BusinessNo = walletLogExt.BusinessNo // 业务单号
	merchantWalletLog.Source = walletLogExt.Source         // 变动来源：1-手动调账
	merchantWalletLog.Remark = walletLogExt.Remark         // 备注
	if err := NewMerchantWalletLogModel(m.db).Insert(merchantWalletLog); err != nil {
		return err
	}
	return nil
}

// 减余额
func (m *MerchantModel) MinusBalance(merchantId int64, amount int64, walletLogExt WalletLogExt) error {
	sqlStr := fmt.Sprintf("UPDATE %s SET balance=balance-? WHERE id=? AND balance >= ?", MerchantTableName)
	result := m.db.Exec(sqlStr, amount, merchantId, amount)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("The number of affected rows is 0 ")
	}

	// 查询余额
	balance, err := m.QueryBalance(merchantId)
	if err != nil {
		return err
	}

	// 3.插入余额变动记录
	merchantWalletLog := new(MerchantWalletLog)
	merchantWalletLog.MerchantId = merchantId              // 商户id
	merchantWalletLog.OpType = OpTypeMinusBalance          // 变动类型：1+，2-
	merchantWalletLog.ChangeAmount = amount                // 变动金额
	merchantWalletLog.AfterBalance = balance               // 变动后余额
	merchantWalletLog.BusinessNo = walletLogExt.BusinessNo // 业务单号
	merchantWalletLog.Source = walletLogExt.Source         // 变动来源：1-手动调账
	merchantWalletLog.Remark = walletLogExt.Remark         // 备注
	if err := NewMerchantWalletLogModel(m.db).Insert(merchantWalletLog); err != nil {
		return err
	}

	return nil
}

// 加冻结金额
func (m *MerchantModel) PlusFrozenAmount(merchantId int64, amount int64, walletLogExt WalletLogExt) error {
	sqlStr := fmt.Sprintf("UPDATE %s SET frozen_amount=frozen_amount+? WHERE id=?", MerchantTableName)
	result := m.db.Exec(sqlStr, amount, merchantId)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("The number of affected rows is 0 ")
	}

	// 查询冻结金额
	balance, err := m.QueryFrozenAmount(merchantId)
	if err != nil {
		return err
	}

	// 2.插入余额变动记录
	merchantWalletLog := new(MerchantWalletLog)
	merchantWalletLog.MerchantId = merchantId              // 商户id
	merchantWalletLog.OpType = OpTypeAddFrozen             // 变动类型：1+，2-, 3加冻结金额, 4减冻结金额
	merchantWalletLog.ChangeAmount = amount                // 变动金额
	merchantWalletLog.AfterBalance = balance               // 变动后余额
	merchantWalletLog.BusinessNo = walletLogExt.BusinessNo // 业务单号
	merchantWalletLog.Source = walletLogExt.Source         // 变动来源：1-手动调账
	merchantWalletLog.Remark = walletLogExt.Remark         // 备注
	if err := NewMerchantWalletLogModel(m.db).Insert(merchantWalletLog); err != nil {
		return err
	}
	return nil
}

// 减冻结金额
func (m *MerchantModel) MinusFrozenAmount(merchantId int64, amount int64, walletLogExt WalletLogExt) error {
	sqlStr := fmt.Sprintf("UPDATE %s SET frozen_amount=frozen_amount-? WHERE id=? and frozen_amount >= ?", MerchantTableName)
	result := m.db.Exec(sqlStr, amount, merchantId, amount)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("The number of affected rows is 0 ")
	}

	// 查询冻结金额
	balance, err := m.QueryFrozenAmount(merchantId)
	if err != nil {
		return err
	}

	// 3.插入余额变动记录
	merchantWalletLog := new(MerchantWalletLog)
	merchantWalletLog.MerchantId = merchantId              // 商户id
	merchantWalletLog.OpType = OpTypeMinusFrozen           // 变动类型：1+, 2-, 3加冻结金额, 4减冻结金额
	merchantWalletLog.ChangeAmount = amount                // 变动金额
	merchantWalletLog.AfterBalance = balance               // 变动后余额
	merchantWalletLog.BusinessNo = walletLogExt.BusinessNo // 业务单号
	merchantWalletLog.Source = walletLogExt.Source         // 变动来源：1-手动调账
	merchantWalletLog.Remark = walletLogExt.Remark         // 备注
	if err := NewMerchantWalletLogModel(m.db).Insert(merchantWalletLog); err != nil {
		return err
	}

	return nil
}

func (m *MerchantModel) QueryBalance(id int64) (int64, error) {
	sqlStr := fmt.Sprintf("SELECT balance FROM %s WHERE id=?", MerchantTableName)
	var balance int64
	result := m.db.Raw(sqlStr, id).Scan(&balance)
	return balance, result.Error
}

func (m *MerchantModel) QueryFrozenAmount(id int64) (int64, error) {
	sqlStr := fmt.Sprintf("SELECT frozen_amount FROM %s WHERE id=?", MerchantTableName)
	var balance int64
	result := m.db.Raw(sqlStr, id).Scan(&balance)
	return balance, result.Error
}

// 刷新md5密钥
func (m *MerchantModel) UpdateMd5Key(id int64, newMd5Key string) error {
	return m.db.Model(Merchant{Id: id}).Update("md5_key", newMd5Key).Error
}

// 修改支付密码
func (m *MerchantModel) UpdatePayPwd(id int64, newPayPwd string) error {
	return m.db.Model(Merchant{Id: id}).Update("pay_password", newPayPwd).Error
}

type MerchantChannelData struct {
	Merchant
	PlatformChannelIds string `gorm:"platform_channel_id"`
}

func (m *MerchantModel) FindMerchantChannel(merchantId int64) (*MerchantChannelData, error) {
	merchant := fmt.Sprintf("%s m", MerchantTableName)
	channel := fmt.Sprintf("left join %v c on c.merchant_id = m.id", MerchantChannelTableName)

	whereStr := "m.id = ? group by m.id"

	data := &MerchantChannelData{}
	result := m.db.Table(merchant).
		Select("m.id, m.status, m.currency, GROUP_CONCAT(c.platform_channel_id) as platform_channel_id ").
		Joins(channel).
		Where(whereStr, merchantId).Scan(data)
	return data, result.Error
}

// 修改白名单
func (m *MerchantModel) UpdateIpWhite(id int64, ipWhite string) error {
	return m.db.Model(Merchant{Id: id}).Update("ip_white_list", ipWhite).Error
}
