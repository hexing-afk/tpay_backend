package model

import "gorm.io/gorm"

const AreaTableName = "area"

type Area struct {
	Id         int64  `gorm:"id"`
	AreaName   string `gorm:"area_name"`
	CreateTime int64  `gorm:"create_time"`
}

func (t *Area) TableName() string {
	return AreaTableName
}

func NewAreaModel(db *gorm.DB) *AreaModel {
	return &AreaModel{db: db}
}

type AreaModel struct {
	db *gorm.DB
}

func (m *AreaModel) FindMany() ([]*Area, error) {
	var list []*Area
	err := m.db.Table(AreaTableName).Scan(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}

//查询地区是否存在
func (m *AreaModel) Check(id int64) (bool, error) {
	var cnt int64
	err := m.db.Table(AreaTableName).Where("id", id).Count(&cnt).Error
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}
