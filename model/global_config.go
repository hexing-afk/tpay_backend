package model

import (
	"strings"
	"time"

	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

const GlobalConfigTableName = "global_config"

const (
	// configKey 常量
	ConfigSiteConfig             = "site_config"               //网站配置
	ConfigImageBaseUrl           = "image_base_url"            // 图片域名地址
	ConfigPayapiHostAddr         = "payapi_host_addr"          // payapi站点域名地址
	ConfigPayTradeTypeSlice      = "pay_trade_type_slice"      // 代收交易方式
	ConfigTransferTradeTypeSlice = "transfer_trade_type_slice" // 代付交易类型

	// isChange 常量
	IsChangeTrue  = 1 // 是
	IsChangeFalse = 2 // 否

	// 所有后台登录谷歌验证码是否关闭(1关闭，0开启),此开关主要是方便开发环境
	ConfigTotpIsClose    = "totp_is_close"
	ConfigTotpIsCloseYes = "1" // 关闭

	BaseBatchTransferFileName = "base_batch_transfer_file_name" //批量付款文件名称

)

type GlobalConfig struct {
	ConfigKey   string `gorm:"config_key"`   // 全局配置key
	ConfigValue string `gorm:"config_value"` // 全局配置value
	Remark      string `gorm:"remark"`       // 备注
	CreateTime  int64  `gorm:"create_time"`
	IsChange    int64  `gorm:"is_change"` // 是否可改: 1-是, 2-否
}

func (t *GlobalConfig) TableName() string {
	return GlobalConfigTableName
}

func NewGlobalConfigModel(db *gorm.DB) *GlobalConfigModel {
	return &GlobalConfigModel{db: db}
}

type GlobalConfigModel struct {
	db *gorm.DB
}

func (m *GlobalConfigModel) Insert(data *GlobalConfig) error {
	data.CreateTime = time.Now().Unix()
	return m.db.Create(data).Error
}

func (m *GlobalConfigModel) InsertOrUpdate(data GlobalConfig) error {
	result := m.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "config_key"}},
		DoUpdates: clause.AssignmentColumns([]string{"config_value"}),
	}).Create(&data)

	return result.Error
}

func (m *GlobalConfigModel) FindValueByKey(key string) (string, error) {
	var value string
	err := m.db.Model(&GlobalConfig{}).Select("config_value").Where("config_key", key).Scan(&value).Error
	if err != nil {
		return "", err
	}
	return value, nil
}

// 是否开启google验证码
func (m *GlobalConfigModel) TotpIsClose() (bool, error) {
	val, err := m.FindValueByKey(ConfigTotpIsClose)
	if err != nil {
		return false, err
	}

	val = strings.TrimSpace(val)

	if val == ConfigTotpIsCloseYes { // 关闭
		return true, err
	}

	return false, nil
}

func (m *GlobalConfigModel) BaseBatchTransferFileName() (string, error) {
	val, err := m.FindValueByKey(BaseBatchTransferFileName)
	if err != nil {
		return "", err
	}

	val = strings.TrimSpace(val)

	return val, nil
}
