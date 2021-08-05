package model

import (
	"time"

	"gorm.io/gorm"
)

const AdminTableName = "admin"

const (
	// 管理员状态
	AdminEnableStatus  = 1 // 启用
	AdminDisableStatus = 2 // 禁用
)

type Admin struct {
	Id           int64  `gorm:"id"`
	Username     string `gorm:"username"`      // 用户名
	Password     string `gorm:"password"`      // 密码
	CreateTime   int64  `gorm:"create_time"`   // 创建时间
	UpdateTime   int64  `gorm:"update_time"`   // 更新时间
	EnableStatus int64  `gorm:"enable_status"` // 启用状态: 2-禁用, 1-启用
	Phone        string `gorm:"phone"`         // 手机号
	Email        string `gorm:"email"`         // 邮箱
	TotpSecret   string `gorm:"totp_secret"`   // TOTP认证秘钥
}

func (t *Admin) TableName() string {
	return AdminTableName
}

type AdminModel struct {
	db *gorm.DB
}

func NewAdminModel(db *gorm.DB) *AdminModel {
	return &AdminModel{db: db}
}

// 插入一条记录，返回自增id
func (m *AdminModel) Insert(data *Admin) (e error) {
	data.CreateTime = time.Now().Unix()
	result := m.db.Create(data)
	return result.Error
}

func (m *AdminModel) Save(data *Admin) error {
	result := m.db.Save(data)
	return result.Error
}

// 通过id字段查询1条记录
func (m *AdminModel) FindOneById(id int64) (*Admin, error) {
	a := &Admin{}
	result := m.db.Model(a).Where("id=?", id).First(a)
	return a, result.Error
}

// 通过username字段查询1条记录
func (m *AdminModel) FindOneByUsername(username string) (*Admin, error) {
	a := &Admin{}
	result := m.db.Model(a).Where("username=?", username).First(a)
	return a, result.Error
}

type FindAdminList struct {
	Page     int64
	PageSize int64
}

// 查询列表数据
func (m *AdminModel) FindList(f FindAdminList) ([]*Admin, int64, error) {
	if f.Page == 0 {
		f.Page = 1
	}
	if f.PageSize == 0 {
		f.PageSize = 10
	}

	var cnt int64
	if err := m.db.Table("admin").Count(&cnt).Error; err != nil {
		return nil, 0, err
	}

	whereStr := "1=1 order by create_time desc limit ? offset ? "
	args := []interface{}{f.PageSize, (f.Page - 1) * f.PageSize}

	var list []*Admin
	if err := m.db.Table("admin").Where(whereStr, args...).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, cnt, nil
}

// 通过id字段查询1条记录
func (m *AdminModel) CheckById(id int64) (bool, error) {
	var a int64
	result := m.db.Model(&Admin{}).Where("id = ?", id).Count(&a)
	return a > 0, result.Error
}

// 启用管理员
func (m *AdminModel) Enable(id int64) error {
	result := m.db.Model(&Admin{Id: id}).Update("enable_status", AdminEnableStatus)
	return result.Error
}

// 禁用管理员
func (m *AdminModel) Disable(id int64) error {
	result := m.db.Model(&Admin{Id: id}).Update("enable_status", AdminDisableStatus)
	return result.Error
}

// 修改登录密码
func (m *AdminModel) UpdatePassword(id int64, newPassword string) error {
	result := m.db.Model(&Admin{}).Where("id=?", id).Update("password", newPassword)
	return result.Error
}
