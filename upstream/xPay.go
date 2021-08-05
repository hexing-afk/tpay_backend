package upstream

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"net/url"
	"strings"
	"time"
	"tpay_backend/utils"
)

const (
	// 接口请求地址
	XPayPayUrl                = "/pay"                  // 代收下单
	XPayPayOrderQueryUrl      = "/pay-order-query"      // 代收订单查询
	XPayTransferUrl           = "/transfer"             // 代付下单
	XPayTransferOrderQueryUrl = "/transfer-order-query" // 代付订单查询
	XPayQueryBalance          = "/query-balance"        // 查询余额

	// 签名字段
	XPaySignFiledName = "sign"

	// 订单状态
	XPayOrderStatusPaying     = 1 // 正在支付中
	XPayOrderStatusPaySuccess = 2 // 支付成功
	XPayOrderStatusPayFail    = 3 // 支付失败

)

type XPay struct {
	upMerchantNo string // 上游账号id
	config       XPayConfig

	payUrl                string
	payOrderQueryUrl      string
	transferUrl           string
	transferOrderQueryUrl string
	queryBalanceUrl       string
}

type XPayConfig struct {
	Host               string `json:"Host"`               // 请求的地址
	SecretKey          string `json:"SecretKey"`          // 商家通信秘钥
	PayNotifyPath      string `json:"PayNotifyPath"`      // 代收异步通知path路径
	TransferNotifyPath string `json:"TransferNotifyPath"` // 代付异步通知path路径
}

type XPayResponse struct {
	Code int64  `json:"code"` // 请求返回代码，值见数据字典
	Msg  string `json:"msg"`  // 出错消息，请求处理失败才会出现
	Data string `json:"data"` // 业务数据
	Sign string `json:"sign"` // 签名
}

func CheckXPayConfig(conf XPayConfig) error {
	if strings.TrimSpace(conf.Host) == "" {
		return errors.New("XPay.Host配置不能为空")
	}

	if strings.TrimSpace(conf.SecretKey) == "" {
		return errors.New("XPay.SecretKey配置不能为空")
	}

	if strings.TrimSpace(conf.PayNotifyPath) == "" {
		return errors.New("XPay.PayNotifyPath配置不能为空")
	}

	if strings.TrimSpace(conf.TransferNotifyPath) == "" {
		return errors.New("XPay.TransferNotifyPath配置不能为空")
	}

	// 更多判断...

	return nil
}

func NewXPay(upMerchantNo string, jsonStrConfig string) (Upstream, error) {
	c := XPayConfig{}

	if strings.TrimSpace(upMerchantNo) == "" {
		return nil, errors.New("upMerchantNo配置不能为空")
	}

	// 解析配置
	if err := json.Unmarshal([]byte(jsonStrConfig), &c); err != nil {
		return nil, err
	}

	// 检查配置
	if err := CheckXPayConfig(c); err != nil {
		return nil, err
	}

	o := &XPay{}
	o.config = c
	o.upMerchantNo = upMerchantNo

	o.payUrl = strings.TrimRight(c.Host, "/") + XPayPayUrl
	o.payOrderQueryUrl = strings.TrimRight(c.Host, "/") + XPayPayOrderQueryUrl
	o.transferUrl = strings.TrimRight(c.Host, "/") + XPayTransferUrl
	o.transferOrderQueryUrl = strings.TrimRight(c.Host, "/") + XPayTransferOrderQueryUrl
	o.queryBalanceUrl = strings.TrimRight(c.Host, "/") + XPayQueryBalance

	return o, nil
}

// XPay上游的公共请求参数
func (o *XPay) RequestCommonField() map[string]interface{} {
	reqData := make(map[string]interface{})
	reqData["merchant_no"] = o.upMerchantNo  // String	M  商户号，平台分配账号
	reqData["timestamp"] = time.Now().Unix() // int64   M  请求时间

	return reqData
}

// 获取上游配置
func (o *XPay) GetUpstreamConfig() *ConfigResponse {
	return &ConfigResponse{
		PayNotifyPath:      o.config.PayNotifyPath,
		TransferNotifyPath: o.config.TransferNotifyPath,
		SecretKey:          o.config.SecretKey,
	}
}

