package model

import (
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
)

const UserCardTableName = "usercard"

type UserCard struct {
	Id         int64          `db:"id"`
	CardNo     string         `db:"card_no"`       // 卡号
	UpNo       sql.NullString `db:"up_no"`         // 签约号
	UpChNo     sql.NullString `db:"up_ch_no"`      // 上游通道号
	CreateTime int64          `gorm:"create_time"` // 创建时间
}

func (t *UserCard) TableName() string {
	return UserCardTableName
}

type FindUserCardList struct {
	UserCardName string
	Page         int64
	PageSize     int64
}

type UserCardModel struct {
	db *gorm.DB
}

func NewUserCardModel(db *gorm.DB) *UserCardModel {
	return &UserCardModel{db: db}
}

func (m *UserCardModel) Insert(data *UserCard) error {
	data.CreateTime = time.Now().Unix()
	return m.db.Create(data).Error
}

func (m *UserCardModel) Update(id int64, data UserCard) error {
	result := m.db.Model(&UserCard{Id: id}).Updates(&data)
	return result.Error
}

func (m *UserCardModel) CheckById(id int64) (bool, error) {
	var cnt int64
	if err := m.db.Model(&UserCard{}).Where("id = ? ", id).Count(&cnt).Error; err != nil {
		return false, err
	}

	return cnt > 0, nil
}

func (m *UserCardModel) FindOneById(id int64) (*UserCard, error) {
	var o = &UserCard{}
	if err := m.db.Model(o).Where("id=?", id).First(o).Error; err != nil {
		return nil, err
	}
	return o, nil
}

type UserCardList struct {
	UserCard
}

func (m *UserCardModel) FindList(f FindUserCardList) ([]*UserCardList, int64, error) {
	var (
		whereStr = " 1=1 "
		args     []interface{}
	)

	if f.UserCardName != "" {
		whereStr += "and up.usercard_name like ? "
		args = append(args, "%"+f.UserCardName+"%")
	}

	if f.Page == 0 {
		f.Page = 1
	}
	if f.PageSize == 0 {
		f.PageSize = 10
	}

	var cnt int64
	err := m.db.Model(&UserCard{}).Where(whereStr, args...).Count(&cnt).Error
	if err != nil {
		return nil, 0, err
	}

	whereStr += "order by up.create_time desc limit ? offset ? "
	args = append(args, f.PageSize, (f.Page-1)*f.PageSize)

	usercardTableName := fmt.Sprintf("%s up", UserCardTableName)
	areaTableName := fmt.Sprintf("left join %s ar on up.area_id = ar.id ", AreaTableName)

	selectStr := "up.*, ar.area_name "
	var resp []*UserCardList
	err = m.db.Table(usercardTableName).Select(selectStr).Joins(areaTableName).Where(whereStr, args...).Scan(&resp).Error
	if err != nil {
		return nil, cnt, err
	}
	return resp, cnt, nil
}

func (m *UserCardModel) FindOneByUserCardNo(cardNo string) (*UserCard, error) {
	var o = &UserCard{}

	result := m.db.Model(o).Where("card_no=?", cardNo).First(o)

	if result.Error != nil {
		return nil, result.Error
	}

	return o, nil
}

func (m *UserCardModel) FindOneByUserCardNoChannel(cardNo string, ch int64) (*UserCard, error) {
	var o = &UserCard{}

	result := m.db.Model(o).Where("card_no=? and up_ch_no=?", cardNo, ch).First(o)

	if result.Error != nil {
		return nil, result.Error
	}

	return o, nil
}
