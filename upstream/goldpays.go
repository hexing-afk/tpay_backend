package upstream

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/tal-tech/go-zero/core/logx"
	"strconv"
	"strings"
	"tpay_backend/utils"
)

const (
	// 接口请求地址
	GoldPaysPayUrl                = "/openApi/pay/createOrder"    // 代收下单
	GoldPaysPayOrderQueryUrl      = "/openApi/pay/queryOrder"     // 代收订单查询
	GoldPaysTransferUrl           = "/openApi/payout/createOrder" // 代付下单
	GoldPaysTransferOrderQueryUrl = "/openApi/payout/queryOrder"  // 代付订单查询
	GoldPaysQueryBalance          = "/openApi/balances"           // 查询余额

	// 签名字段
	GoldPaysSignFiledName = "sign"

	// 币种
	GoldPaysCurrency = "INR" // 卢比，上游金额单位是卢布; 1卢比=100分

	// 订单状态
	GoldPaysOrderStatusCreate     = "PAY_CREATE"  // 创建订单
	GoldPaysOrderStatusPaying     = "PAY_ING"     // 正在支付中
	GoldPaysOrderStatusPayFail    = "PAY_FAIL"    // 支付失败
	GoldPaysOrderStatusPaySuccess = "PAY_SUCCESS" // 支付成功

	// 代付模式：mode
	GoldPaysPayOutModeIMPS = "IMPS" // 使用IFSC和银行卡号代付
	GoldPaysPayOutModeUPI  = "UPI"  // 使用UPI账号代付(通道不一定支持，需先咨询商务)

)

type GoldPays struct {
	upMerchantNo string // 上游账号id
	config       GoldPaysConfig

	payUrl                string
	payOrderQueryUrl      string
	transferUrl           string
	transferOrderQueryUrl string
	queryBalanceUrl       string
}

type GoldPaysConfig struct {
	Host               string `json:"Host"`               // 请求的地址
	SecretKey          string `json:"SecretKey"`          // 商家通信秘钥
	PayNotifyPath      string `json:"PayNotifyPath"`      // 代收异步通知path路径
	TransferNotifyPath string `json:"TransferNotifyPath"` // 代付异步通知path路径
}

func CheckGoldPaysConfig(conf GoldPaysConfig) error {
	if strings.TrimSpace(conf.Host) == "" {
		return errors.New("GoldPays.Host配置不能为空")
	}

	if strings.TrimSpace(conf.SecretKey) == "" {
		return errors.New("GoldPays.SecretKey配置不能为空")
	}

	if strings.TrimSpace(conf.PayNotifyPath) == "" {
		return errors.New("GoldPays.PayNotifyPath配置不能为空")
	}

	if strings.TrimSpace(conf.TransferNotifyPath) == "" {
		return errors.New("GoldPays.TransferNotifyPath配置不能为空")
	}

	// 更多判断...

	return nil
}

func NewGoldPays(upMerchantNo string, jsonStrConfig string) (Upstream, error) {
	c := GoldPaysConfig{}

	if strings.TrimSpace(upMerchantNo) == "" {
		return nil, errors.New("upMerchantNo配置不能为空")
	}

	// 解析配置
	if err := json.Unmarshal([]byte(jsonStrConfig), &c); err != nil {
		return nil, err
	}

	// 检查配置
	if err := CheckGoldPaysConfig(c); err != nil {
		return nil, err
	}

	o := &GoldPays{}
	o.config = c
	o.upMerchantNo = upMerchantNo

	o.payUrl = strings.TrimRight(c.Host, "/") + GoldPaysPayUrl
	o.payOrderQueryUrl = strings.TrimRight(c.Host, "/") + GoldPaysPayOrderQueryUrl
	o.transferUrl = strings.TrimRight(c.Host, "/") + GoldPaysTransferUrl
	o.transferOrderQueryUrl = strings.TrimRight(c.Host, "/") + GoldPaysTransferOrderQueryUrl
	o.queryBalanceUrl = strings.TrimRight(c.Host, "/") + GoldPaysQueryBalance

	return o, nil
}

// 获取上游配置
func (o *GoldPays) GetUpstreamConfig() *ConfigResponse {
	return &ConfigResponse{
		PayNotifyPath:      o.config.PayNotifyPath,
		TransferNotifyPath: o.config.TransferNotifyPath,
		SecretKey:          o.config.SecretKey,
	}
}