func (o *XPay) Pay(req *PayRequest) (*PayResponse, error) {
	// 1.拼接参数
	reqParam := o.RequestCommonField()
	reqParam["subject"] = req.Subject        // String   O
	reqParam["amount"] = req.Amount          // String	M  金额
	reqParam["currency"] = req.Currency      // String	M  币种
	reqParam["mch_order_no"] = req.OrderNo   // String	M  商户订单号（唯一）
	reqParam["trade_type"] = req.ProductType // String	M  交易类型
	reqParam["notify_url"] = req.NotifyUrl   // String	M  异步通知回调地址
	reqParam["return_url"] = req.ReturnUrl   // String	O  页面回跳地址
	reqParam["attach"] = req.Attach          // String	M  备注

	signData, err := json.Marshal(reqParam)
	if err != nil {
		return nil, err
	}

	signDataStr := string(signData)

	params := make(map[string]interface{})
	params["data"] = signDataStr
	sign := o.GenerateSign(params)

	urlValue := url.Values{}
	urlValue.Set("data", signDataStr)
	urlValue.Set("sign", sign)

	funcName := utils.RunFuncName()

	// 2.发送请求
	logx.Infof(funcName+":request:%v", urlValue)
	body, err := utils.PostForm(o.payUrl, urlValue)
	logx.Infof(funcName+":response:%v", string(body))

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("response body is empty")
	}

	// 3.解析返回结果
	var resp XPayResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	// 4.验证业务是否成功
	if resp.Code != 0 {
		return nil, errors.New(fmt.Sprintf("response code failed,code:%v,msg:%v", resp.Code, resp.Msg))
	}

	// 5.验签
	checkSignMap := make(map[string]interface{})
	checkSignMap["data"] = resp.Data
	checkSignMap[XPaySignFiledName] = resp.Sign
	if err := o.CheckSign(checkSignMap); err != nil {
		return nil, errors.New(fmt.Sprintf("sign verification failed, err:%v, body:%v", err, string(body)))
	}

	// 6.返回结果
	var data struct {
		MchOrderNo string `json:"mch_order_no"` // 外部订单号(商户系统内部的订单号)
		OrderNo    string `json:"order_no"`     // 平台订单号
		PayUrl     string `json:"pay_url"`      // 付款收银台地址
	}
	if err := json.Unmarshal([]byte(resp.Data), &data); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json data failed, err:%v, data:%v", err, resp.Data))
	}

	return &PayResponse{
		UpstreamOrderNo: data.OrderNo,
		PayUrl:          data.PayUrl,
	}, nil
}

func (o *XPay) PayOrderQuery(req *PayOrderQueryRequest) (*PayOrderQueryResponse, error) {
	// 1.拼接参数
	reqParam := o.RequestCommonField()
	reqParam["mch_order_no"] = req.OrderNo
	reqParam["order_no"] = req.UpstreamOrderNo

	dataByte, jErr := json.Marshal(reqParam)
	if jErr != nil {
		return nil, jErr
	}

	dataByteStr := string(dataByte)

	params := make(map[string]interface{})
	params["data"] = dataByteStr
	sign := o.GenerateSign(params)

	urlValue := url.Values{}
	urlValue.Set("data", dataByteStr)
	urlValue.Set("sign", sign)

	funcName := utils.RunFuncName()

	// 2.发送请求
	logx.Infof(funcName+":request:%v", urlValue)
	body, err := utils.PostForm(o.payOrderQueryUrl, urlValue)
	logx.Infof(funcName+":response:%v", string(body))

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("response body is empty")
	}

	var resp XPayResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	if resp.Code != 0 {
		return nil, errors.New(fmt.Sprintf("response code failed,code:%v,msg:%v", resp.Code, resp.Msg))
	}

	// 5.验签
	checkSignMap := make(map[string]interface{})
	checkSignMap[XPaySignFiledName] = resp.Sign
	checkSignMap["data"] = resp.Data
	if err := o.CheckSign(checkSignMap); err != nil {
		return nil, errors.New(fmt.Sprintf("sign verification failed, err:%v, body:%v", err, string(body)))
	}

	// 6.返回结果
	var data struct {
		MchOrderNo string `json:"mch_order_no"` // 外部订单号(商户系统内部的订单号)
		OrderNo    string `json:"order_no"`     // 平台订单号
		Amount     int64  `json:"amount"`       // 订单金额
		Status     int64  `json:"status"`       // 订单状态
		Currency   string `json:"currency"`     // 币种
	}
	if err := json.Unmarshal([]byte(resp.Data), &data); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json data failed, err:%v, data:%v", err, resp.Data))
	}

	var orderStatus int64
	switch data.Status {
	case XPayOrderStatusPaying:
		orderStatus = PayPaying // 支付中
	case XPayOrderStatusPaySuccess:
		orderStatus = PaySuccess // 支付成功
	case XPayOrderStatusPayFail:
		orderStatus = PayFail // 支付失败
	default:
		return nil, errors.New(fmt.Sprintf("response unknown order status, data:%v", resp.Data))
	}

	return &PayOrderQueryResponse{
		OrderNo:         data.MchOrderNo,
		UpstreamOrderNo: data.OrderNo,
		OrderStatus:     orderStatus,
	}, nil
}

