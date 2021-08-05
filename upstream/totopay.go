package upstream

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"tpay_backend/utils"

	"github.com/tal-tech/go-zero/core/logx"
)

const (
	TotopayPayUrl                = "/api/pay/prepay"
	TotopayPayOrderQueryUrl      = "/api/pay/query"
	TotopayTransferUrl           = "/api/pay/transfer"
	TotopayTransferOrderQueryUrl = "/api/pay/transfer_query"
	TotopayQueryBalance          = "/api/pay/query_balance"

	TotopayOrderStatusPending = "1" // 待支付
	TotopayOrderStatusSuccess = "2" // 支付成功
	TotopayOrderStatusFail    = "3" // 支付失败

	TotopaySignFiledName = "sign"
)

type Totopay struct {
	upMerchantNo string // 上游账号id
	config       TotopayConfig

	payUrl                string
	payOrderQueryUrl      string
	transferUrl           string
	transferOrderQueryUrl string
	queryBalanceUrl       string
}

type TotopayConfig struct {
	Host               string `json:"Host"`               // 请求的地址
	SecretKey          string `json:"SecretKey"`          // 商家通信秘钥
	PayNotifyPath      string `json:"PayNotifyPath"`      // 代收异步通知path路径
	TransferNotifyPath string `json:"TransferNotifyPath"` // 代付异步通知path路径
}

func CheckTotopayConfig(conf TotopayConfig) error {
	if strings.TrimSpace(conf.Host) == "" {
		return errors.New("Totopay.Host配置不能为空")
	}

	if strings.TrimSpace(conf.SecretKey) == "" {
		return errors.New("Totopay.SecretKey配置不能为空")
	}

	if strings.TrimSpace(conf.PayNotifyPath) == "" {
		return errors.New("Totopay.PayNotifyPath配置不能为空")
	}

	if strings.TrimSpace(conf.TransferNotifyPath) == "" {
		return errors.New("Totopay.TransferNotifyPath配置不能为空")
	}

	// 更多判断...

	return nil
}

func NewTotopay(upMerchantNo string, jsonStrConfig string) (Upstream, error) {
	c := TotopayConfig{}

	if strings.TrimSpace(upMerchantNo) == "" {
		return nil, errors.New("upMerchantNo配置不能为空")
	}

	// 解析配置
	if err := json.Unmarshal([]byte(jsonStrConfig), &c); err != nil {
		return nil, err
	}

	// 检查配置
	if err := CheckTotopayConfig(c); err != nil {
		return nil, err
	}

	o := &Totopay{}
	o.config = c
	o.upMerchantNo = upMerchantNo

	o.payUrl = strings.TrimRight(c.Host, "/") + TotopayPayUrl
	o.payOrderQueryUrl = strings.TrimRight(c.Host, "/") + TotopayPayOrderQueryUrl
	o.transferUrl = strings.TrimRight(c.Host, "/") + TotopayTransferUrl
	o.transferOrderQueryUrl = strings.TrimRight(c.Host, "/") + TotopayTransferOrderQueryUrl
	o.queryBalanceUrl = strings.TrimRight(c.Host, "/") + TotopayQueryBalance

	return o, nil
}

// 获取上游配置
func (o *Totopay) GetUpstreamConfig() *ConfigResponse {
	return &ConfigResponse{
		PayNotifyPath:      o.config.PayNotifyPath,
		TransferNotifyPath: o.config.TransferNotifyPath,
		SecretKey:          o.config.SecretKey,
	}
}

func (o *Totopay) Pay(req *PayRequest) (*PayResponse, error) {
	params := make(map[string]interface{})
	params["amount"] = req.Amount
	params["req_no"] = req.OrderNo
	params["notify_url"] = req.NotifyUrl
	params["frontend_url"] = req.ReturnUrl
	params["product_type"] = req.ProductType
	params["acc_no"] = o.upMerchantNo
	params["sign"] = o.GenerateSign(params)

	dataByte, jerr := json.Marshal(params)
	if jerr != nil {
		return nil, jerr
	}

	funcName := utils.RunFuncName()

	logx.Infof(funcName+":request:%v", string(dataByte))
	body, err := utils.PostJson(o.payUrl, dataByte)
	logx.Infof(funcName+":response:%v", string(body))

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("response body is empty")
	}

	var res struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			CodeUrl string `json:"code_url"`
			OrderNo string `json:"order_no"`
		}
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	if res.Code != "0" {
		return nil, errors.New(fmt.Sprintf("response code failed,code:%v,msg:%v", res.Code, res.Msg))
	}

	return &PayResponse{
		UpstreamOrderNo: res.Data.OrderNo,
		PayUrl:          res.Data.CodeUrl,
	}, nil
}

