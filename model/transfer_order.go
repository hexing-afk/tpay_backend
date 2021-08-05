package model

import (
	"bytes"
	"fmt"
	"time"
	"tpay_backend/utils"

	"gorm.io/gorm"
)

const TransferOrderTableName = "transfer_order"

const (
	// 代付订单状态
	TransferOrderStatusPending = 1 // 待支付
	TransferOrderStatusPaid    = 2 // 已支付
	TransferOrderStatusFail    = 3 // 支付失败

	// 代付订单异步通知状态
	TransferNotifyStatusNot       = 0 // 未通知
	TransferNotifyStatusSuccess   = 1 // 成功
	TransferNotifyStatusNotifying = 2 // 通知进行中
	TransferNotifyStatusTimeOut   = 3 // 超时

	// 代付订单来源
	TransferOrderSourceInterface       = 1 // 接口
	TransferOrderSourceWithdrawAllot   = 2 // 平台提现派单
	TransferOrderSourceMerchantPayment = 3 // 商户后台付款

	// 模式
	TransferModeTest = "test" // 测试
	TransferModePro  = "pro"  // 生产
)

type TransferOrder struct {
	Id                 int64   `gorm:"id"`
	OrderNo            string  `gorm:"order_no"`             // 平台订单号
	MerchantOrderNo    string  `gorm:"merchant_order_no"`    // 商户订单号
	UpstreamOrderNo    string  `gorm:"upstream_order_no"`    // 上游订单号
	MerchantNo         string  `gorm:"merchant_no"`          // 商户编号
	ReqAmount          int64   `gorm:"req_amount"`           // 订单请求金额
	MerchantFee        int64   `gorm:"merchant_fee"`         // 商户手续费
	DecreaseAmount     int64   `gorm:"decrease_amount"`      // 账户扣除的金额
	UpstreamAmount     int64   `gorm:"upstream_amount"`      // 请求上游的金额
	UpstreamFee        int64   `gorm:"upstream_fee"`         // 上游手续费
	Currency           string  `gorm:"currency"`             // 币种
	OrderStatus        int64   `gorm:"order_status"`         // 订单状态
	CreateTime         int64   `gorm:"create_time"`          // 创建时间
	PlatformChannelId  int64   `gorm:"platform_channel_id"`  // 平台通道id
	UpstreamChannelId  int64   `gorm:"upstream_channel_id"`  // 上游通道id
	UpdateTime         int64   `gorm:"update_time"`          // 更新时间
	NotifyUrl          string  `gorm:"notify_url"`           // 异步通知地址
	ReturnUrl          string  `gorm:"return_url"`           // 同步跳转地址
	BankName           string  `gorm:"bank_name"`            // 银行名称
	AccountName        string  `gorm:"account_name"`         // 银行卡开户名
	CardNumber         string  `gorm:"card_number"`          // 收款卡号
	BranchName         string  `gorm:"branch_name"`          // 支行名称
	NotifyStatus       int64   `gorm:"notify_status"`        // 异步通知状态(0未通知,1成功,2通知进行中,3超时)
	NotifyFailTimes    int64   `gorm:"notify_fail_times"`    // 通知失败次数
	NextNotifyTime     int64   `gorm:"next_notify_time"`     // 下次通知时间
	PayeeRealAmount    int64   `gorm:"payee_real_amount"`    // 收款方实际到账金额
	FeeDeductType      int64   `gorm:"fee_deduct_type"`      // 手续费扣款方式(1内扣,2外扣)
	UpstreamFailReason string  `gorm:"upstream_fail_reason"` // 上游失败原因
	OrderSource        int64   `gorm:"order_source"`         // 订单来源：1-接口; 2-平台提现派单；3-商户后台付款
	Remark             string  `gorm:"remark"`               // 付款备注
	BankCode           string  `gorm:"bank_code"`            // 银行代码(ifsc_code)
	AreaId             int64   `gorm:"area_id"`
	Mode               string  `gorm:"mode"`                // 模式：test|pro(测试|生产)
	MerchantRate       float64 `gorm:"merchant_rate"`       // 商户费率
	MerchantSingleFee  int64   `gorm:"merchant_single_fee"` // 商户单笔手续费
	BatchNo            string  `gorm:"batch_no"`            // 批量付款批次号
	BatchRowNo         string  `gorm:"batch_row_no"`        // 批量付款批次行号
}