func (o *GoldPays) Pay(req *PayRequest) (*PayResponse, error) {
	// 0.参数处理
	// 金额除以100，再四舍五入保留两位小数点
	amount := decimal.NewFromInt(req.Amount).Div(decimal.NewFromInt(100)).Round(2).String()

	if req.CustomName == "" {
		req.CustomName = "李明"
		req.CustomMobile = "138690108"
		req.CustomEmail = "h132@193.com"
	}

	// 1.拼接参数
	params := make(map[string]interface{})
	params["merchant"] = o.upMerchantNo       // String	M 商户号，平台分配账号
	params["orderId"] = req.OrderNo           // String	M 商户订单号（唯一），字符长度40以内
	params["amount"] = amount                 // String	M  金额，单位卢币(最多保留两位小数)
	params["customName"] = req.CustomName     // String	M  客户姓名
	params["customMobile"] = req.CustomMobile // String	M  客户电话
	params["customEmail"] = req.CustomEmail   // String	M  客户email地址
	params["notifyUrl"] = req.NotifyUrl       // String	M  异步通知回调地址
	params["callbackUrl"] = req.ReturnUrl     // String	M  页面回跳地址（客户操作 支付成功或失败后跳转页面。）
	params["sign"] = o.GenerateSign(params)   // String	M  签名

	dataByte, jErr := json.Marshal(params)
	if jErr != nil {
		return nil, jErr
	}

	funcName := utils.RunFuncName()

	// 2.发送请求
	logx.Infof(funcName+":request:%v", string(dataByte))
	body, err := utils.PostJson(o.payUrl, dataByte)
	logx.Infof(funcName+":response:%v", string(body))

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("response body is empty")
	}

	// 3.解析返回结果
	var resp struct {
		Code          int64  `json:"code"`          // 请求返回代码，值见数据字典
		Success       bool   `json:"success"`       // 请求是否成功
		ErrorMessages string `json:"errorMessages"` // 出错消息，请求处理失败才会出现
		Data          struct {
			Merchant    string `json:"merchant"`    // 商户号
			OrderId     string `json:"orderId"`     // 商户订单号
			PlatOrderId string `json:"platOrderId"` // 平台订单号
			Url         string `json:"url"`         // 收银台地址
			Sign        string `json:"sign"`        // 签名
		}
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	// 4.验证业务是否成功
	if !resp.Success {
		return nil, errors.New(fmt.Sprintf("response code failed,code:%v,msg:%v", resp.Code, resp.ErrorMessages))
	}

	// 5.验签
	checkSignMap := make(map[string]interface{})
	checkSignMap["merchant"] = resp.Data.Merchant
	checkSignMap["orderId"] = resp.Data.OrderId
	checkSignMap["platOrderId"] = resp.Data.PlatOrderId
	checkSignMap["url"] = resp.Data.Url
	checkSignMap["sign"] = resp.Data.Sign

	if err := o.CheckSign(checkSignMap); err != nil {
		return nil, errors.New(fmt.Sprintf("sign verification failed, err:%v, body:%v", err, string(body)))
	}

	// 5.返回结果
	return &PayResponse{
		UpstreamOrderNo: resp.Data.PlatOrderId,
		PayUrl:          resp.Data.Url,
	}, nil
}

