package model

import (
	"fmt"

	"gorm.io/gorm"
)

const MerchantChannelTableName = "merchant_channel"

const (
	MerchantChannelStatusEnable  = 1 // 启用
	MerchantChannelStatusDisable = 2 // 禁用
)

type MerchantChannel struct {
	Id                int64   `gorm:"id"`
	MerchantId        int64   `gorm:"merchant_id"`         // 商户Id
	PlatformChannelId int64   `gorm:"platform_channel_id"` // 平台通道id
	Rate              float64 `gorm:"rate"`                // 商户通道费率
	SingleFee         int64   `gorm:"single_fee"`          // 单笔手续费
	Status            int64   `gorm:"status"`              // 1-启用, 2-禁用
	CreateTime        int64   `gorm:"create_time"`
}

func (t *MerchantChannel) TableName() string {
	return MerchantChannelTableName
}

func NewMerchantChannelModel(db *gorm.DB) *MerchantChannelModel {
	return &MerchantChannelModel{db: db}
}

type MerchantChannelModel struct {
	db *gorm.DB
}

// 插入一条记录
func (m *MerchantChannelModel) Insert(data *MerchantChannel) error {
	result := m.db.Create(data)
	return result.Error
}

// 修改费率
func (m *MerchantChannelModel) UpdateRate(id int64, data MerchantChannel) error {
	setMap := map[string]interface{}{
		"single_fee": data.SingleFee,
		"rate":       data.Rate,
	}
	return m.db.Model(&MerchantChannel{Id: id}).Updates(&setMap).Error
}

func (m *MerchantChannelModel) Delete(id int64) error {
	return m.db.Delete(&MerchantChannel{Id: id}).Error
}

// 根据id查询记录
func (m *MerchantChannelModel) FindOneById(id int64) (*MerchantChannelList, error) {
	selectField := "mc.id, mc.merchant_id, mc.platform_channel_id, mc.rate, " +
		"mc.single_fee, mc.status, pc.channel_type AS platform_channel_type "

	sqlStr := fmt.Sprintf("SELECT %s FROM %s mc LEFT JOIN %s pc ON pc.id = mc.platform_channel_id ",
		selectField, MerchantChannelTableName, PlatformChannelTableName)
	sqlStr += "WHERE mc.id = ? "

	var o = &MerchantChannelList{}
	result := m.db.Raw(sqlStr, id).Scan(o)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrRecordNotFound
	}

	return o, nil
}

// 根据平台通道id查询记录
func (m *MerchantChannelModel) FindOneByPlatformId(platformId int64) ([]MerchantChannel, error) {
	var list []MerchantChannel
	m.db.Model(&MerchantChannel{}).Where("platform_channel_id=?", platformId).Scan(&list)
	return list, nil
}

type MerchantChannelList struct {
	MerchantChannel
	PlatformChannelName   string `gorm:"platform_channel_name"`   // 平台通道名称
	PlatformChannelType   int64  `gorm:"platform_channel_type"`   // 平台通道类型
	PlatformChannelCode   string `gorm:"platform_channel_code"`   // 平台通道代码
	UpstreamChannelWeight string `gorm:"upstream_channel_weight"` // 上游通道和权重
}

type FindMerchantChannelListReq struct {
	MerchantId            int64
	ChannelType           int64
	MerchantChannelStatus int64
}

// 查询列表
func (m *MerchantChannelModel) FindMerchantChannelList(req FindMerchantChannelListReq) ([]*MerchantChannelList, int64, error) {
	merCh := fmt.Sprintf("%s mc ", MerchantChannelTableName)
	merChUp := fmt.Sprintf("left join %s mcu on mcu.merchant_id = mc.merchant_id and mcu.merchant_channel_id = mc.id ", MerchantChannelUpstreamTableName)
	platCh := fmt.Sprintf("left join %s pc on pc.id = mc.platform_channel_id ", PlatformChannelTableName)
	upCh := fmt.Sprintf("left join upstream_channel up on up.id = mcu.upstream_channel_id ")

	selectField := "mc.id, mc.merchant_id, mc.platform_channel_id, mc.rate, mc.single_fee, mc.status, " +
		"pc.channel_name AS platform_channel_name, pc.channel_type AS platform_channel_type, pc.channel_code AS platform_channel_code, " +
		"GROUP_CONCAT(CONCAT_WS('-', up.channel_name, mcu.weight)) AS upstream_channel_weight "

	where := ` mc.merchant_id = ? `
	args := []interface{}{req.MerchantId}

	if req.ChannelType != 0 {
		where += ` AND pc.channel_type = ? `
		args = append(args, req.ChannelType)
	}
	if req.MerchantChannelStatus != 0 {
		where += ` AND mc.status = ? `
		args = append(args, req.MerchantChannelStatus)
	}

	where += ` GROUP BY mc.id `

	var total int64
	result := m.db.Table(merCh).Joins(merChUp).Joins(platCh).Joins(upCh).Where(where, args...).Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	if total == 0 {
		return nil, 0, nil
	}

	var list []*MerchantChannelList
	result = m.db.Table(merCh).Select(selectField).Joins(merChUp).Joins(platCh).Joins(upCh).
		Where(where, args...).Order("mc.create_time desc").Scan(&list)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return list, total, nil
}