func (o *Totopay) PayOrderQuery(req *PayOrderQueryRequest) (*PayOrderQueryResponse, error) {
	params := make(map[string]interface{})
	params["req_no"] = req.OrderNo
	params["acc_no"] = o.upMerchantNo
	params["sign"] = o.GenerateSign(params)

	dataByte, jerr := json.Marshal(params)
	if jerr != nil {
		return nil, jerr
	}

	funcName := utils.RunFuncName()

	logx.Infof(funcName+":request:%v", string(dataByte))
	body, err := utils.PostJson(o.payOrderQueryUrl, dataByte)
	logx.Infof(funcName+":response:%v", string(body))

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("response body is empty")
	}

	var res struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Amount      string `json:"amount"`
			OrderNo     string `json:"order_no"` // 上游单号
			OrderStatus string `json:"order_status"`
			ProductType string `json:"product_type"`
			ReqAmount   string `json:"req_amount"`
			ReqNo       string `json:"req_no"`
		}
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	if res.Code != "0" {
		return nil, errors.New(fmt.Sprintf("response code failed,code:%v,msg:%v", res.Code, res.Msg))
	}

	var orderStatus int64
	if res.Data.OrderStatus == TotopayOrderStatusSuccess {
		orderStatus = PaySuccess
	} else if res.Data.OrderStatus == TotopayOrderStatusFail {
		orderStatus = PayFail
	}

	return &PayOrderQueryResponse{
		OrderNo:         res.Data.ReqNo,
		UpstreamOrderNo: res.Data.OrderNo,
		OrderStatus:     orderStatus,
	}, nil
}

func (o *Totopay) Transfer(req *TransferRequest) (*TransferResponse, error) {
	params := make(map[string]interface{})
	params["amount"] = req.Amount
	params["req_no"] = req.OrderNo
	params["notify_url"] = req.NotifyUrl
	params["transfer_type"] = "withdraw"
	params["product_type"] = req.ProductType
	params["bankname"] = req.BankName
	params["card_no"] = req.BankCardNo
	params["card_acc"] = req.BankCardHolderName
	params["acc_no"] = o.upMerchantNo
	params["sign"] = o.GenerateSign(params)

	dataByte, jErr := json.Marshal(params)
	if jErr != nil {
		logx.Errorf("map转json失败, err=%v", jErr)
		return nil, jErr
	}

	funcName := utils.RunFuncName()

	logx.Infof(funcName+":request:%v", string(dataByte))
	body, err := utils.PostJson(o.transferUrl, dataByte)
	logx.Infof(funcName+":response:%v", string(body))

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("response body is empty")
	}

	var res struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			OrderStatus string `json:"order_status"`
			OrderNo     string `json:"order_no"`
		}
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	if res.Code != "0" {
		return nil, errors.New(fmt.Sprintf("response code failed,code:%v,msg:%v", res.Code, res.Msg))
	}

	var orderStatus int64
	if res.Data.OrderStatus == TotopayOrderStatusSuccess {
		orderStatus = TransferSuccess
	} else if res.Data.OrderStatus == TotopayOrderStatusFail {
		orderStatus = TransferFail
	}

	return &TransferResponse{
		UpstreamOrderNo: res.Data.OrderNo,
		OrderStatus:     orderStatus,
	}, nil

}

func (o *Totopay) TransferOrderQuery(req *TransferOrderQueryRequest) (*TransferOrderQueryResponse, error) {
	params := make(map[string]interface{})
	params["req_no"] = req.OrderNo
	params["acc_no"] = o.upMerchantNo
	params["sign"] = o.GenerateSign(params)

	dataByte, jerr := json.Marshal(params)
	if jerr != nil {
		return nil, jerr
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

	var res struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Amount      string `json:"amount"`
			OrderNo     string `json:"order_no"` // 上游单号
			OrderStatus string `json:"order_status"`
		}
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	if res.Code != "0" {
		return nil, errors.New(fmt.Sprintf("response code failed,code:%v,msg:%v", res.Code, res.Msg))
	}

	var orderStatus int64
	if res.Data.OrderStatus == TotopayOrderStatusSuccess {
		orderStatus = TransferSuccess
	} else if res.Data.OrderStatus == TotopayOrderStatusFail {
		orderStatus = TransferFail
	}

	return &TransferOrderQueryResponse{
		OrderNo:         res.Data.OrderNo,
		UpstreamOrderNo: res.Data.OrderNo,
		OrderStatus:     orderStatus,
	}, nil
}

func (o *Totopay) QueryBalance() (*QueryBalanceResponse, error) {
	params := make(map[string]interface{})
	params["acc_no"] = o.upMerchantNo
	params["sign"] = o.GenerateSign(params)

	dataByte, jerr := json.Marshal(params)
	if jerr != nil {
		return nil, jerr
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

	var res struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			BalanceAmount string `json:"balance_amount"`
		}
	}

	if err := json.Unmarshal(body, &res); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	if res.Code != "0" {
		return nil, errors.New(fmt.Sprintf("response code failed,code:%v,msg:%v", res.Code, res.Msg))
	}

	balance, err := strconv.ParseFloat(res.Data.BalanceAmount, 64)
	//balance, err := strconv.ParseInt(res.Data.BalanceAmount, 10, 64)
	if err != nil {
		return nil, err
	}

	return &QueryBalanceResponse{
		Balance: balance,
	}, nil
}

// 生成签名
func (o *Totopay) GenerateSign(data map[string]interface{}) string {
	// 将参数排序并拼接成字符串
	dataStr := utils.ParamsMapToString(data, TotopaySignFiledName)
	dataStr = fmt.Sprintf("%s&key=%s", dataStr, o.config.SecretKey)
	return strings.ToLower(utils.Md5(dataStr))
}

// 校验签名
func (o *Totopay) CheckSign(data map[string]interface{}) error {
	sign, exist := data[TotopaySignFiledName]
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