func (o *GoldPays) PayOrderQuery(req *PayOrderQueryRequest) (*PayOrderQueryResponse, error) {
	// 1.拼接参数
	params := make(map[string]interface{})
	params["merchant"] = o.upMerchantNo     // String	M	商户号，平台分配账号
	params["orderId"] = req.OrderNo         // String	M	商户订单号
	params["sign"] = o.GenerateSign(params) // String	M	签名

	dataByte, jErr := json.Marshal(params)
	if jErr != nil {
		return nil, jErr
	}

	funcName := utils.RunFuncName()

	// 2.发送请求
	logx.Infof(funcName+":request:%v", string(dataByte))
	body, err := utils.PostJson(o.payOrderQueryUrl, dataByte)
	logx.Infof(funcName+":response:%v", string(body))

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("response body is empty")
	}

	var resp struct {
		Code          int64  `json:"code"`          // 请求返回代码，值见数据字典
		Success       bool   `json:"success"`       // 请求是否成功
		ErrorMessages string `json:"errorMessages"` // 出错消息，请求处理失败才会出现
		Data          struct {
			Merchant    string `json:"merchant"`    // 商户号
			OrderId     string `json:"orderId"`     // 商户订单号
			PlatOrderId string `json:"platOrderId"` // 平台订单号
			Amount      string `json:"amount"`      // 金额
			Msg         string `json:"msg"`         // 处理消息
			Status      string `json:"status"`      // 状态
			Sign        string `json:"sign"`        // 签名
		}
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	// 4.验证业务是否成功
	if !resp.Success {
		return nil, errors.New(fmt.Sprintf("response code failed,code:%v,msg:%v", resp.Code, resp.ErrorMessages))
	}

	// 5.验签
	checkSignMap := make(map[string]interface{})
	checkSignMap["merchant"] = resp.Data.Merchant
	checkSignMap["orderId"] = resp.Data.OrderId
	checkSignMap["platOrderId"] = resp.Data.PlatOrderId
	checkSignMap["amount"] = resp.Data.Amount
	checkSignMap["msg"] = resp.Data.Msg
	checkSignMap["status"] = resp.Data.Status
	checkSignMap["sign"] = resp.Data.Sign

	if err := o.CheckSign(checkSignMap); err != nil {
		return nil, errors.New(fmt.Sprintf("sign verification failed, err:%v, body:%v", err, string(body)))
	}

	// 6.处理结果
	var orderStatus int64
	switch resp.Data.Status {
	case GoldPaysOrderStatusCreate:
		orderStatus = PayPaying // 支付中
	case GoldPaysOrderStatusPaying:
		orderStatus = PayPaying // 支付中
	case GoldPaysOrderStatusPayFail:
		orderStatus = PayFail // 支付失败
	case GoldPaysOrderStatusPaySuccess:
		orderStatus = PaySuccess // 支付成功
	default:
		return nil, errors.New(fmt.Sprintf("response unknown order status, data:%v", resp.Data))
	}

	return &PayOrderQueryResponse{
		OrderNo:         resp.Data.OrderId,
		UpstreamOrderNo: resp.Data.PlatOrderId,
		OrderStatus:     orderStatus,
	}, nil
}