func (o *XPay) Transfer(req *TransferRequest) (*TransferResponse, error) {
	// 1.拼接参数
	reqParam := o.RequestCommonField()
	reqParam["amount"] = req.Amount                            // String M	订单金额
	reqParam["currency"] = req.Currency                        // String M	币种
	reqParam["mch_order_no"] = req.OrderNo                     // String M	商户系统内部的订单号
	reqParam["trade_type"] = req.ProductType                   // String M	交易类型
	reqParam["notify_url"] = req.NotifyUrl                     // String M	异步通知地址
	reqParam["return_url"] = req.ReturnUrl                     // String O	同步跳转地址
	reqParam["attach"] = req.Attach                            // String O	原样返回字段
	reqParam["bank_name"] = req.BankName                       // String M	收款银行名称
	reqParam["bank_card_holder_name"] = req.BankCardHolderName // String M	银行卡持卡人姓名
	reqParam["bank_card_no"] = req.BankCardNo                  // String M	银行卡号
	reqParam["bank_branch_name"] = req.BankBranchName          // String O	收款银行支行名称
	reqParam["bank_code"] = req.BankCode                       // String O	银行代码
	reqParam["remark"] = req.Remark                            // String O	付款备注

	dataByte, jErr := json.Marshal(reqParam)
	if jErr != nil {
		logx.Errorf("map转json失败, err=%v", jErr)
		return nil, jErr
	}

	dataByteStr := string(dataByte)

	// 2.生成签名
	params := make(map[string]interface{})
	params["data"] = dataByteStr
	sign := o.GenerateSign(params)

	urlValue := url.Values{}
	urlValue.Set("data", dataByteStr)
	urlValue.Set("sign", sign)

	funcName := utils.RunFuncName()

	// 3.发送请求
	logx.Infof(funcName+":request:%v", urlValue)
	body, err := utils.PostForm(o.transferUrl, urlValue)
	logx.Infof(funcName+":response:%v", string(body))

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("response body is empty")
	}

	// 4.解析返回结果
	var resp XPayResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	if resp.Code != 0 {
		return nil, errors.New(fmt.Sprintf("response code failed,code:%v,msg:%v", resp.Code, resp.Msg))
	}

	// 5.验签
	checkSignMap := make(map[string]interface{})
	checkSignMap["data"] = resp.Data
	checkSignMap[XPaySignFiledName] = resp.Sign
	if err := o.CheckSign(checkSignMap); err != nil {
		return nil, errors.New(fmt.Sprintf("sign verification failed, err:%v, body:%v", err, string(body)))
	}

	// 6.处理结果
	var data struct {
		MchOrderNo string `json:"mch_order_no"` // 外部订单号(商户系统内部的订单号)
		OrderNo    string `json:"order_no"`     // 平台订单号
		Status     int64  `json:"status"`       // 订单状态
	}

	if err := json.Unmarshal([]byte(resp.Data), &data); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json data failed, err:%v, data:%v", err, resp.Data))
	}

	var orderStatus int64
	switch data.Status {
	case XPayOrderStatusPaying:
		orderStatus = PayPaying // 支付中
	case XPayOrderStatusPayFail:
		orderStatus = PayFail // 支付失败
	case XPayOrderStatusPaySuccess:
		orderStatus = PaySuccess // 支付成功
	default:
		return nil, errors.New(fmt.Sprintf("response unknown order status, data:%v", resp.Data))
	}

	return &TransferResponse{
		UpstreamOrderNo: data.OrderNo,
		OrderStatus:     orderStatus,
	}, nil

}

