package model

import (
	"fmt"
	"gorm.io/gorm"
)

const PlatformChannelUpstreamTableName = "platform_channel_upstream"

type PlatformChannelUpstream struct {
	Id                int64 `gorm:"id"`
	PlatformChannelId int64 `gorm:"platform_channel_id"` // 平台通道id
	UpstreamChannelId int64 `gorm:"upstream_channel_id"` // 上游通道id
}

type PlatformChannelUpstreamChannel struct {
	PlatformChannelUpstream
	PlatformChannelName string `gorm:"platform_channel_name"` // 平台通道名称
	UpstreamChannelName string `gorm:"upstream_channel_name"` // 上游通道名称
}

func (m *PlatformChannelUpstream) TableName() string {
	return PlatformChannelUpstreamTableName
}

func NewPlatformChannelUpstreamModel(db *gorm.DB) *PlatformChannelUpstreamModel {
	return &PlatformChannelUpstreamModel{db: db}
}

type PlatformChannelUpstreamModel struct {
	db *gorm.DB
}

// 插入多条记录
func (m *PlatformChannelUpstreamModel) Inserts(data []*PlatformChannelUpstream) error {
	result := m.db.Create(&data)
	return result.Error
}

// 根据下游通道id删除记录
func (m *PlatformChannelUpstreamModel) DeleteByPlatform(channelId int64) error {
	result := m.db.Where("platform_channel_id = ?", channelId).Delete(&PlatformChannelUpstream{})
	return result.Error
}

func (m *PlatformChannelUpstreamModel) DeleteByWhere(platformId, upstreamId int64) error {
	result := m.db.Where("platform_channel_id = ? and upstream_channel_id=?", platformId, upstreamId).Delete(&PlatformChannelUpstream{})
	return result.Error
}

// 根据id检查下游通道是否有关联上游通道
func (m *PlatformChannelUpstreamModel) CheckPlatformById(channelId int64) (bool, error) {
	platformChannelUpstream := fmt.Sprintf("%s pcu", PlatformChannelUpstreamTableName)
	upstreamChannel := fmt.Sprintf("left join %s upc on upc.id = pcu.upstream_channel_id ", UpstreamChannelTableName)

	var cnt int64
	result := m.db.Table(platformChannelUpstream).Joins(upstreamChannel).Where("pcu.platform_channel_id = ? ", channelId).Count(&cnt)
	if result.Error != nil {
		return false, result.Error
	}

	return cnt > 0, nil
}

// 根据id检查上游通道是否有下游通道关联
func (m *PlatformChannelUpstreamModel) CheckUpstreamById(channelId int64) (bool, error) {
	platformChannelUpstream := fmt.Sprintf("%s pcu", PlatformChannelUpstreamTableName)
	platformChannel := fmt.Sprintf("left join %s pc on pc.id = pcu.platform_channel_id ", PlatformChannelTableName)

	var cnt int64
	result := m.db.Table(platformChannelUpstream).Joins(platformChannel).Where("pcu.upstream_channel_id = ? ", channelId).Count(&cnt)

	if result.Error != nil {
		return false, result.Error
	}

	return cnt > 0, nil
}

// 根据id和type检查下游通道是否有关联上游通道
func (m *PlatformChannelUpstreamModel) CheckPlatformByIdAndType(channelId, channelType int64) (bool, error) {
	platformChannelUpstream := fmt.Sprintf("%s pcu", PlatformChannelUpstreamTableName)
	upstreamChannel := fmt.Sprintf("left join %s upc on upc.id = pcu.upstream_channel_id and upc.channel_type = ? ", UpstreamChannelTableName)

	var cnt int64
	result := m.db.Table(platformChannelUpstream).Joins(upstreamChannel, channelType).Where("pcu.platform_channel_id = ? ", channelId).Count(&cnt)
	if result.Error != nil {
		return false, result.Error
	}

	return cnt > 0, nil
}

// 根据id和type检查上游通道是否有下游通道关联
func (m *PlatformChannelUpstreamModel) CheckUpstreamByIdAndType(channelId, channelType int64) (bool, error) {
	platformChannelUpstream := fmt.Sprintf("%s pcu", PlatformChannelUpstreamTableName)
	platformChannel := fmt.Sprintf("left join %s pch on pch.id = pcu.platform_channel_id and pch.channel_type = ? ", PlatformChannelTableName)

	var cnt int64
	result := m.db.Table(platformChannelUpstream).Joins(platformChannel, channelType).Where("pcu.upstream_channel_id = ? ", channelId).Count(&cnt)
	if result.Error != nil {
		return false, result.Error
	}

	return cnt > 0, nil
}

// 根据下游通道id查询上游通道
func (m *PlatformChannelUpstreamModel) FindUpstreamByPlatform(channelId int64) ([]*PlatformChannelUpstreamChannel, error) {
	platformChannelUpstream := fmt.Sprintf("%s pcu", PlatformChannelUpstreamTableName)
	upstreamChannel := fmt.Sprintf("left join %s upc on upc.id = pcu.upstream_channel_id ", UpstreamChannelTableName)

	selectField := "pcu.upstream_channel_id, upc.channel_name as upstream_channel_name "

	var o []*PlatformChannelUpstreamChannel
	result := m.db.Table(platformChannelUpstream).Select(selectField).Joins(upstreamChannel).Where("pcu.platform_channel_id = ? ", channelId).Scan(&o)
	if result.Error != nil {
		return nil, result.Error
	}

	return o, nil
}

// 根据下游通道id和商户币种查询商户可用的上游通道
func (m *PlatformChannelUpstreamModel) FindMerchantUpstreamChannel(channelId int64, currency string) ([]*PlatformChannelUpstreamChannel, error) {
	platformChannelUpstream := fmt.Sprintf("%s pcu", PlatformChannelUpstreamTableName)
	upstreamChannel := fmt.Sprintf("left join %s upc on upc.id = pcu.upstream_channel_id", UpstreamChannelTableName)

	selectField := "pcu.upstream_channel_id, upc.channel_name as upstream_channel_name "

	var o []*PlatformChannelUpstreamChannel
	result := m.db.Table(platformChannelUpstream).
		Select(selectField).
		Joins(upstreamChannel).
		Where("pcu.platform_channel_id = ?  and upc.currency = ? ", channelId, currency).Scan(&o)
	if result.Error != nil {
		return nil, result.Error
	}

	return o, nil
}