func (o *GoldPays) Transfer(req *TransferRequest) (*TransferResponse, error) {
	// 0.参数处理
	// 金额除以100，再四舍五入保留两位小数点
	amount := decimal.NewFromInt(req.Amount).Div(decimal.NewFromInt(100)).Round(2).String()

	// 1.拼接参数
	// todo 代付模式暂时固定为：IMPS
	params := make(map[string]interface{})
	params["merchant"] = o.upMerchantNo           // String	M	商户号，平台分配账号
	params["orderId"] = req.OrderNo               // String M	商户订单号（唯一），字符长度40以内
	params["amount"] = amount                     // String M	金额，单位卢币(最多保留两位小数)
	params["customName"] = req.BankCardHolderName // String M	客户姓名
	params["customMobile"] = req.CardHolderMobile // String M	客户电话
	params["customEmail"] = req.CardHolderEmail   // String M	客户email地址
	params["mode"] = req.ProductType              // String M	代付模式，值见数据字典
	params["bankAccount"] = req.BankCardNo        // String O	收款人银行账号（mode是IMPS必须）
	params["ifscCode"] = req.BankCode             // String O	收款人IFSC CODE（mode是IMPS必须）
	params["upiAccount"] = ""                     // String O	收款人UPI账户（mode是UPI必须）
	params["notifyUrl"] = req.NotifyUrl           // String M	异步通知回调地址
	params["sign"] = o.GenerateSign(params)       // String M	签名

	dataByte, jErr := json.Marshal(params)
	if jErr != nil {
		logx.Errorf("map转json失败, err=%v", jErr)
		return nil, jErr
	}

	funcName := utils.RunFuncName()

	// 2.发送请求
	logx.Infof(funcName+":request:%v", string(dataByte))
	body, err := utils.PostJson(o.transferUrl, dataByte)
	logx.Infof(funcName+":response:%v", string(body))

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("response body is empty")
	}

	// 3.解析返回结果
	var resp struct {
		Code          int64  `json:"code"`          // 请求返回代码，值见数据字典
		Success       bool   `json:"success"`       // 请求是否成功
		ErrorMessages string `json:"errorMessages"` // 出错消息，请求处理失败才会出现
		Data          struct {
			Merchant    string `json:"merchant"`    // 商户号
			OrderId     string `json:"orderId"`     // 商户订单号
			PlatOrderId string `json:"platOrderId"` // 平台订单号
			Amount      string `json:"amount"`      // 金额
			Msg         string `json:"msg"`         // 处理消息
			Status      string `json:"status"`      // 交易状态，值见数据字典
			Sign        string `json:"sign"`        // 签名
		}
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	// 4.验证业务是否成功
	if !resp.Success {
		return nil, errors.New(fmt.Sprintf("response code failed,code:%v,msg:%v", resp.Code, resp.ErrorMessages))
	}

	// 5.验签
	checkSignMap := make(map[string]interface{})
	checkSignMap["merchant"] = resp.Data.Merchant
	checkSignMap["orderId"] = resp.Data.OrderId
	checkSignMap["platOrderId"] = resp.Data.PlatOrderId
	checkSignMap["amount"] = resp.Data.Amount
	checkSignMap["msg"] = resp.Data.Msg
	checkSignMap["status"] = resp.Data.Status
	checkSignMap["sign"] = resp.Data.Sign

	if err := o.CheckSign(checkSignMap); err != nil {
		return nil, errors.New(fmt.Sprintf("sign verification failed, err:%v, body:%v", err, string(body)))
	}

	// 6.处理结果
	var orderStatus int64
	switch resp.Data.Status {
	case GoldPaysOrderStatusCreate:
		orderStatus = PayPaying // 支付中
	case GoldPaysOrderStatusPaying:
		orderStatus = PayPaying // 支付中
	case GoldPaysOrderStatusPayFail:
		orderStatus = PayFail // 支付失败
	case GoldPaysOrderStatusPaySuccess:
		orderStatus = PaySuccess // 支付成功
	default:
		return nil, errors.New(fmt.Sprintf("response unknown order status, data:%v", resp.Data))
	}

	return &TransferResponse{
		UpstreamOrderNo: resp.Data.PlatOrderId,
		OrderStatus:     orderStatus,
	}, nil

}

func (o *GoldPays) TransferOrderQuery(req *TransferOrderQueryRequest) (*TransferOrderQueryResponse, error) {
	params := make(map[string]interface{})
	params["merchant"] = o.upMerchantNo     // String M	商户号，平台分配账号
	params["orderId"] = req.OrderNo         // String M	商户订单号
	params["sign"] = o.GenerateSign(params) // String M	签名

	dataByte, jErr := json.Marshal(params)
	if jErr != nil {
		return nil, jErr
	}

	funcName := utils.RunFuncName()

	logx.Infof(funcName+":request:%v", string(dataByte))
	body, err := utils.PostJson(o.transferOrderQueryUrl, dataByte)
	logx.Infof(funcName+":response:%v", string(body))

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("response body is empty")
	}

	// 3.解析返回结果
	var resp struct {
		Code          int64  `json:"code"`          // 请求返回代码，值见数据字典
		Success       bool   `json:"success"`       // 请求是否成功
		ErrorMessages string `json:"errorMessages"` // 出错消息，请求处理失败才会出现
		Data          struct {
			Merchant    string `json:"merchant"`    // 商户号
			OrderId     string `json:"orderId"`     // 商户订单号
			PlatOrderId string `json:"platOrderId"` // 平台订单号
			Amount      string `json:"amount"`      // 金额
			Msg         string `json:"msg"`         // 处理消息
			Status      string `json:"status"`      // 交易状态，值见数据字典
			Sign        string `json:"sign"`        // 签名
		}
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	// 4.验证业务是否成功
	if !resp.Success {
		return nil, errors.New(fmt.Sprintf("response code failed,code:%v,msg:%v", resp.Code, resp.ErrorMessages))
	}

	// 5.验签
	checkSignMap := make(map[string]interface{})
	checkSignMap["merchant"] = resp.Data.Merchant
	checkSignMap["orderId"] = resp.Data.OrderId
	checkSignMap["platOrderId"] = resp.Data.PlatOrderId
	checkSignMap["amount"] = resp.Data.Amount
	checkSignMap["msg"] = resp.Data.Msg
	checkSignMap["status"] = resp.Data.Status
	checkSignMap["sign"] = resp.Data.Sign

	if err := o.CheckSign(checkSignMap); err != nil {
		return nil, errors.New(fmt.Sprintf("sign verification failed, err:%v, body:%v", err, string(body)))
	}

	// 6.处理结果
	var orderStatus int64
	switch resp.Data.Status {
	case GoldPaysOrderStatusCreate:
		orderStatus = PayPaying // 支付中
	case GoldPaysOrderStatusPaying:
		orderStatus = PayPaying // 支付中
	case GoldPaysOrderStatusPayFail:
		orderStatus = PayFail // 支付失败
	case GoldPaysOrderStatusPaySuccess:
		orderStatus = PaySuccess // 支付成功
	default:
		return nil, errors.New(fmt.Sprintf("response unknown order status, data:%v", resp.Data))
	}

	return &TransferOrderQueryResponse{
		OrderNo:         resp.Data.OrderId,
		UpstreamOrderNo: resp.Data.PlatOrderId,
		OrderStatus:     orderStatus,
	}, nil
}

