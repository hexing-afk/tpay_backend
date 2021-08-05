package model

import (
	"gorm.io/gorm"
	"time"
)

const UploadFileLogTableName = "upload_file_log"

const (
	// 账号类型
	UploadFileLogAccountTypeMerchant int64 = 1 // 1商户
)

type UploadFileLog struct {
	Id          int64  `gorm:"id"`           //
	FileName    string `gorm:"file_name"`    //文件名称
	AccountId   int64  `gorm:"account_id"`   //账号id
	AccountType int64  `gorm:"account_type"` //账号类型
	CreateTime  int64  `gorm:"create_time"`  //创建时间
}

func (t *UploadFileLog) TableName() string {
	return UploadFileLogTableName
}

func NewUploadFileLogModel(db *gorm.DB) *UploadFileLogModel {
	return &UploadFileLogModel{db: db}
}

type UploadFileLogModel struct {
	db *gorm.DB
}

// 插入一条记录
func (m *UploadFileLogModel) Insert(o *UploadFileLog) error {
	o.CreateTime = time.Now().Unix()
	result := m.db.Create(o)
	return result.Error
}

// 确认是否存在记录
func (m *UploadFileLogModel) FindOne(fileName string) (UploadFileLog, error) {
	var data UploadFileLog
	result := m.db.Model(&UploadFileLog{}).Where("file_name=?", fileName).Find(&data)
	return data, result.Error
}