func (t *TransferOrder) TableName() string {
	return TransferOrderTableName
}

func NewTransferOrderModel(db *gorm.DB) *TransferOrderModel {
	return &TransferOrderModel{db: db}
}

type TransferOrderModel struct {
	db *gorm.DB
}

// 生成订单号
func (t *TransferOrder) GenerateOrderNo(mode string) string {
	orderNo := fmt.Sprintf("%d%d%d",
		TransferOrderNoPrefix,
		time.Now().UnixNano()/1000,
		utils.RandInt64(10000, 99999),
	)
	switch mode {
	case TransferModeTest:
		// 1位+16位+5位 = 22位
		return TransferModeTest + orderNo
	case TransferModePro:
		fallthrough
	default:
		return orderNo
	}
}

func (m *TransferOrderModel) Insert(data *TransferOrder) error {
	data.CreateTime = time.Now().Unix()
	result := m.db.Create(data)
	return result.Error
}

func (m *TransferOrderModel) UpdateUpstreamOrderInfo(id int64, data TransferOrder) error {
	setMap := map[string]interface{}{
		"update_time":       time.Now().Unix(),
		"upstream_order_no": data.UpstreamOrderNo,
	}

	if data.OrderStatus != 0 {
		setMap["order_status"] = data.OrderStatus
	}

	result := m.db.Model(&TransferOrder{}).Where("id=?", id).Updates(setMap)
	return result.Error
}

// 修改订单状态
func (m *TransferOrderModel) UpdateOrderStatus(id, status int64) error {
	result := m.db.Model(&TransferOrder{}).Where("id=?", id).Update("order_status", status)
	return result.Error
}

// 修改订单为已支付
func (m *TransferOrderModel) UpdateOrderPaidById(id int64) error {
	result := m.db.Model(&TransferOrder{}).Where("id=?", id).Update("order_status", TransferOrderStatusPaid)
	return result.Error
}

// 修改订单为失败
func (m *TransferOrderModel) UpdateOrderFailById(id int64, failReason string) error {
	setMap := map[string]interface{}{
		"order_status":         TransferOrderStatusFail,
		"upstream_fail_reason": failReason,
	}
	result := m.db.Model(&TransferOrder{}).Where("id=?", id).Updates(setMap)
	return result.Error
}

type TransferNotifyInfo struct {
	NotifyStatus    int64 `gorm:"notify_status"`     // 异步通知状态(0未通知,1成功,2通知进行中,3超时)
	NotifyFailTimes int64 `gorm:"notify_fail_times"` // 通知失败次数
	NextNotifyTime  int64 `gorm:"next_notify_time"`  // 下次通知时间
}

// 修改订单-通知信息
func (m *TransferOrderModel) UpdateNotify(orderId int64, data TransferNotifyInfo) error {
	setMap := map[string]interface{}{
		"notify_status":     data.NotifyStatus,
		"notify_fail_times": data.NotifyFailTimes,
		"next_notify_time":  data.NextNotifyTime,
	}
	result := m.db.Model(TransferOrder{}).Where("id = ?", orderId).Updates(&setMap)
	return result.Error
}

func (m *TransferOrderModel) MerchantOrderNoExist(merchantNo, merchantOrderNo string) (bool, error) {
	var cnt int64
	result := m.db.Model(&TransferOrder{}).Where("merchant_no=? and merchant_order_no=?", merchantNo, merchantOrderNo).Count(&cnt)
	if result.Error != nil {
		return false, result.Error
	}

	return cnt > 0, nil
}

