package model

import (
	"bytes"
	"fmt"
	"time"
	"tpay_backend/utils"

	"gorm.io/gorm"
)

const PayOrderTableName = "pay_order"

const (
	// 代收订单状态
	PayOrderStatusPending = 1 // 待支付
	PayOrderStatusPaid    = 2 // 已支付
	PayOrderStatusFail    = 3 // 支付失败

	// 代收订单异步通知状态
	PayNotifyStatusNot       = 0 // 未通知
	PayNotifyStatusSuccess   = 1 // 成功
	PayNotifyStatusNotifying = 2 // 通知进行中
	PayNotifyStatusTimeOut   = 3 // 超时

	// 模式
	PayModeTest = "test" // 测试
	PayModePro  = "pro"  // 生产
)

type PayOrder struct {
	Id                int64   `gorm:"id"`
	MerchantNo        string  `gorm:"merchant_no"`         // 商户编号
	OrderNo           string  `gorm:"order_no"`            // 平台单号
	MerchantOrderNo   string  `gorm:"merchant_order_no"`   // 商户订单号(即下游订单号)
	UpstreamOrderNo   string  `gorm:"upstream_order_no"`   // 上游订单号
	ReqAmount         int64   `gorm:"req_amount"`          // 订单请求金额
	IncreaseAmount    int64   `gorm:"increase_amount"`     // 账户增加的金额
	MerchantFee       int64   `gorm:"merchant_fee"`        // 商户手续费
	UpstreamAmount    int64   `gorm:"upstream_amount"`     // 请求上游的金额
	UpstreamFee       int64   `gorm:"upstream_fee"`        // 上游手续费
	OrderStatus       int64   `gorm:"order_status"`        // 订单状态
	Currency          string  `gorm:"currency"`            // 币种
	CreateTime        int64   `gorm:"create_time"`         // 创建时间
	UpdateTime        int64   `gorm:"update_time"`         // 更新时间
	NotifyUrl         string  `gorm:"notify_url"`          // 异步通知url
	ReturnUrl         string  `gorm:"return_url"`          // 同步跳转url
	PlatformChannelId int64   `gorm:"platform_channel_id"` // 平台通道id
	UpstreamChannelId int64   `gorm:"upstream_channel_id"` // 上游通道id
	NotifyStatus      int64   `gorm:"notify_status"`       // 异步通知状态(0未通知,1成功,2通知进行中,3超时)
	NotifyFailTimes   int64   `gorm:"notify_fail_times"`   // 通知失败次数
	NextNotifyTime    int64   `gorm:"next_notify_time"`    // 下次通知时间
	Subject           string  `gorm:"subject"`             // 商品的标题/交易标题/订单标题/订单关键字等
	AreaId            int64   `gorm:"area_id"`             // 地区
	Mode              string  `gorm:"mode"`                // 模式：test|pro(测试|生产)
	MerchantRate      float64 `gorm:"merchant_rate"`       // 商户费率
	MerchantSingleFee int64   `gorm:"merchant_single_fee"` // 商户单笔手续费
	PaymentAmount     int64   `gorm:"payment_amount"`      // 实际支付金额
}

func (t *PayOrder) TableName() string {
	return PayOrderTableName
}

// 生成订单号
func (t *PayOrder) GenerateOrderNo(mode string) string {
	orderNo := fmt.Sprintf("%d%d%d",
		PayOrderNoPrefix,
		time.Now().UnixNano()/1000,
		utils.RandInt64(10000, 99999),
	)
	switch mode {
	case PayModeTest:
		return PayModeTest + orderNo
	case PayModePro:
		fallthrough
	default:
		return orderNo
	}
}

func NewPayOrderModel(db *gorm.DB) *PayOrderModel {
	return &PayOrderModel{db: db}
}

type PayOrderModel struct {
	db *gorm.DB
}

