package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

const UpstreamChannelTableName = "upstream_channel"

const (
	// 状态
	UpstreamChannelStatusEnable  = 1 // 启用
	UpstreamChannelStatusDisable = 2 // 禁用

	// 通道类型
	UpstreamChannelTypeCollection = 1 // 代收
	UpstreamChannelTypeTransfer   = 2 // 代付

	// 扣费方式
	UpstreamChannelDeductionInner = 1 // 内扣
	UpstreamChannelDeductionOut   = 2 // 外扣
)

type UpstreamChannel struct {
	Id              int64   `gorm:"id"`
	ChannelName     string  `gorm:"channel_name"`      // 通道名称
	CreateTime      int64   `gorm:"create_time"`       // 创建时间
	ChannelCode     string  `gorm:"channel_code"`      // 通道代码
	ChannelDesc     string  `gorm:"channel_desc"`      // 通道描述
	Status          int64   `gorm:"status"`            // 状态: 1-启用， 2-禁用
	Currency        string  `gorm:"currency"`          // 通道币种
	UpstreamId      int64   `gorm:"upstream_id"`       // 上游id
	ChannelType     int64   `gorm:"channel_type"`      // 通道类型:1-代收, 2-代付
	UpdateTime      int64   `gorm:"update_time"`       // 更新时间
	DeductionMethod int64   `gorm:"deduction_method"`  // 扣费方式: 1-内扣， 2-外扣
	Rate            float64 `gorm:"rate"`              // 通道费率
	SingleFee       int64   `gorm:"single_fee"`        // 单笔手续费
	SingleMaxAmount int64   `gorm:"single_max_amount"` // 单笔最大金额
}

func (m *UpstreamChannel) TableName() string {
	return UpstreamChannelTableName
}

type UpstreamChannelModel struct {
	db *gorm.DB
}

func NewUpstreamChannelModel(db *gorm.DB) *UpstreamChannelModel {
	return &UpstreamChannelModel{db: db}
}

// 插入一条记录
func (m *UpstreamChannelModel) Insert(data *UpstreamChannel) error {
	data.CreateTime = time.Now().Unix()
	return m.db.Create(data).Error
}

// 更新有值的字段
func (m *UpstreamChannelModel) Update(id int64, data UpstreamChannel) error {
	data.UpdateTime = time.Now().Unix()
	result := m.db.Model(&UpstreamChannel{Id: id}).Updates(&data)
	return result.Error
}

// 删除一条记录
func (m *UpstreamChannelModel) Delete(id int64) error {
	result := m.db.Delete(&UpstreamChannel{Id: id})
	return result.Error
}

// 检查
func (m *UpstreamChannelModel) CheckById(id int64) (bool, error) {
	var cnt int64
	if err := m.db.Model(&UpstreamChannel{}).Where("id = ? ", id).Count(&cnt).Error; err != nil {
		return false, err
	}

	return cnt > 0, nil
}

// 根据Id关联查一个
func (m *UpstreamChannelModel) FindOneById(id int64) (*UpstreamChannelList, error) {
	upstreamChannel := fmt.Sprintf("%s upc ", UpstreamChannelTableName)
	upstream := fmt.Sprintf("left join %s up on up.id = upc.upstream_id ", UpstreamTableName)

	selectField := "upc.id, upc.channel_name, upc.create_time, upc.channel_code, upc.channel_desc, " +
		"upc.status, upc.currency, upc.upstream_id, upc.channel_type, upc.update_time, " +
		"upc.deduction_method, upc.rate, upc.single_fee, upc.single_max_amount, up.upstream_name "

	var channel = &UpstreamChannelList{}
	result := m.db.Table(upstreamChannel).Select(selectField).Joins(upstream).Where("upc.id = ? ", id).Scan(channel)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrRecordNotFound
	}

	return channel, nil
}

// 查询ids外的数据
func (m *UpstreamChannelModel) FindManyNotIds(ids []int64, channelType, areaId int64) ([]*UpstreamChannel, error) {
	var channel []*UpstreamChannel
	selectStr := "uc.*"
	upstreamJoins := fmt.Sprintf("left join %s u on u.id = uc.upstream_id ", UpstreamTableName)

	where := "uc.channel_type = ? and uc.status = ? "
	args := []interface{}{channelType, UpstreamChannelStatusEnable}

	if ids != nil {
		where += " and uc.id NOT IN ? "
		args = append(args, ids)
	}

	if areaId != 0 {
		where += " and u.area_id = ? "
		args = append(args, areaId)
	}
	upstreamChannel := fmt.Sprintf("%s uc", UpstreamChannelTableName)

	result := m.db.Table(upstreamChannel).Select(selectStr).Joins(upstreamJoins).
		Where(where, args...).
		Find(&channel)
	if result.Error != nil {
		return nil, result.Error
	}
	return channel, nil
}

