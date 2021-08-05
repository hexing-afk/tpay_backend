package model

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
	"tpay_backend/utils"
)

const TransferBatchOrderTableName = "transfer_batch_order"

const (
	// 处理状态
	BatchStatusInit    = 1 // 初始化
	BatchStatusSuccess = 2 // 成功
	BatchStatusFailed  = 3 // 失败
	BatchStatusPending = 4 // 处理中

	// 是否已全部生成订单
	GenerateAllDone   = 1 // 完成
	GenerateAllUndone = 2 // 未完成
)

type TransferBatchOrder struct {
	Id               int64         `gorm:"id"`
	BatchNo          string        `gorm:"batch_no"`     // 批量号
	TotalNumber      int64         `gorm:"total_number"` // 订单总笔数
	TotalAmount      int64         `gorm:"total_amount"` // 订单总金额
	Status           int64         `gorm:"status"`       // 批量状态  1 初始化，2-成功，3-失败, 4-处理中
	MerchantId       int64         `gorm:"merchant_id"`  // 商户id
	CreateTime       int64         `gorm:"create_time"`  // 创建时间
	FinishTime       int64         `gorm:"finish_time"`  // 完成时间
	FileContent      string        `gorm:"file_content"` // 文件内容json
	GenerateAll      int64         `gorm:"generate_all"` // 是否已全部生成订单 1-完成，2-未完成
	FileContentSlice []FileContent `gorm:"-"`            // 文件内容
}

type FileContent struct {
	Row            string `json:"row"`              //行号
	ChannelCode    string `json:"channel_code"`     //平台代付通道code
	AccountName    string `json:"account_name"`     //收款人姓名
	CardNumber     string `json:"card_number"`      //收款卡号
	BankName       string `json:"bank_name"`        //银行名称
	BankBranchName string `json:"bank_branch_name"` //支行名称
	Remark         string `json:"remark"`           //备注
	Amount         int64  `json:"amount"`           //金额
}

func (t *TransferBatchOrder) TableName() string {
	return TransferBatchOrderTableName
}

type TransferBatchOrderModel struct {
	db *gorm.DB
}

func NewTransferBatchOrderModel(db *gorm.DB) *TransferBatchOrderModel {
	return &TransferBatchOrderModel{db: db}
}

func (m *TransferBatchOrderModel) Insert(data *TransferBatchOrder) error {
	data.CreateTime = time.Now().Unix()
	result := m.db.Create(data)
	return result.Error
}

// 生成订单号
func (t *TransferBatchOrderModel) GenerateOrderNo() string {
	// 1位+16位+5位 = 22位
	return fmt.Sprintf("%s%d%d",
		PayOrderNoPrefix,
		time.Now().UnixNano()/1000,
		utils.RandInt64(10000, 99999),
	)
}

func (m *TransferBatchOrderModel) FindOne(id int64) (*TransferBatchOrder, error) {
	var obj = &TransferBatchOrder{}
	result := m.db.Model(&TransferBatchOrder{Id: id}).First(obj)
	if result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (m *TransferBatchOrderModel) FindByBatchNo(merchantId int64, batchNo string) (*TransferBatchOrder, error) {
	var obj = &TransferBatchOrder{}
	result := m.db.Model(&TransferBatchOrder{}).Where("merchant_id=? and batch_no=?", merchantId, batchNo).First(obj)
	if result.Error != nil {
		return nil, result.Error
	}

	if obj.FileContent != "" {
		var fileContentSlice []FileContent
		if err := json.Unmarshal([]byte(obj.FileContent), &fileContentSlice); err != nil {
			return nil, err
		}

		obj.FileContentSlice = fileContentSlice
	}

	return obj, nil
}

func (m *TransferBatchOrderModel) UpdateStatus(id int64, status, generateAll int64) error {
	result := m.db.Model(&TransferBatchOrder{Id: id}).
		Update("status", status).
		Update("generate_all", generateAll)
	return result.Error
}

func (m *TransferBatchOrderModel) FindGenerateAllUndone() ([]*TransferBatchOrder, error) {
	var dataList []*TransferBatchOrder
	result := m.db.Table(TransferBatchOrderTableName+" batch").
		Joins("LEFT JOIN "+MerchantTableName+" m ON m.merchant_no=batch.merchant_no").
		Select("batch.batch_no, m.merchant_no, m.balance").
		Where("batch.status=? AND batch.generate_all=? AND m.status=?", BatchStatusPending, GenerateAllUndone, MerchantStatusEnable).
		Scan(&dataList)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, v := range dataList {
		var fileContent []FileContent
		if err := json.Unmarshal([]byte(v.FileContent), &fileContent); err != nil {
			continue
		}
		v.FileContentSlice = fileContent
	}

	return dataList, nil
}
