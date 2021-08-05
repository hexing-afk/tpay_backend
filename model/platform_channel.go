package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

const PlatformChannelTableName = "platform_channel"

const (
	// 状态
	PlatformChannelStatusEnable  = 1 // 启用
	PlatformChannelStatusDisable = 2 // 禁用

	// 通道类型
	PlatformChannelTypeCollection = 1 // 代收
	PlatformChannelTypeTransfer   = 2 // 代付
)

type PlatformChannel struct {
	Id          int64  `gorm:"id"`           // 下游通道id
	ChannelName string `gorm:"channel_name"` // 通道名称
	ChannelCode string `gorm:"channel_code"` // 通道代码
	CreateTime  int64  `gorm:"create_time"`  // 创建时间
	ChannelDesc string `gorm:"channel_desc"` // 通道描述
	Status      int64  `gorm:"status"`       // 状态: 1-启用， 2-禁用
	ChannelType int64  `gorm:"channel_type"` // 通道类型:1-代收, 2-代付
	UpdateTime  int64  `gorm:"update_time"`  // 更新时间
	AreaId      int64  `gorm:"area_id"`      // 地区id
	StartTime   int64  `gorm:"start_time"`   // 开始时间
	EndTime     int64  `gorm:"end_time"`     // 结束时间
	StartAmount int64  `gorm:"start_amount"` // 开始金额
	EndAmount   int64  `gorm:"end_amount"`   // 结束金额
}

func (m *PlatformChannel) TableName() string {
	return PlatformChannelTableName
}

func NewPlatformChannelModel(db *gorm.DB) *PlatformChannelModel {
	return &PlatformChannelModel{db: db}
}

type PlatformChannelModel struct {
	db *gorm.DB
}

// 插入一条记录
func (m *PlatformChannelModel) Insert(data *PlatformChannel) error {
	data.CreateTime = time.Now().Unix()
	return m.db.Create(data).Error
}

// 更新有值的字段
func (m *PlatformChannelModel) Update(id int64, data PlatformChannel) error {
	data.UpdateTime = time.Now().Unix()
	result := m.db.Model(&PlatformChannel{Id: id}).Updates(&data)
	return result.Error
}

// 删除一条记录
func (m *PlatformChannelModel) Delete(id int64) error {
	result := m.db.Delete(&PlatformChannel{Id: id})
	return result.Error
}

// 查询全部通道(id，名称, 类型)
func (m *PlatformChannelModel) FindManyNotIds(ids []int64, status, channelType, areaId int64) ([]*PlatformChannel, error) {
	var list []*PlatformChannel

	whereStr := " 1 = 1 "

	var whereParams []interface{}

	if status != 0 {
		whereStr += " and status = ? "
		whereParams = append(whereParams, status)
	}

	if channelType != 0 {
		whereStr += " and channel_type = ? "
		whereParams = append(whereParams, channelType)
	}
	if areaId != 0 {
		whereStr += " and area_id = ? "
		whereParams = append(whereParams, areaId)
	}

	result := m.db.Model(&PlatformChannel{}).Not(ids).
		Select("id, channel_name, channel_type ").
		Where(whereStr, whereParams...).
		Order("create_time desc ").Scan(&list)
	if result.Error != nil {
		return nil, result.Error
	}

	return list, nil
}

// 根据id查询
func (m *PlatformChannelModel) FindOneById(id int64) (*PlatformChannel, error) {
	var o = &PlatformChannel{}
	result := m.db.Model(&PlatformChannel{}).Where("id = ? ", id).First(o)
	return o, result.Error
}

type FindPlatformChannelList struct {
	Search   string
	Page     int64
	PageSize int64
}

type PlatformChannelList struct {
	PlatformChannel
	UpstreamChannelName string `gorm:"upstream_channel_name"` // 上游通道名称
	AreaName            string `gorm:"area_name"`             // 地区名称
}

// 查询列表
func (m *PlatformChannelModel) FindList(f FindPlatformChannelList) ([]*PlatformChannelList, int64, error) {
	var (
		whereStr = " 1=1 "
		args     []interface{}
	)

	if f.Search != "" {
		whereStr += "and pch.channel_name like ? "
		args = append(args, "%"+f.Search+"%")
	}

	if f.Page == 0 {
		f.Page = 1
	}
	if f.PageSize == 0 {
		f.PageSize = 10
	}

	platformChannel := fmt.Sprintf("%s pch ", PlatformChannelTableName)
	platformChannelUpstream := fmt.Sprintf("left join %s pcu on pcu.platform_channel_id = pch.id ", PlatformChannelUpstreamTableName)
	upstreamChannel := fmt.Sprintf("left join %s up on up.id = pcu.upstream_channel_id ", UpstreamChannelTableName)
	area := fmt.Sprintf("left join %s ar on ar.id = pch.area_id ", AreaTableName)

	selectField := "pch.id, pch.channel_name, pch.channel_code, pch.channel_desc, " +
		"pch.channel_type, pch.status, pch.create_time, pch.update_time, ar.area_name, pch.area_id, " +
		"GROUP_CONCAT(up.channel_name) as upstream_channel_name "

	whereStr += "GROUP BY pch.id "

	var total int64
	err := m.db.Table(platformChannel).
		Joins(platformChannelUpstream).
		Joins(upstreamChannel).
		Joins(area).
		Where(whereStr, args...).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	whereStr += "ORDER BY pch.create_time DESC "
	whereStr += "limit ? offset ? "
	args = append(args, f.PageSize, (f.Page-1)*f.PageSize)

	var list []*PlatformChannelList
	err = m.db.Table(platformChannel).
		Select(selectField).
		Joins(platformChannelUpstream).
		Joins(upstreamChannel).
		Joins(area).
		Where(whereStr, args...).
		Scan(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

// 启用通道
func (m *PlatformChannelModel) EnableChannel(id int64) error {
	result := m.db.Model(&PlatformChannel{Id: id}).Update("status", PlatformChannelStatusEnable)
	return result.Error
}

// 禁用通道
func (m *PlatformChannelModel) DisableChannel(id int64) error {
	result := m.db.Model(&PlatformChannel{Id: id}).Update("status", PlatformChannelStatusDisable)
	return result.Error
}

type PlatformChannelUpChannel struct {
	Id                 int64  `gorm:"id"`
	ChannelType        int64  `gorm:"channel_type"`
	UpstreamChannelIds string `gorm:"upstream_channel_ids"`
}

func (m *PlatformChannelModel) FindUpChannelById(id int64) (*PlatformChannelUpChannel, error) {
	sqlStr := "select pc.id, pc.channel_type, GROUP_CONCAT(pcu.upstream_channel_id) as upstream_channel_ids " +
		"from platform_channel pc " +
		"left join platform_channel_upstream pcu on pcu.platform_channel_id = pc.id " +
		"where pc.id = ? GROUP BY pc.id "

	var o = &PlatformChannelUpChannel{}
	result := m.db.Raw(sqlStr, id).Scan(o)

	if result.Error != nil {
		return nil, result.Error
	}

	return o, nil
}