// 查询ids的数据
func (m *UpstreamChannelModel) FindManyByIds(ids []int64, channelType int64) ([]*UpstreamChannel, error) {
	var channel []*UpstreamChannel
	result := m.db.Model(&UpstreamChannel{}).
		Where("id in ? and channel_type = ? and status = ? ", ids, channelType, UpstreamChannelStatusEnable).
		Find(&channel)
	if result.Error != nil {
		return nil, result.Error
	}
	return channel, nil
}

type FindUpstreamChannelList struct {
	Search          string
	Currency        string
	ChannelType     int64
	StartCreateTime int64
	EndCreateTime   int64
	DeductionMethod int64
	Status          int64
	Page            int64
	PageSize        int64
}

type UpstreamChannelList struct {
	UpstreamChannel
	UpstreamName string `gorm:"upstream_name"` // 上游通道名
}

// 查列表
func (m *UpstreamChannelModel) FindList(f FindUpstreamChannelList) ([]*UpstreamChannelList, int64, error) {
	var (
		upstreamChannel = fmt.Sprintf("%s upc", UpstreamChannelTableName)
		upstream        = fmt.Sprintf("left join %s up on up.id = upc.upstream_id ", UpstreamTableName)

		selectField = "upc.id, upc.channel_name, upc.channel_code, upc.channel_desc, upc.currency, " +
			"upc.channel_type, upc.rate, upc.deduction_method, upc.update_time, upc.status, up.upstream_name," +
			"upc.single_fee "

		whereStr = " 1=1 "
		args     []interface{}
	)

	if f.Search != "" {
		whereStr += "and (upc.channel_name like ? or up.upstream_name like ? )"
		args = append(args, "%"+f.Search+"%", "%"+f.Search+"%")
	}

	if f.ChannelType != 0 {
		whereStr += "and upc.channel_type = ? "
		args = append(args, f.ChannelType)
	}

	if f.Currency != "" {
		whereStr += "and upc.currency = ? "
		args = append(args, f.Currency)
	}

	if f.StartCreateTime > 0 {
		whereStr += "and upc.create_time >= ? "
		args = append(args, f.StartCreateTime)
	}

	if f.EndCreateTime > 0 {
		whereStr += "and upc.create_time <= ? "
		args = append(args, f.EndCreateTime)
	}

	if f.DeductionMethod != 0 {
		whereStr += "and upc.deduction_method = ? "
		args = append(args, f.DeductionMethod)
	}

	if f.Status != 0 {
		whereStr += "and upc.status = ? "
		args = append(args, f.Status)
	}

	if f.Page == 0 {
		f.Page = 1
	}
	if f.PageSize == 0 {
		f.PageSize = 10
	}

	var cnt int64
	err := m.db.Table(upstreamChannel).Joins(upstream).Where(whereStr, args...).Count(&cnt).Error
	if err != nil {
		return nil, 0, err
	}

	if cnt == 0 {
		return nil, 0, nil
	}

	whereStr += "order by upc.create_time desc "
	whereStr += "limit ? offset ? "
	args = append(args, f.PageSize, (f.Page-1)*f.PageSize)

	var list []*UpstreamChannelList
	err = m.db.Table(upstreamChannel).Select(selectField).Joins(upstream).Where(whereStr, args...).Scan(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, cnt, nil
}

// 启用通道
func (m *UpstreamChannelModel) EnableChannel(id int64) error {
	result := m.db.Model(&UpstreamChannel{Id: id}).Update("status", UpstreamChannelStatusEnable)
	return result.Error
}

// 禁用通道
func (m *UpstreamChannelModel) DisableChannel(id int64) error {
	result := m.db.Model(&UpstreamChannel{Id: id}).Update("status", UpstreamChannelStatusDisable)
	return result.Error
}

// 通过上游通道id找到上游
func (m *UpstreamChannelModel) FindUpstreamByChannelId(upstreamChannelId int64) (*Upstream, error) {
	var up *Upstream

	result := m.db.Table(UpstreamChannelTableName+" AS uc ").
		Select("u.*").
		Joins("LEFT JOIN "+UpstreamTableName+" AS u ON uc.upstream_id=u.id").
		Where("uc.id=?", upstreamChannelId).First(&up)

	if result.Error != nil {
		return nil, result.Error
	}

	return up, nil
}
