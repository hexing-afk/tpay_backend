package model

import (
	"gorm.io/gorm"
	"time"
)

const AdminWebLogTableName = "admin_web_log"

const (
	//日志类型
	LogTypePlatformConfig = 1 // 平台配置相关
	LogTypeMerchant       = 2 // 商户相关
)

type AdminWebLog struct {
	LogNo       string `gorm:"log_no"`      // 日志号
	AdminId     int64  `gorm:"admin_id"`    // 管理员id
	Description string `gorm:"description"` // 描述
	Type        int64  `gorm:"type"`        // 日志类型：1-平台配置相关， 2-商家相关，
	CreateTime  int64  `gorm:"create_time"` // 创建时间
}

func (t *AdminWebLog) TableName() string {
	return AdminWebLogTableName
}

type AdminWebLogModel struct {
	db *gorm.DB
}

func NewAdminWebLogModel(db *gorm.DB) *AdminWebLogModel {
	return &AdminWebLogModel{db: db}
}

// 插入一条记录
func (m *AdminWebLogModel) Insert(data *AdminWebLog) error {
	data.CreateTime = time.Now().Unix()
	result := m.db.Create(data)
	return result.Error
}