func (m *TransferOrderModel) FindByOrderNo(orderNo string) (*TransferOrder, error) {
	var o = &TransferOrder{}
	result := m.db.Model(o).Where("order_no=?", orderNo).First(o)
	return o, result.Error
}

func (m *TransferOrderModel) FindByMerchantId(merchantId int64, orderNo string) (*TransferOrder, error) {
	var o = &TransferOrder{}
	result := m.db.Table(TransferOrderTableName+" o").
		Select("o.*").
		Joins("left join "+MerchantTableName+" m on m.merchant_no=o.merchant_no").
		Where("m.id=? and o.order_no=?", merchantId, orderNo).
		Find(o)
	return o, result.Error
}

type FindByOrderNoAndMerchantNoData struct {
	TransferOrder
	PlatformChannelName string `gorm:"platform_channel_name"` // 银行代码(ifsc_code)
}

func (m *TransferOrderModel) FindByOrderNoAndMerchantNo(orderNo, merchantNo string) (FindByOrderNoAndMerchantNoData, error) {
	var o FindByOrderNoAndMerchantNoData
	selectField := "t.*, plat.channel_name AS platform_channel_name "
	platformCh := fmt.Sprintf("left join %s plat on plat.id = t.platform_channel_id", PlatformChannelTableName)
	result := m.db.Table(TransferOrderTableName+" t").
		Select(selectField).
		Joins(platformCh).
		Where("t.order_no=? and t.merchant_no = ? ", orderNo, merchantNo).Scan(&o)
	return o, result.Error
}

func (m *TransferOrderModel) FindByMchOrderNo(merchantNo, merchantOrderNo string) (*TransferOrder, error) {
	var o = &TransferOrder{}
	result := m.db.Model(o).Where("merchant_no=? and merchant_order_no=?", merchantNo, merchantOrderNo).First(o)
	return o, result.Error
}

type FindTransferOrderList struct {
	Page              int64
	PageSize          int64
	OrderNo           string
	MerchantOrderNo   string
	UpstreamOrderNo   string
	MerchantName      string
	Currency          string
	PlatformChannelId int64
	StartCreateTime   int64
	EndCreateTime     int64
	OrderStatus       int64
	OrderSourceList   []int64
	MerchantNo        string
	OrderType         string
}

type TransferOrderList struct {
	TransferOrder
	MerchantName        string `gorm:"merchant_name"`
	PlatformChannelName string `gorm:"platform_channel_name"`
	UpstreamName        string `gorm:"upstream_name"`
}