// 查询商户已添加的平台通道id
func (m *MerchantChannelModel) FindMerchantPlatformChannelId(merchantId, channelType int64) ([]int64, error) {

	where := " mc.merchant_id = ? "
	args := []interface{}{merchantId}

	if channelType != 0 {
		where += "and plc.channel_type = ? "
		args = append(args, channelType)
	}

	var ids []int64
	result := m.db.Table(MerchantChannelTableName+" mc").
		Select("mc.platform_channel_id").
		Joins("left join "+PlatformChannelTableName+" plc on plc.id=mc.platform_channel_id").
		Where(where, args...).Scan(&ids)
	if result.Error != nil {
		return nil, result.Error
	}

	return ids, nil
}

// 启用商户通道
func (m *MerchantChannelModel) EnableMerchantChannel(id int64) error {
	result := m.db.Model(&MerchantChannel{Id: id}).Update("status", MerchantChannelStatusEnable)
	return result.Error
}

// 禁用商户通道
func (m *MerchantChannelModel) DisableMerchantChannel(id int64) error {
	result := m.db.Model(&MerchantChannel{Id: id}).Update("status", MerchantChannelStatusDisable)
	return result.Error
}

type MerchantUpstreamChannelData struct {
	MerchantId        int64   `gorm:"merchant_id"`
	Rate              float64 `gorm:"rate"`
	SingleFee         int64   `gorm:"single_fee"`
	Weight            int64   `gorm:"weight"`
	PlatformChannelId int64   `gorm:"platform_channel_id"`
	UpstreamChannelId int64   `gorm:"upstream_channel_id"`

	UpChannelDeductionMethod int64   `gorm:"up_channel_deduction_method"`
	UpChannelRate            float64 `gorm:"up_channel_rate"`
	UpChannelSingleFee       int64   `gorm:"up_channel_single_fee"`
	UpChannelCode            string  `gorm:"up_channel_code"`
}

// 查询商户代收通道
func (m *MerchantChannelModel) QueryPayUpstreamChannel(merchantId, merchantAreaId int64, platformChannelCode string, merchantCurrency string) ([]MerchantUpstreamChannelData, error) {
	var list []MerchantUpstreamChannelData

	selectStr := "mc.merchant_id, mc.platform_channel_id, mc.rate, mc.single_fee, mcu.weight, mcu.upstream_channel_id, " +
		" uc.channel_code AS up_channel_code, uc.deduction_method AS up_channel_deduction_method, uc.rate AS up_channel_rate, uc.single_fee AS up_channel_single_fee "

	whereStr := "mc.merchant_id=? AND mc.status=? AND mcu.weight > 0 " +
		"AND pc.channel_code=? AND pc.channel_type=? AND pc.status=? " +
		"AND uc.status=? AND uc.currency=? " +
		"AND pc.area_id=? "
	whereParams := []interface{}{merchantId, MerchantStatusEnable,
		platformChannelCode, PlatformChannelTypeCollection, PlatformChannelStatusEnable,
		UpstreamChannelStatusEnable, merchantCurrency,
		merchantAreaId,
	}

	result := m.db.Table(MerchantChannelTableName+" AS mc ").
		Select(selectStr).
		Joins("LEFT JOIN "+MerchantTableName+" AS m ON m.id = mc.merchant_id").
		Joins("LEFT JOIN "+MerchantChannelUpstreamTableName+" AS mcu ON mc.id=mcu.merchant_channel_id").
		Joins("LEFT JOIN "+PlatformChannelTableName+" AS pc ON mc.platform_channel_id=pc.id").
		Joins("LEFT JOIN "+UpstreamChannelTableName+" AS uc ON mcu.upstream_channel_id=uc.id").
		Where(whereStr, whereParams...).Scan(&list)

	if result.Error != nil {
		return nil, result.Error
	}

	return list, nil
}

