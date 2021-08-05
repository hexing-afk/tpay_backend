package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

const UpstreamTableName = "upstream"

type Upstream struct {
	Id                 int64  `gorm:"id"`
	UpstreamName       string `gorm:"upstream_name"`        // 上游名称
	UpstreamMerchantNo string `gorm:"upstream_merchant_no"` // 上游账号id
	UpstreamCode       string `gorm:"upstream_code"`        // 上游代码
	CreateTime         int64  `gorm:"create_time"`          // 创建时间
	CallConfig         string `gorm:"call_config"`          // 通信配置
	AreaId             int64  `gorm:"area_id"`              // 地区id
}

func (t *Upstream) TableName() string {
	return UpstreamTableName
}

type FindUpstreamList struct {
	UpstreamName string
	Page         int64
	PageSize     int64
}

type UpstreamModel struct {
	db *gorm.DB
}

func NewUpstreamModel(db *gorm.DB) *UpstreamModel {
	return &UpstreamModel{db: db}
}

func (m *UpstreamModel) Insert(data *Upstream) error {
	data.CreateTime = time.Now().Unix()
	return m.db.Create(data).Error
}

func (m *UpstreamModel) Update(id int64, data Upstream) error {
	result := m.db.Model(&Upstream{Id: id}).Updates(&data)
	return result.Error
}

func (m *UpstreamModel) CheckById(id int64) (bool, error) {
	var cnt int64
	if err := m.db.Model(&Upstream{}).Where("id = ? ", id).Count(&cnt).Error; err != nil {
		return false, err
	}

	return cnt > 0, nil
}

func (m *UpstreamModel) FindOneById(id int64) (*Upstream, error) {
	var o = &Upstream{}
	if err := m.db.Model(o).Where("id=?", id).First(o).Error; err != nil {
		return nil, err
	}
	return o, nil
}

type UpstreamList struct {
	Upstream
	AreaName string `gorm:"area_name"` // 地区名称
}

func (m *UpstreamModel) FindList(f FindUpstreamList) ([]*UpstreamList, int64, error) {
	var (
		whereStr = " 1=1 "
		args     []interface{}
	)

	if f.UpstreamName != "" {
		whereStr += "and up.upstream_name like ? "
		args = append(args, "%"+f.UpstreamName+"%")
	}

	if f.Page == 0 {
		f.Page = 1
	}
	if f.PageSize == 0 {
		f.PageSize = 10
	}

	var cnt int64
	err := m.db.Model(&Upstream{}).Where(whereStr, args...).Count(&cnt).Error
	if err != nil {
		return nil, 0, err
	}

	whereStr += "order by up.create_time desc limit ? offset ? "
	args = append(args, f.PageSize, (f.Page-1)*f.PageSize)

	upstreamTableName := fmt.Sprintf("%s up", UpstreamTableName)
	areaTableName := fmt.Sprintf("left join %s ar on up.area_id = ar.id ", AreaTableName)

	selectStr := "up.*, ar.area_name "
	var resp []*UpstreamList
	err = m.db.Table(upstreamTableName).Select(selectStr).Joins(areaTableName).Where(whereStr, args...).Scan(&resp).Error
	if err != nil {
		return nil, cnt, err
	}
	return resp, cnt, nil
}

func (m *UpstreamModel) FindOneByUpstreamMerchantNo(upMerchantNo string) (*Upstream, error) {
	var o = &Upstream{}

	result := m.db.Model(o).Where("upstream_merchant_no=?", upMerchantNo).First(o)

	if result.Error != nil {
		return nil, result.Error
	}

	return o, nil
}