func (m *TransferOrderModel) FindList(f FindTransferOrderList) ([]*TransferOrderList, int64, error) {
	var (
		order      = fmt.Sprintf("%s o", TransferOrderTableName)
		merchant   = fmt.Sprintf("left join %s m on m.merchant_no = o.merchant_no", MerchantTableName)
		platformCh = fmt.Sprintf("left join %s plat on plat.id = o.platform_channel_id", PlatformChannelTableName)
		upstreamCh = fmt.Sprintf("left join %s upst on upst.id = o.upstream_channel_id", UpstreamChannelTableName)
		upstream   = fmt.Sprintf("left join %s up on up.id = upst.upstream_id", UpstreamTableName)

		selectField = "o.order_no, o.merchant_order_no, o.currency, o.req_amount, o.merchant_fee, " +
			"o.payee_real_amount, o.upstream_fee, o.order_status, o.create_time, o.update_time," +
			"o.bank_name, o.account_name, o.card_number, o.branch_name, o.order_source," +
			"o.upstream_order_no, o.remark, " +
			"m.username AS merchant_name, plat.channel_name as platform_channel_name, up.upstream_name "

		whereBuffer = bytes.NewBufferString(" 1=1 ")
		args        []interface{}
	)

	if f.OrderNo != "" {
		whereBuffer.WriteString("and o.order_no like ? ")
		args = append(args, "%"+f.OrderNo+"%")
	}

	if f.MerchantOrderNo != "" {
		whereBuffer.WriteString("and o.merchant_order_no like ? ")
		args = append(args, "%"+f.MerchantOrderNo+"%")
	}

	if f.UpstreamOrderNo != "" {
		whereBuffer.WriteString("and o.upstream_order_no like ? ")
		args = append(args, "%"+f.UpstreamOrderNo+"%")
	}

	if f.MerchantNo != "" {
		whereBuffer.WriteString("and m.merchant_no = ? ")
		args = append(args, f.MerchantNo)
	}

	if f.MerchantName != "" {
		whereBuffer.WriteString("and m.username like ? ")
		args = append(args, "%"+f.MerchantName+"%")
	}

	if f.Currency != "" {
		whereBuffer.WriteString("and o.currency = ? ")
		args = append(args, f.Currency)
	}

	if f.PlatformChannelId != 0 {
		whereBuffer.WriteString("and o.platform_channel_id = ? ")
		args = append(args, f.PlatformChannelId)
	}

	if f.StartCreateTime != 0 {
		whereBuffer.WriteString("and o.create_time >= ? ")
		args = append(args, f.StartCreateTime)
	}

	if f.EndCreateTime != 0 {
		whereBuffer.WriteString("and o.create_time <= ? ")
		args = append(args, f.EndCreateTime)
	}

	if f.OrderStatus != 0 {
		whereBuffer.WriteString("and o.order_status = ? ")
		args = append(args, f.OrderStatus)
	}

	if f.OrderType != "" {
		whereBuffer.WriteString("and o.mode = ? ")
		args = append(args, f.OrderType)
	}

	if f.OrderSourceList != nil {
		whereBuffer.WriteString("and o.order_source in ? ")
		args = append(args, f.OrderSourceList)

	}

	if f.Page == 0 {
		f.Page = 1
	}
	if f.PageSize == 0 {
		f.PageSize = 10
	}

	var total int64
	ret := m.db.Table(order).Joins(merchant).Joins(platformCh).Joins(upstreamCh).Joins(upstream).Where(whereBuffer.String(), args...).Count(&total)
	if ret.Error != nil {
		return nil, 0, ret.Error
	}

	if total == 0 {
		return nil, 0, nil
	}

	whereBuffer.WriteString(" order by o.create_time desc limit ? offset ?")
	args = append(args, f.PageSize, (f.Page-1)*f.PageSize)

	var list []*TransferOrderList
	ret = m.db.Table(order).Select(selectField).Joins(merchant).Joins(platformCh).Joins(upstreamCh).Joins(upstream).Where(whereBuffer.String(), args...).Scan(&list)
	if ret.Error != nil {
		return nil, 0, ret.Error
	}

	return list, total, nil
}

// 查询异步通知遗漏的订单
func (m *TransferOrderModel) FindNotifyOmissionOrderNo(payTime int64) ([]string, error) {
	orderStatus := []int64{TransferOrderStatusPaid, TransferOrderStatusFail}

	selectField := "order_no"
	whereStr := "order_status in ? and notify_status= ? and update_time < ?"

	var orderNos []string
	result := m.db.Model(&TransferOrder{}).Select(selectField).Where(whereStr, orderStatus, TransferNotifyStatusNot, payTime).Find(&orderNos)
	if result.Error != nil {
		return nil, result.Error
	}

	return orderNos, nil
}

// 查询异步通知中断的订单
func (m *TransferOrderModel) FindNotifyBreakOrderNo(nextTime int64) ([]string, error) {
	orderStatus := []int64{TransferOrderStatusPaid, TransferOrderStatusFail}

	selectField := "order_no"
	whereStr := "order_status in ? and notify_status = ? and next_notify_time < ?"

	var orderNos []string
	result := m.db.Model(&TransferOrder{}).Select(selectField).Where(whereStr, orderStatus, TransferNotifyStatusNotifying, nextTime).Find(&orderNos)
	if result.Error != nil {
		return nil, result.Error
	}

	return orderNos, nil
}