// 查询商户代付通道
func (m *MerchantChannelModel) QueryTransferUpstreamChannel(merchantId, merchantAreaId int64, platformChannelCode string, merchantCurrency string) ([]MerchantUpstreamChannelData, error) {
	var list []MerchantUpstreamChannelData

	selectStr := "mc.merchant_id, mc.platform_channel_id, mc.rate, mc.single_fee, mcu.weight, mcu.upstream_channel_id, " +
		" uc.channel_code AS up_channel_code, uc.deduction_method AS up_channel_deduction_method, uc.rate AS up_channel_rate, uc.single_fee AS up_channel_single_fee "

	whereStr := "mc.merchant_id=? AND mc.status=? AND mcu.weight > 0 " +
		"AND pc.channel_code=? AND pc.channel_type=? AND pc.status=? " +
		"AND uc.status=? AND uc.currency=? " +
		"AND pc.area_id=? "
	whereParams := []interface{}{merchantId, MerchantStatusEnable,
		platformChannelCode, PlatformChannelTypeTransfer, PlatformChannelStatusEnable,
		UpstreamChannelStatusEnable, merchantCurrency,
		merchantAreaId,
	}

	result := m.db.Table(MerchantChannelTableName+" AS mc ").
		Select(selectStr).
		Joins("LEFT JOIN "+MerchantTableName+" AS m ON m.id = mc.merchant_id").
		Joins("LEFT JOIN "+MerchantChannelUpstreamTableName+" AS mcu ON mc.id=mcu.merchant_channel_id").
		Joins("LEFT JOIN "+PlatformChannelTableName+" AS pc ON mc.platform_channel_id=pc.id").
		Joins("LEFT JOIN "+UpstreamChannelTableName+" AS uc ON mcu.upstream_channel_id=uc.id").
		Where(whereStr, whereParams...).Scan(&list)

	if result.Error != nil {
		return nil, result.Error
	}

	return list, nil
}

// 查询商户代付通道
func (m *MerchantChannelModel) FindPlfMchTransferChannel(merchantId, merchantAreaId int64, merchantCurrency string) ([]MerchantChannelList, error) {
	var list []MerchantChannelList

	selectStr := "mc.id, mc.merchant_id, mc.rate, mc.single_fee, pc.channel_code AS platform_channel_code "

	whereStr := "mc.merchant_id=? AND mc.status=? AND mcu.weight > 0 " +
		"AND pc.channel_type=? AND pc.status=? " +
		"AND uc.status=? AND uc.currency=? " +
		"AND pc.area_id=? "
	whereParams := []interface{}{merchantId, MerchantStatusEnable,
		PlatformChannelTypeTransfer, PlatformChannelStatusEnable,
		UpstreamChannelStatusEnable, merchantCurrency,
		merchantAreaId,
	}

	result := m.db.Table(MerchantChannelTableName+" AS mc ").
		Select(selectStr).
		Joins("LEFT JOIN "+MerchantTableName+" AS m ON m.id = mc.merchant_id").
		Joins("LEFT JOIN "+MerchantChannelUpstreamTableName+" AS mcu ON mcu.merchant_id=mc.merchant_id").
		Joins("LEFT JOIN "+PlatformChannelTableName+" AS pc ON pc.id=mc.platform_channel_id").
		Joins("LEFT JOIN "+UpstreamChannelTableName+" AS uc ON uc.id=mcu.upstream_channel_id").
		Where(whereStr, whereParams...).Scan(&list)

	if result.Error != nil {
		return nil, result.Error
	}

	return list, nil
}

// 查询商户代付通道是否存在
func (m *MerchantChannelModel) FindPlfMchTransferChannelByPlfChannelId(merchantId, merchantAreaId int64, merchantCurrency string, channelId int64) (MerchantChannelList, error) {
	var data MerchantChannelList

	selectStr := "mc.id, mc.merchant_id, mc.rate, mc.single_fee, pc.channel_code AS platform_channel_code "

	whereStr := "mc.merchant_id=? AND mc.status=? AND mcu.weight > 0 " +
		"AND pc.channel_type=? AND pc.status=? " +
		"AND uc.status=? AND uc.currency=? " +
		"AND pc.area_id=? AND mc.platform_channel_id=? "
	whereParams := []interface{}{merchantId, MerchantStatusEnable,
		PlatformChannelTypeTransfer, PlatformChannelStatusEnable,
		UpstreamChannelStatusEnable, merchantCurrency,
		merchantAreaId, channelId,
	}

	result := m.db.Table(MerchantChannelTableName+" AS mc ").
		Select(selectStr).
		Joins("LEFT JOIN "+MerchantTableName+" AS m ON m.id = mc.merchant_id").
		Joins("LEFT JOIN "+MerchantChannelUpstreamTableName+" AS mcu ON mcu.merchant_id=mc.merchant_id").
		Joins("LEFT JOIN "+PlatformChannelTableName+" AS pc ON pc.id=mc.platform_channel_id").
		Joins("LEFT JOIN "+UpstreamChannelTableName+" AS uc ON uc.id=mcu.upstream_channel_id").
		Where(whereStr, whereParams...).Scan(&data)

	if result.Error != nil {
		return data, result.Error
	}

	return data, nil
}