func (o *XPay) TransferOrderQuery(req *TransferOrderQueryRequest) (*TransferOrderQueryResponse, error) {
	// 1.参数拼接
	reqData := o.RequestCommonField()
	reqData["order_no"] = req.UpstreamOrderNo // String O	上游订单号
	reqData["mch_order_no"] = req.OrderNo     // String O	商户订单号

	dataByte, jErr := json.Marshal(reqData)
	if jErr != nil {
		return nil, jErr
	}

	dataByteStr := string(dataByte)

	// 2.生成签名
	params := make(map[string]interface{})
	params["data"] = dataByteStr
	sign := o.GenerateSign(params)

	urlValue := url.Values{}
	urlValue.Set("data", dataByteStr)
	urlValue.Set("sign", sign)

	funcName := utils.RunFuncName()

	// 3.发起请求
	logx.Infof(funcName+":request:%v", urlValue)
	body, err := utils.PostForm(o.transferOrderQueryUrl, urlValue)
	logx.Infof(funcName+":response:%v", string(body))
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("response body is empty")
	}

	// 4.解析返回结果
	var resp XPayResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	// 5.验证业务是否成功
	if resp.Code != 0 {
		return nil, errors.New(fmt.Sprintf("response code failed,code:%v,msg:%v", resp.Code, resp.Msg))
	}

	// 6.验签
	checkSignMap := make(map[string]interface{})
	checkSignMap["data"] = resp.Data
	checkSignMap[XPaySignFiledName] = resp.Sign
	if err := o.CheckSign(checkSignMap); err != nil {
		return nil, errors.New(fmt.Sprintf("sign verification failed, err:%v, body:%v", err, string(body)))
	}

	// 7.处理结果
	var data struct {
		MchOrderNo string `json:"mch_order_no"` // 外部订单号(商户系统内部的订单号)
		OrderNo    string `json:"order_no"`     // 平台订单号
		Amount     int64  `json:"amount"`       // 订单金额
		Status     int64  `json:"status"`       // 订单状态
		Currency   string `json:"currency"`     // 币种
	}
	if err := json.Unmarshal([]byte(resp.Data), &data); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json data failed, err:%v, data:%v", err, resp.Data))
	}

	var orderStatus int64
	switch data.Status {
	case XPayOrderStatusPaying:
		orderStatus = PayPaying // 支付中
	case XPayOrderStatusPaySuccess:
		orderStatus = PaySuccess // 支付成功
	case XPayOrderStatusPayFail:
		orderStatus = PayFail // 支付失败
	default:
		return nil, errors.New(fmt.Sprintf("response unknown order status, data:%v", resp.Data))
	}

	return &TransferOrderQueryResponse{
		OrderNo:         data.MchOrderNo,
		UpstreamOrderNo: data.OrderNo,
		OrderStatus:     orderStatus,
	}, nil
}

func (o *XPay) QueryBalance() (*QueryBalanceResponse, error) {
	// 1.参数拼接
	reqData := o.RequestCommonField()
	dataByte, jErr := json.Marshal(reqData)
	if jErr != nil {
		return nil, jErr
	}

	// 2.生成签名
	params := make(map[string]interface{})
	params["data"] = string(dataByte)
	sign := o.GenerateSign(params)

	urlValue := url.Values{}
	urlValue.Set("data", string(dataByte))
	urlValue.Set("sign", sign)

	funcName := utils.RunFuncName()

	// 3.发送请求
	logx.Infof(funcName+":request:%v", urlValue)
	body, err := utils.PostForm(o.queryBalanceUrl, urlValue)
	logx.Infof(funcName+":response:%v", string(body))

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("response body is empty")
	}

	var resp XPayResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	if resp.Code != 0 {
		return nil, errors.New(fmt.Sprintf("response code failed,code:%v,msg:%v", resp.Code, resp.Msg))
	}

	// 5.处理结果
	var data struct {
		Balance  int64  `json:"balance"`
		Currency string `json:"currency"`
	}

	if err := json.Unmarshal([]byte(resp.Data), &data); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json data failed, err:%v, data:%v", err, resp.Data))
	}

	return &QueryBalanceResponse{
		Balance:  float64(data.Balance),
		Currency: data.Currency,
	}, nil
}

// 生成签名
func (o *XPay) GenerateSign(data map[string]interface{}) string {
	signData := data["data"].(string)
	return strings.ToLower(utils.Md5(signData + o.config.SecretKey))
}

// 校验签名
func (o *XPay) CheckSign(data map[string]interface{}) error {
	sign, exist := data[XPaySignFiledName]
	if !exist {
		return errors.New("no sign field")
	}

	if sign == "" {
		return errors.New("sign field empty")
	}

	if sign != o.GenerateSign(data) {
		return errors.New("sign not match")
	}

	return nil
}