func (o *GoldPays) QueryBalance() (*QueryBalanceResponse, error) {
	params := make(map[string]interface{})
	params["merchant"] = o.upMerchantNo     // String	M	商户号，平台分配账号
	params["sign"] = o.GenerateSign(params) // String	M	签名

	dataByte, jErr := json.Marshal(params)
	if jErr != nil {
		return nil, jErr
	}

	funcName := utils.RunFuncName()

	logx.Infof(funcName+":request:%v", string(dataByte))
	body, err := utils.PostJson(o.queryBalanceUrl, dataByte)
	logx.Infof(funcName+":response:%v", string(body))

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("response body is empty")
	}

	var resp struct {
		Code          int64  `json:"code"`          // 请求返回代码，值见数据字典
		Success       bool   `json:"success"`       // 请求是否成功
		ErrorMessages string `json:"errorMessages"` // 出错消息，请求处理失败才会出现
		Data          struct {
			Balance           string `json:"balance"`           // 余额
			PayoutBalance     string `json:"payoutBalance"`     // 代付余额
			PayAmountLimit    string `json:"payAmountLimit"`    // 今日剩余代收额度
			PayoutAmountLimit string `json:"payoutAmountLimit"` // 今日剩余代付额度
		}
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	// 5.处理结果
	if !resp.Success {
		return nil, errors.New(fmt.Sprintf("response code failed,code:%v,msg:%v", resp.Code, resp.ErrorMessages))
	}

	balance, err := strconv.ParseFloat(resp.Data.Balance, 64)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("'balance' string to float64 failed, err:%v, balance:%v", err, resp.Data.Balance))
	}

	payoutBalance, err := strconv.ParseFloat(resp.Data.PayoutBalance, 64)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("'payoutBalance' string to float64 failed, err:%v, payoutBalance:%v", err, resp.Data.PayoutBalance))
	}

	payAmountLimit, err := strconv.ParseFloat(resp.Data.PayAmountLimit, 64)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("'payAmountLimit' string to float64 failed, err:%v, payAmountLimit:%v", err, resp.Data.PayAmountLimit))
	}

	payoutAmountLimit, err := strconv.ParseFloat(resp.Data.PayoutAmountLimit, 64)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("'payoutAmountLimit' string to float64 failed, err:%v, payoutAmountLimit:%v", err, resp.Data.PayoutAmountLimit))
	}

	// 上游金额是卢布, 到平台后需要转成 1卢布=100戈比
	return &QueryBalanceResponse{
		Balance:           balance,
		PayOutBalance:     payoutBalance * 100,
		PayAmountLimit:    payAmountLimit * 100,
		PayoutAmountLimit: payoutAmountLimit * 100,
		Currency:          GoldPaysCurrency,
	}, nil
}

// 生成签名
func (o *GoldPays) GenerateSign(data map[string]interface{}) string {
	// 将参数排序并拼接成字符串
	dataStr := utils.ParamsMapToString(data, GoldPaysSignFiledName)
	dataStr = fmt.Sprintf("%s&key=%s", dataStr, o.config.SecretKey)
	return strings.ToLower(utils.Md5(dataStr))
}

// 校验签名
func (o *GoldPays) CheckSign(data map[string]interface{}) error {
	sign, exist := data[GoldPaysSignFiledName]
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