type TransferOrderDetail struct {
	TransferOrder
	MerchantName        string `gorm:"merchant_name"`
	PlatformChannelName string `gorm:"platform_channel_name"`
	UpstreamName        string `gorm:"upstream_name"`
}

func (m *TransferOrderModel) FindDetail(orderNo string) (*TransferOrderDetail, error) {
	var (
		order    = fmt.Sprintf("%s o", TransferOrderTableName)
		merchant = fmt.Sprintf("left join %s m on m.merchant_no = o.merchant_no", MerchantTableName)

		selectField = "o.order_no, o.merchant_order_no, o.upstream_order_no, o.merchant_no, o.req_amount," +
			"o.decrease_amount, o.merchant_fee, o.upstream_fee, o.payee_real_amount, o.fee_deduct_type," +
			"o.upstream_amount, o.currency, o.create_time, o.update_time, o.notify_url, " +
			"o.bank_name, o.card_number, o.account_name, o.branch_name, o.notify_status, " +
			"m.username AS merchant_name "
	)

	var data TransferOrderDetail
	ret := m.db.Table(order).Select(selectField).Joins(merchant).Where(" o.order_no = ? ", orderNo).Scan(&data)
	if ret.Error != nil {
		return nil, ret.Error
	}

	return &data, nil
}

type FindTransferExportData struct {
	MerchantId        int64
	OrderNo           string
	MerchantOrderNo   string
	UpstreamOrderNo   string
	MerchantName      string
	Currency          string
	PlatformChannelId int64
	StartCreateTime   int64
	EndCreateTime     int64
	OrderStatus       int64
	OrderSourceList   []int64
	MerchantNo        string
	OrderType         string
}

type TransferExportData struct {
	OrderList           []TransferOrderList
	Total               int64
	TotalReqAmount      int64
	TotalFee            int64
	TotalIncreaseAmount int64
}

