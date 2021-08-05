package payapi

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"
	"tpay_backend/utils"
)

type Tpay struct {
	Host      string // 请求的地址
	AccNo     string // 商户账号id
	SecretKey string // 商家通信秘钥
}

func NewTpay(host, appId, appSecret string) *Tpay {
	return &Tpay{Host: host, AccNo: appId, SecretKey: appSecret}
}

type CommonRequest struct {
	MerchantNo string `json:"merchant_no"`
	Timestamp  int64  `json:"timestamp"`
}

//--------------------------Pay---------------------------------------------
type PayReq struct {
	CommonRequest
	Subject    string `json:"subject"`
	Amount     int64  `json:"amount"`
	Currency   string `json:"currency"`
	MchOrderNo string `json:"mch_order_no"`
	TradeType  string `json:"trade_type"`
	NotifyUrl  string `json:"notify_url"`
	ReturnUrl  string `json:"return_url"`
	Attach     string `json:"attach"`
}

// 返回示例:
// {"code":"AA300098","msg":"订单号重复"}
type PayRes struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		CodeUrl string `json:"code_url"`
	}
}

// 下单
func (t *Tpay) Pay(req PayReq) (*PayRes, error) {
	req.CommonRequest.MerchantNo = t.AccNo
	req.CommonRequest.Timestamp = time.Now().Unix()

	dataByte, jerr := json.Marshal(req)
	if jerr != nil {
		return nil, jerr
	}

	data := url.Values{}
	data.Set("data", string(dataByte))
	data.Set("sign", utils.Md5(string(dataByte)+t.SecretKey))

	// post请求
	body, resErr := utils.PostForm(strings.TrimRight(t.Host, "/")+"/pay", data)
	if resErr != nil {
		return nil, resErr
	}

	fmt.Println("body:", string(body))

	return nil, nil
}

//--------------------------PayOrderQuery-----------------------------------
type PayOrderQueryReq struct {
	CommonRequest
	MchOrderNo string `json:"mch_order_no"`
	OrderNo    string `json:"order_no"`
}

// 返回示例:
type PayOrderQueryRes struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		CodeUrl string `json:"code_url"`
		OrderNo string `json:"order_no"`
	}
}

func (t *Tpay) PayOrderQuery(req PayOrderQueryReq) (*PayOrderQueryRes, error) {
	req.CommonRequest.MerchantNo = t.AccNo
	req.CommonRequest.Timestamp = time.Now().Unix()

	dataByte, jerr := json.Marshal(req)
	if jerr != nil {
		return nil, jerr
	}

	data := url.Values{}
	data.Set("data", string(dataByte))
	data.Set("sign", utils.Md5(string(dataByte)+t.SecretKey))

	// post请求
	body, resErr := utils.PostForm(strings.TrimRight(t.Host, "/")+"/pay-order-query", data)
	if resErr != nil {
		return nil, resErr
	}

	fmt.Println("body:", string(body))

	return nil, nil
}

//--------------------------Transfer----------------------------------------
type TransferReq struct {
	CommonRequest
	Amount             int64  `json:"amount"`                    // 订单金额
	Currency           string `json:"currency"`                  // 币种
	MchOrderNo         string `json:"mch_order_no"`              // 外部订单号(商户系统内部的订单号)
	TradeType          string `json:"trade_type"`                // 交易类型
	NotifyUrl          string `json:"notify_url"`                // 异步通知地址
	ReturnUrl          string `json:"return_url"`                // 同步跳转地址
	Attach             string `json:"attach"`                    // 原样返回字段
	BankName           string `json:"bank_name"`                 // 收款银行名称
	BankCardHolderName string `json:"bank_card_holder_name"`     // 银行卡持卡人姓名
	BankCardNo         string `json:"bank_card_no"`              // 银行卡号
	BankBranchName     string `json:"bank_branch_name,optional"` // 收款银行支行名称
}

// 返回示例:
// {"code":"AA300098","msg":"订单号重复"}
type TransferRes struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		CodeUrl string `json:"code_url"`
	}
}

// 下单
func (t *Tpay) Transfer(req TransferReq) (*TransferRes, error) {
	req.CommonRequest.MerchantNo = t.AccNo
	req.CommonRequest.Timestamp = time.Now().Unix()

	dataByte, jerr := json.Marshal(req)
	if jerr != nil {
		return nil, jerr
	}

	data := url.Values{}
	data.Set("data", string(dataByte))
	data.Set("sign", utils.Md5(string(dataByte)+t.SecretKey))

	// post请求
	body, resErr := utils.PostForm(strings.TrimRight(t.Host, "/")+"/transfer", data)
	if resErr != nil {
		return nil, resErr
	}

	fmt.Println("body:", string(body))

	return nil, nil
}

//--------------------------TransferOrderQuery-----------------------------
type TransferOrderQueryReq struct {
	CommonRequest
	MchOrderNo string `json:"mch_order_no"`
	OrderNo    string `json:"order_no"`
}

// 返回示例:
type TransferOrderQueryRes struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		CodeUrl string `json:"code_url"`
	}
}

func (t *Tpay) TransferOrderQuery(req TransferOrderQueryReq) (*TransferOrderQueryRes, error) {
	req.CommonRequest.MerchantNo = t.AccNo
	req.CommonRequest.Timestamp = time.Now().Unix()

	dataByte, jerr := json.Marshal(req)
	if jerr != nil {
		return nil, jerr
	}

	data := url.Values{}
	data.Set("data", string(dataByte))
	data.Set("sign", utils.Md5(string(dataByte)+t.SecretKey))

	// post请求
	body, resErr := utils.PostForm(strings.TrimRight(t.Host, "/")+"/transfer-order-query", data)
	if resErr != nil {
		return nil, resErr
	}

	fmt.Println("body:", string(body))

	return nil, nil
}

//--------------------------QueryBalance-----------------------------------
// 返回示例:
type QueryBalanceRes struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Balance  int64  `json:"balance"`
		Currency string `json:"currency"`
	}
}

func (t *Tpay) QueryBalance() (*PayOrderQueryRes, error) {
	req := CommonRequest{}
	req.MerchantNo = t.AccNo
	req.Timestamp = time.Now().Unix()

	dataByte, jerr := json.Marshal(req)
	if jerr != nil {
		return nil, jerr
	}

	data := url.Values{}
	data.Set("data", string(dataByte))
	data.Set("sign", utils.Md5(string(dataByte)+t.SecretKey))

	// post请求
	body, resErr := utils.PostForm(strings.TrimRight(t.Host, "/")+"/query-balance", data)
	if resErr != nil {
		return nil, resErr
	}

	fmt.Println("body:", string(body))

	return nil, nil
}

//--------------------------------------------------------------------------
