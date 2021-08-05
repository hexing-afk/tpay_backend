package model

import (
	"fmt"
	"gorm.io/gorm"
)

const MerchantChannelUpstreamTableName = "merchant_channel_upstream"

type MerchantChannelUpstream struct {
	Id                int64 `gorm:"id"`
	MerchantId        int64 `gorm:"merchant_id"`         // 商户id
	MerchantChannelId int64 `gorm:"merchant_channel_id"` // 商户通道id
	UpstreamChannelId int64 `gorm:"upstream_channel_id"` // 上游通道id
	Weight            int64 `gorm:"weight"`              // 权重
}

func (t *MerchantChannelUpstream) TableName() string {
	return MerchantChannelUpstreamTableName
}

func NewMerchantChannelUpstreamModel(db *gorm.DB) *MerchantChannelUpstreamModel {
	return &MerchantChannelUpstreamModel{db: db}
}

type MerchantChannelUpstreamModel struct {
	db *gorm.DB
}

func (m *MerchantChannelUpstreamModel) Inserts(data []*MerchantChannelUpstream) error {
	result := m.db.Create(&data)
	return result.Error
}

func (m *MerchantChannelUpstreamModel) UpdateWeight(id, weight int64) error {
	result := m.db.Model(&MerchantChannelUpstream{Id: id}).Update("weight", weight)
	return result.Error
}

func (m *MerchantChannelUpstreamModel) Delete(merchantChannelId int64) error {
	return m.db.Where("merchant_channel_id = ? ", merchantChannelId).Delete(&MerchantChannelUpstream{}).Error
}

func (m *MerchantChannelUpstreamModel) DeleteByUpChannel(merChannelId, upChannelId int64) error {
	result := m.db.Where("merchant_channel_id = ? and upstream_channel_id=?", merChannelId, upChannelId).Delete(&MerchantChannelUpstream{})
	return result.Error
}

type FindMerchantChannelUpstreamList struct {
	MerchantId        int64
	MerchantChannelId int64
	Page              int64
	PageSize          int64
}

type MerchantChannelUpstreamList struct {
	MerchantChannelUpstream
	UpstreamChannelName string `gorm:"upstream_channel_name"` // 上游通道名称
}

func (m *MerchantChannelUpstreamModel) FindList(f FindMerchantChannelUpstreamList) ([]*MerchantChannelUpstreamList, int64, error) {
	var (
		whereStr = " 1=1 "
		args     []interface{}
	)

	if f.MerchantId != 0 {
		whereStr += "and mcu.merchant_id = ? "
		args = append(args, f.MerchantId)
	}

	if f.MerchantChannelId != 0 {
		whereStr += "and mcu.merchant_channel_id = ? "
		args = append(args, f.MerchantChannelId)
	}

	if f.Page == 0 {
		f.Page = 1
	}
	if f.PageSize == 0 {
		f.PageSize = 10
	}

	merchantChannelUpstream := fmt.Sprintf("%s mcu ", MerchantChannelUpstreamTableName)
	upstreamChannel := fmt.Sprintf("left join %s uc on uc.id = mcu.upstream_channel_id ", UpstreamChannelTableName)

	selectFiled := "mcu.id, mcu.weight, uc.channel_name as upstream_channel_name"

	var total int64
	result := m.db.Table(merchantChannelUpstream).Joins(upstreamChannel).Where(whereStr, args...).Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	whereStr += "order by mcu.id desc "
	whereStr += "limit ? offset ? "
	args = append(args, f.PageSize, (f.Page-1)*f.PageSize)

	var list []*MerchantChannelUpstreamList
	result = m.db.Table(merchantChannelUpstream).Select(selectFiled).Joins(upstreamChannel).Where(whereStr, args...).Scan(&list)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return list, total, nil
}