// 查询导出数据
func (m *TransferOrderModel) FindExportData(f FindTransferExportData) (*TransferExportData, error) {
	var (
		order      = fmt.Sprintf("%s o", TransferOrderTableName)
		merchant   = fmt.Sprintf("left join %s m on m.merchant_no = o.merchant_no", MerchantTableName)
		platformCh = fmt.Sprintf("left join %s plat on plat.id = o.platform_channel_id", PlatformChannelTableName)

		selectField = "o.id, o.order_no, o.merchant_order_no, o.currency, o.req_amount, " +
			"o.merchant_rate, o.merchant_fee, o.merchant_single_fee, o.payee_real_amount, o.order_status, " +
			"o.create_time, o.update_time," +
			"m.username AS merchant_name, plat.channel_name as platform_channel_name "

		whereBuffer = bytes.NewBufferString(" 1=1 ")
		args        []interface{}
	)

	if f.MerchantId != 0 {
		whereBuffer.WriteString("and m.id = ? ")
		args = append(args, f.MerchantId)
	}

	if f.OrderNo != "" {
		whereBuffer.WriteString("and o.order_no like ? ")
		args = append(args, "%"+f.OrderNo+"%")
	}

	if f.MerchantOrderNo != "" {
		whereBuffer.WriteString("and o.merchant_order_no like ? ")
		args = append(args, "%"+f.MerchantOrderNo+"%")
	}

	if f.UpstreamOrderNo != "" {
		whereBuffer.WriteString("and o.upstream_order_no like ? ")
		args = append(args, "%"+f.UpstreamOrderNo+"%")
	}

	if f.MerchantNo != "" {
		whereBuffer.WriteString("and m.merchant_no = ? ")
		args = append(args, f.MerchantNo)
	}

	if f.MerchantName != "" {
		whereBuffer.WriteString("and m.username like ? ")
		args = append(args, "%"+f.MerchantName+"%")
	}

	if f.Currency != "" {
		whereBuffer.WriteString("and o.currency = ? ")
		args = append(args, f.Currency)
	}

	if f.PlatformChannelId != 0 {
		whereBuffer.WriteString("and o.platform_channel_id = ? ")
		args = append(args, f.PlatformChannelId)
	}

	if f.StartCreateTime != 0 {
		whereBuffer.WriteString("and o.create_time >= ? ")
		args = append(args, f.StartCreateTime)
	}

	if f.EndCreateTime != 0 {
		whereBuffer.WriteString("and o.create_time <= ? ")
		args = append(args, f.EndCreateTime)
	}

	if f.OrderStatus != 0 {
		whereBuffer.WriteString("and o.order_status = ? ")
		args = append(args, f.OrderStatus)
	}

	if f.OrderType != "" {
		whereBuffer.WriteString("and o.mode = ? ")
		args = append(args, f.OrderType)
	}

	if f.OrderSourceList != nil {
		whereBuffer.WriteString("and o.order_source in ? ")
		args = append(args, f.OrderSourceList)

	}

	var total, totalReqAmount, totalFee, totalIncreaseAmount int64

	result := m.db.Table(order).Joins(merchant).Joins(platformCh).Where(whereBuffer.String(), args...).Count(&total)
	if result.Error != nil {
		return nil, result.Error
	}

	if total == 0 {
		return nil, nil
	}

	// 请求金额总数
	result = m.db.Table(order).Select("IFNULL(sum(o.req_amount), 0)").Joins(merchant).Joins(platformCh).Where(whereBuffer.String(), args...).Scan(&totalReqAmount)
	if result.Error != nil {
		return nil, result.Error
	}

	// 商户手续费总数
	result = m.db.Table(order).Select("IFNULL(sum(o.merchant_fee), 0)").Joins(merchant).Joins(platformCh).Where(whereBuffer.String(), args...).Scan(&totalFee)
	if result.Error != nil {
		return nil, result.Error
	}

	// 商户实际入账金额总数
	result = m.db.Table(order).Select("IFNULL(sum(o.payee_real_amount), 0)").Joins(merchant).Joins(platformCh).
		Where(whereBuffer.String()+fmt.Sprintf("and o.order_status = %v ", TransferOrderStatusPaid), args...).Scan(&totalIncreaseAmount)
	if result.Error != nil {
		return nil, result.Error
	}

	whereBuffer.WriteString(" order by o.create_time asc ")

	var list []TransferOrderList
	result = m.db.Table(order).Select(selectField).Joins(merchant).Joins(platformCh).Where(whereBuffer.String(), args...).Scan(&list)
	if result.Error != nil {
		return nil, result.Error
	}

	return &TransferExportData{
		OrderList:           list,
		Total:               total,
		TotalReqAmount:      totalReqAmount,
		TotalFee:            totalFee,
		TotalIncreaseAmount: totalIncreaseAmount,
	}, nil
}

func (m *TransferOrderModel) FindByBatchNoAndRowNo(merchantNo string, batchNo string, rowNos []string) ([]string, error) {
	var rowNo []string
	result := m.db.Model(&TransferOrder{}).Select("batch_row_no").
		Where("merchant_no=? and batch_no=? and batch_row_no in ?", merchantNo, batchNo, rowNos).Scan(&rowNo)

	return rowNo, result.Error
}

func (m *TransferOrderModel) CountByBatchNo(merchantNo string, batchNo string) (int64, error) {
	var total int64
	result := m.db.Model(&TransferOrder{}).Where("merchant_no=? and batch_no=?", merchantNo, batchNo).Count(&total)

	return total, result.Error
}