// 插入一条记录
func (m *PayOrderModel) Insert(o *PayOrder) error {
	o.CreateTime = time.Now().Unix()
	o.UpdateTime = o.CreateTime
	result := m.db.Create(o)
	return result.Error
}

// 商户订单号是否已经存在
func (m *PayOrderModel) MerchantOrderNoExist(merchantNo, merchantOrderNo string) (bool, error) {
	var count int64

	result := m.db.Model(&PayOrder{}).Where("merchant_no=? AND merchant_order_no=? ", merchantNo, merchantOrderNo).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

// 更新上游订单相关信息
func (m *PayOrderModel) UpdateUpstreamInfo(orderId int64, upOrderNo string) error {
	updates := map[string]interface{}{
		"upstream_order_no": upOrderNo,
	}
	result := m.db.Model(PayOrder{}).Where("id = ?", orderId).Updates(updates)
	return result.Error
}

// 修改订单状态
func (m *PayOrderModel) UpdateOrderStatus(orderId, status, paymentAmount int64) error {
	sqlStr := "UPDATE pay_order SET order_status=?, payment_amount=? WHERE id=?"

	result := m.db.Exec(sqlStr, status, paymentAmount, orderId)
	return result.Error
}

// 修改订单-已支付
func (m *PayOrderModel) UpdateOrderPaid(orderId int64, order *PayOrder) error {
	setMap := map[string]interface{}{
		"order_status":    PayOrderStatusPaid,
		"update_time":     time.Now().Unix(),
		"payment_amount":  order.PaymentAmount,
		"increase_amount": order.IncreaseAmount,
		"merchant_fee":    order.MerchantFee,
	}

	result := m.db.Model(PayOrder{}).Where("id = ?", orderId).Updates(setMap)
	return result.Error
}

// 修改订单-超时
func (m *PayOrderModel) UpdateOrderFail(orderId int64) error {
	setMap := map[string]interface{}{
		"order_status": PayOrderStatusFail,
		"update_time":  time.Now().Unix(),
	}

	result := m.db.Model(PayOrder{}).Where("id = ?", orderId).Updates(setMap)
	return result.Error
}

type PayNotifyInfo struct {
	NotifyStatus    int64 `gorm:"notify_status"`     // 异步通知状态(0未通知,1成功,2通知进行中,3超时)
	NotifyFailTimes int64 `gorm:"notify_fail_times"` // 通知失败次数
	NextNotifyTime  int64 `gorm:"next_notify_time"`  // 下次通知时间
}

// 修改订单-通知信息
func (m *PayOrderModel) UpdateNotify(orderId int64, data PayNotifyInfo) error {
	setMap := map[string]interface{}{
		"notify_status":     data.NotifyStatus,
		"notify_fail_times": data.NotifyFailTimes,
		"next_notify_time":  data.NextNotifyTime,
	}
	result := m.db.Model(PayOrder{}).Where("id = ?", orderId).Updates(&setMap)
	return result.Error
}

// 根据订单号查询
func (m *PayOrderModel) FindOneByOrderNo(orderNo string) (*PayOrder, error) {
	a := &PayOrder{}
	result := m.db.Model(a).Where("order_no=?", orderNo).Limit(1).First(&a)
	return a, result.Error
}

// 根据商户编号-订单号查询
func (m *PayOrderModel) FindMerchantOrder(merchantNo, mchOrderNo string) (*PayOrder, error) {
	a := &PayOrder{}
	result := m.db.Model(a).Where("merchant_no=? AND merchant_order_no=?", merchantNo, mchOrderNo).Limit(1).First(&a)
	return a, result.Error
}

// 根据商户id-订单号查询
func (m *PayOrderModel) FindOrderByMerchantId(merchantId int64, orderNo string) (*PayOrder, error) {
	a := &PayOrder{}

	result := m.db.Table(PayOrderTableName+" o").
		Select("o.*").
		Joins("left join "+MerchantTableName+" m on m.merchant_no=o.merchant_no").
		Where("m.id=? and o.order_no=?", merchantId, orderNo).
		Find(a)
	return a, result.Error
}

type FindPayOrderList struct {
	Page                int64
	PageSize            int64
	OrderNo             string
	MerchantOrderNo     string
	UpstreamOrderNo     string
	MerchantId          int64
	MerchantName        string
	Currency            string
	PlatformChannelId   int64
	StartCreateTime     int64
	EndCreateTime       int64
	OrderStatus         int64
	PlatformChannelName string
	OrderType           string
}

type PayOrderList struct {
	PayOrder
	MerchantName        string `gorm:"merchant_name"`
	PlatformChannelName string `gorm:"platform_channel_name"`
	UpstreamName        string `gorm:"upstream_name"`
}

// 查询列表
func (m *PayOrderModel) FindList(f FindPayOrderList) ([]*PayOrderList, int64, error) {
	var (
		payOrder   = fmt.Sprintf("%s pay", PayOrderTableName)
		merchant   = fmt.Sprintf("left join %s m on m.merchant_no = pay.merchant_no", MerchantTableName)
		platformCh = fmt.Sprintf("left join %s plat on plat.id = pay.platform_channel_id", PlatformChannelTableName)
		upstreamCh = fmt.Sprintf("left join %s upstc on upstc.id = pay.upstream_channel_id", UpstreamChannelTableName)
		upstream   = fmt.Sprintf("left join %s up on up.id = upstc.upstream_id", UpstreamTableName)

		selectField = "pay.order_no, pay.merchant_order_no, pay.upstream_order_no, pay.currency, pay.req_amount, " +
			"pay.increase_amount, pay.merchant_fee, pay.order_status, pay.create_time, pay.update_time, pay.payment_amount, " +
			"m.username as merchant_name, plat.channel_name as platform_channel_name, up.upstream_name "

		whereBuffet = bytes.NewBufferString(" 1=1 ")
		args        []interface{}
	)

	if f.OrderNo != "" {
		whereBuffet.WriteString("and pay.order_no like ? ")
		args = append(args, "%"+f.OrderNo+"%")
	}

	if f.MerchantOrderNo != "" {
		whereBuffet.WriteString("and pay.merchant_order_no like ? ")
		args = append(args, "%"+f.MerchantOrderNo+"%")
	}

	if f.UpstreamOrderNo != "" {
		whereBuffet.WriteString("and pay.upstream_order_no like ? ")
		args = append(args, "%"+f.UpstreamOrderNo+"%")
	}

	if f.MerchantId != 0 {
		whereBuffet.WriteString("and m.id = ? ")
		args = append(args, f.MerchantId)
	}

	if f.MerchantName != "" {
		whereBuffet.WriteString("and m.username like ? ")
		args = append(args, "%"+f.MerchantName+"%")
	}

	if f.Currency != "" {
		whereBuffet.WriteString("and pay.currency = ? ")
		args = append(args, f.Currency)
	}

	if f.PlatformChannelId != 0 {
		whereBuffet.WriteString("and pay.platform_channel_id = ? ")
		args = append(args, f.PlatformChannelId)
	}

	if f.PlatformChannelName != "" {
		whereBuffet.WriteString("and plat.channel_name like ? ")
		args = append(args, "%"+f.PlatformChannelName+"%")
	}

	if f.OrderStatus != 0 {
		whereBuffet.WriteString("and pay.order_status = ? ")
		args = append(args, f.OrderStatus)
	}

	if f.StartCreateTime != 0 {
		whereBuffet.WriteString("and pay.create_time >= ? ")
		args = append(args, f.StartCreateTime)
	}

	if f.EndCreateTime != 0 {
		whereBuffet.WriteString("and pay.create_time <= ? ")
		args = append(args, f.EndCreateTime)
	}

	if f.OrderType != "" {
		whereBuffet.WriteString("and pay.mode = ? ")
		args = append(args, f.OrderType)
	}

	if f.Page == 0 {
		f.Page = 1
	}
	if f.PageSize == 0 {
		f.PageSize = 10
	}

	var total int64
	result := m.db.Table(payOrder).
		Joins(merchant).Joins(platformCh).Joins(upstreamCh).Joins(upstream).
		Where(whereBuffet.String(), args...).Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	if total == 0 {
		return nil, 0, nil
	}

	whereBuffet.WriteString("order by pay.create_time desc limit ? offset ? ")
	args = append(args, f.PageSize, (f.Page-1)*f.PageSize)

	var list []*PayOrderList
	result = m.db.Table(payOrder).
		Select(selectField).
		Joins(merchant).Joins(platformCh).Joins(upstreamCh).Joins(upstream).
		Where(whereBuffet.String(), args...).Scan(&list)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return list, total, nil
}

type PayOrderDetail struct {
	PayOrder
	MerchantName        string `gorm:"merchant_name"`
	PlatformChannelName string `gorm:"platform_channel_name"`
	UpstreamChannelName string `gorm:"upstream_channel_name"`
}

// 查询订单详情
func (m *PayOrderModel) FindOneDetailByOrderNo(orderNo string) (*PayOrderDetail, error) {
	var (
		payOrder   = fmt.Sprintf("%s pay", PayOrderTableName)
		merchant   = fmt.Sprintf("left join %s m on m.merchant_no = pay.merchant_no", MerchantTableName)
		platformCh = fmt.Sprintf("left join %s plat on plat.id = pay.platform_channel_id", PlatformChannelTableName)
		upstreamCh = fmt.Sprintf("left join %s upstc on upstc.id = pay.upstream_channel_id", UpstreamChannelTableName)

		selectField = "pay.order_no, pay.merchant_order_no, pay.upstream_order_no, pay.currency, pay.req_amount, " +
			"pay.increase_amount, pay.merchant_fee, pay.order_status, pay.create_time, pay.update_time, " +
			"pay.merchant_no, pay.upstream_fee, pay.notify_url, pay.return_url, pay.platform_channel_id," +
			"pay.upstream_channel_id, pay.notify_status, pay.subject, pay.payment_amount, " +
			"m.username as merchant_name, plat.channel_name as platform_channel_name, upstc.channel_name AS upstream_channel_name "
	)

	var data PayOrderDetail
	result := m.db.Table(payOrder).
		Select(selectField).
		Joins(merchant).Joins(platformCh).Joins(upstreamCh).
		Where(" pay.order_no = ? ", orderNo).Limit(1).Scan(&data)
	if result.Error != nil {
		return nil, result.Error
	}
	return &data, nil
}

// 查询异步通知遗漏的订单
func (m *PayOrderModel) FindNotifyOmissionOrderNo(orderStatus []int64, payTime int64) ([]string, error) {
	selectField := "order_no"
	whereStr := "order_status in ? and notify_status= ? and update_time < ?"

	var orderNos []string
	result := m.db.Model(&PayOrder{}).Select(selectField).Where(whereStr, orderStatus, PayNotifyStatusNot, payTime).Find(&orderNos)
	if result.Error != nil {
		return nil, result.Error
	}

	return orderNos, nil
}

// 查询异步通知中断的订单
func (m *PayOrderModel) FindNotifyBreakOrderNo(orderStatus []int64, nextTime int64) ([]string, error) {
	selectField := "order_no"
	whereStr := "order_status in ? and notify_status = ? and next_notify_time < ?"

	var orderNos []string
	result := m.db.Model(&PayOrder{}).Select(selectField).Where(whereStr, orderStatus, PayNotifyStatusNotifying, nextTime).Find(&orderNos)
	if result.Error != nil {
		return nil, result.Error
	}

	return orderNos, nil
}

type ToDayCountDetail struct {
	OrderNumber        int64 //今日收款总订单数
	SuccessOrderNumber int64 `gorm:"success_order_number"` //今日收款成功订单数
	SuccessAmount      int64 `gorm:"success_amount"`       //今日成功收款金额
}

/**
查询当天统计
*/
func (m *PayOrderModel) FindMerchantToDayCountByDay(merchantNo string) (*ToDayCountDetail, error) {

	var data ToDayCountDetail

	var orderNumber int64
	//今日收款总订单数
	where := ` merchant_no = ? AND DATE(FROM_UNIXTIME(create_time)) = ? `

	args := []interface{}{merchantNo, time.Now().Format("2006-01-02")}
	if err := m.db.Table(PayOrderTableName).Where(where, args...).Count(&orderNumber).Error; err != nil {
		return nil, err
	}

	//今日收款成功订单数、今日成功收款金额
	where += " AND order_status = ? "
	args = append(args, PayOrderStatusPaid)
	selectStr := "COUNT(1) AS success_order_number, IFNULL(SUM(increase_amount),0) AS success_amount "
	if err := m.db.Table(PayOrderTableName).Select(selectStr).Where(where, args...).Scan(&data).Error; err != nil {
		return nil, err
	}

	data.OrderNumber = orderNumber

	return &data, nil
}

/**

根据日期查询统计信息
dayStr格式: yyyy-MM-dd
*/
func (m *PayOrderModel) FindMerchantCountByDay(merchantNo, dayStr string) (int64, error) {

	var successAmount int64

	//今日成功收款金额
	selectStr := "IFNULL(SUM(increase_amount),0) AS success_amount "

	where := ` merchant_no = ? AND DATE(FROM_UNIXTIME(create_time)) = ?  AND order_status = ?  `
	args := []interface{}{merchantNo, dayStr, PayOrderStatusPaid}

	if err := m.db.Table(PayOrderTableName).Select(selectStr).Where(where, args...).Scan(&successAmount).Error; err != nil {
		return 0, err
	}

	return successAmount, nil
}

type FindExportData struct {
	OrderNo             string
	MerchantOrderNo     string
	UpstreamOrderNo     string
	MerchantId          int64
	MerchantName        string
	Currency            string
	PlatformChannelId   int64
	StartCreateTime     int64
	EndCreateTime       int64
	OrderStatus         int64
	PlatformChannelName string
	OrderType           string
}

type ExportData struct {
	OrderList           []PayOrderList
	Total               int64
	TotalReqAmount      int64
	TotalPayAmount      int64
	TotalMerchantFee    int64
	TotalIncreaseAmount int64
}

// 查询导出订单
func (m *PayOrderModel) FindExportData(f FindExportData) (*ExportData, error) {
	var (
		payOrder   = fmt.Sprintf("%s pay", PayOrderTableName)
		merchant   = fmt.Sprintf("left join %s m on m.merchant_no = pay.merchant_no", MerchantTableName)
		platformCh = fmt.Sprintf("left join %s plat on plat.id = pay.platform_channel_id", PlatformChannelTableName)

		selectField = "pay.id, pay.order_no, pay.merchant_order_no, pay.currency, pay.req_amount, " +
			"pay.payment_amount, pay.merchant_fee, pay.merchant_rate, pay.merchant_single_fee, pay.increase_amount, " +
			"pay.order_status, pay.create_time, pay.update_time, pay.currency, " +
			"m.username as merchant_name, plat.channel_name as platform_channel_name "

		whereBuffet = bytes.NewBufferString(" 1=1 ")
		args        []interface{}
	)

	if f.OrderNo != "" {
		whereBuffet.WriteString("and pay.order_no like ? ")
		args = append(args, "%"+f.OrderNo+"%")
	}

	if f.MerchantOrderNo != "" {
		whereBuffet.WriteString("and pay.merchant_order_no like ? ")
		args = append(args, "%"+f.MerchantOrderNo+"%")
	}

	if f.UpstreamOrderNo != "" {
		whereBuffet.WriteString("and pay.upstream_order_no like ? ")
		args = append(args, "%"+f.UpstreamOrderNo+"%")
	}

	if f.MerchantId != 0 {
		whereBuffet.WriteString("and m.id = ? ")
		args = append(args, f.MerchantId)
	}

	if f.MerchantName != "" {
		whereBuffet.WriteString("and m.username like ? ")
		args = append(args, "%"+f.MerchantName+"%")
	}

	if f.Currency != "" {
		whereBuffet.WriteString("and pay.currency = ? ")
		args = append(args, f.Currency)
	}

	if f.PlatformChannelId != 0 {
		whereBuffet.WriteString("and pay.platform_channel_id = ? ")
		args = append(args, f.PlatformChannelId)
	}

	if f.PlatformChannelName != "" {
		whereBuffet.WriteString("and plat.channel_name like ? ")
		args = append(args, "%"+f.PlatformChannelName+"%")
	}

	if f.OrderStatus != 0 {
		whereBuffet.WriteString("and pay.order_status = ? ")
		args = append(args, f.OrderStatus)
	}

	if f.StartCreateTime != 0 {
		whereBuffet.WriteString("and pay.create_time >= ? ")
		args = append(args, f.StartCreateTime)
	}

	if f.EndCreateTime != 0 {
		whereBuffet.WriteString("and pay.create_time <= ? ")
		args = append(args, f.EndCreateTime)
	}

	if f.OrderType != "" {
		whereBuffet.WriteString("and pay.mode = ? ")
		args = append(args, f.OrderType)
	}

	var total, totalReqAmount, totalPayAmount, totalFee, totalIncreaseAmount int64

	// 订单总数
	result := m.db.Table(payOrder).Joins(merchant).Joins(platformCh).Where(whereBuffet.String(), args...).Count(&total)
	if result.Error != nil {
		return nil, result.Error
	}

	if total == 0 {
		return nil, nil
	}

	// 请求金额总数
	result = m.db.Table(payOrder).Select("IFNULL(sum(pay.req_amount), 0)").Joins(merchant).Joins(platformCh).Where(whereBuffet.String(), args...).Scan(&totalReqAmount)
	if result.Error != nil {
		return nil, result.Error
	}

	// 实际支付金额总数
	result = m.db.Table(payOrder).Select("IFNULL(sum(pay.payment_amount), 0)").Joins(merchant).Joins(platformCh).Where(whereBuffet.String(), args...).Scan(&totalPayAmount)
	if result.Error != nil {
		return nil, result.Error
	}

	// 商户手续费总数
	result = m.db.Table(payOrder).Select("IFNULL(sum(pay.merchant_fee), 0)").Joins(merchant).Joins(platformCh).Where(whereBuffet.String(), args...).Scan(&totalFee)
	if result.Error != nil {
		return nil, result.Error
	}

	// 商户实际入账金额总数
	result = m.db.Table(payOrder).Select("IFNULL(sum(pay.increase_amount), 0)").Joins(merchant).Joins(platformCh).
		Where(whereBuffet.String()+fmt.Sprintf("and pay.order_status = %v ", PayOrderStatusPaid), args...).Scan(&totalIncreaseAmount)
	if result.Error != nil {
		return nil, result.Error
	}

	whereBuffet.WriteString("order by pay.create_time asc ")

	var list []PayOrderList
	result = m.db.Table(payOrder).Select(selectField).Joins(merchant).Joins(platformCh).Where(whereBuffet.String(), args...).Scan(&list)
	if result.Error != nil {
		return nil, result.Error
	}

	obj := new(ExportData)
	obj.Total = total
	obj.TotalReqAmount = totalReqAmount
	obj.TotalPayAmount = totalPayAmount
	obj.TotalMerchantFee = totalFee
	obj.TotalIncreaseAmount = totalIncreaseAmount
	obj.OrderList = list

	return obj, nil
}
