package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"net/url"
	"strings"
)

// 代付API路径
const PayApiTransferPath = "/system/transfer"

// 生成签名算法
func GenerateSign(data, key string) string {
	return strings.ToLower(Md5(data + key))
}

// 代付API请求参数
type TransferApiRequest struct {
	MerchantNo         string `json:"merchant_no"`               // 商户编号
	Timestamp          int64  `json:"timestamp"`                 // 请求时间
	Amount             int64  `json:"amount"`                    // 订单金额
	Currency           string `json:"currency"`                  // 币种
	MchOrderNo         string `json:"mch_order_no"`              // 商户订单号
	TradeType          string `json:"trade_type"`                // 交易类型
	OrderSource        int64  `json:"order_source"`              // 订单来源
	BankName           string `json:"bank_name"`                 // 收款银行名称
	BankCardHolderName string `json:"bank_card_holder_name"`     // 银行卡持卡人姓名
	BankCardNo         string `json:"bank_card_no"`              // 银行卡号
	BankBranchName     string `json:"bank_branch_name,optional"` // 收款银行支行名称
	BankCode           string `json:"bank_code,optional"`        // 银行代码
	Remark             string `json:"remark,optional"`           // 付款备注
}

// 代付API响应数据
type TransferApiResponse struct {
	Code int64            `json:"code"`
	Msg  string           `json:"msg"`
	Data *TransferApiData `json:"data"`
}

type TransferApiData struct {
	MchOrderNo string `json:"mch_order_no"` // 外部订单号(商户系统内部的订单号)
	OrderNo    string `json:"order_no"`     // 平台订单号
	Status     int64  `json:"status"`       // 订单状态
}

// 发起代付请求
func SendTransfer(url string, req TransferApiRequest, secretKey string) (TransferApiResponse, error) {
	postData, err := packTransferParam(req, secretKey)
	if err != nil {
		return TransferApiResponse{}, errors.New(fmt.Sprintf("打包请求参数失败：%v", err))
	}

	body, err := PostForm(url, postData)
	if err != nil {
		return TransferApiResponse{}, errors.New(fmt.Sprintf("发送请求失败：%v", err))
	}

	var resp struct {
		Code int64  `json:"code"`
		Msg  string `json:"msg"`
		Data string `json:"data"`
		Sign string `json:"sign"`
	}
	logx.Errorf("body: %v", string(body))
	if err := json.Unmarshal(body, &resp); err != nil {
		return TransferApiResponse{}, err
	}

	if resp.Code != 0 {
		return TransferApiResponse{}, errors.New(fmt.Sprintf("response code failed,code:%v,msg:%v", resp.Code, resp.Msg))
	}

	data := &TransferApiData{}
	if err := json.Unmarshal([]byte(resp.Data), data); err != nil {
		return TransferApiResponse{}, err
	}

	if data.OrderNo == "" {
		return TransferApiResponse{}, errors.New(fmt.Sprintf("response data is nil, data:%+v", resp.Data))
	}

	return TransferApiResponse{
		Code: resp.Code,
		Msg:  resp.Msg,
		Data: data,
	}, nil

}

// 打包参数(包括签名)
func packTransferParam(param TransferApiRequest, secretKey string) (url.Values, error) {
	content, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	contentStr := string(content)

	sign := GenerateSign(contentStr, secretKey)

	urlValue := url.Values{}
	urlValue.Set("data", contentStr)
	urlValue.Set("sign", sign)

	return urlValue, nil
}

const (
	// 异步通知类型
	PayNotifyType      = "PAY"      // 代收订单通知
	TransferNotifyType = "TRANSFER" // 代付订单通知
)

type PackTransferNotifyParamsRequest struct {
	MerchantNo string `json:"merchant_no"` // 商户编号
	Timestamp  int64  `json:"timestamp"`   // 时间戳
	NotifyType string `json:"notify_type"` // 通知类型

	OrderNo         string `json:"order_no"`     // 平台订单号
	MerchantOrderNo string `json:"mch_order_no"` // 商户订单号
	ReqAmount       int64  `json:"req_amount"`   // 请求金额
	Currency        string `json:"currency"`     // 币种
	OrderStatus     int64  `json:"order_status"` // 订单状态
	PayTime         int64  `json:"pay_time"`     // 支付时间(时间戳)
}

// 打包代付异步通知参数
func PackTransferNotifyParams(req *PackTransferNotifyParamsRequest, secretKey string) (url.Values, error) {
	jByte, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	jsonStr := string(jByte)

	sign := GenerateSign(jsonStr, secretKey)

	urlValue := url.Values{}
	urlValue.Set("data", jsonStr)
	urlValue.Set("sign", sign)
	urlValue.Set("merchant_no", req.MerchantNo)

	return urlValue, nil
}

type PackPayNotifyParamsRequest struct {
	MerchantNo string `json:"merchant_no"` // 商户编号
	Timestamp  int64  `json:"timestamp"`   // 时间戳
	NotifyType string `json:"notify_type"` // 通知类型

	OrderNo         string `json:"order_no"`       // 平台订单号
	MerchantOrderNo string `json:"mch_order_no"`   // 外部订单号
	ReqAmount       int64  `json:"req_amount"`     // 订单请求金额
	Currency        string `json:"currency"`       // 币种
	OrderStatus     int64  `json:"order_status"`   // 订单状态
	PayTime         int64  `json:"pay_time"`       // 支付时间(时间戳)
	Subject         string `json:"subject"`        // 商品名称
	PaymentAmount   int64  `json:"payment_amount"` // 实际支付金额
}

// 打包代收异步通知参数
func PackPayNotifyParams(req *PackPayNotifyParamsRequest, secretKey string) (url.Values, error) {
	jByte, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	jsonStr := string(jByte)

	sign := GenerateSign(jsonStr, secretKey)

	urlValue := url.Values{}
	urlValue.Set("data", jsonStr)
	urlValue.Set("sign", sign)
	urlValue.Set("merchant_no", req.MerchantNo)

	return urlValue, nil
}

//发起批量代付请求

// 批量代付API请求参数
type TransferBatchApiRequest struct {
	MerchantNo string `json:"merchant_no"` // 商户编号
	Timestamp  int64  `json:"timestamp"`   // 请求时间
	BatchNo    string `json:"batch_no"`    // 批量号
}

// 批量代付API响应数据
type TransferBatchApiResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

// 发起代付请求
func SendTransferBatch(url string, req TransferBatchApiRequest, secretKey string) (TransferBatchApiResponse, error) {
	postData, err := packTransferBatchParam(req, secretKey)
	if err != nil {
		return TransferBatchApiResponse{}, errors.New(fmt.Sprintf("打包请求参数失败：%v", err))
	}

	body, err := PostForm(url, postData)
	if err != nil {
		return TransferBatchApiResponse{}, errors.New(fmt.Sprintf("发送请求失败：%v", err))
	}

	var resp struct {
		Code int64  `json:"code"`
		Msg  string `json:"msg"`
		Data string `json:"data"`
		Sign string `json:"sign"`
	}
	logx.Errorf("body: %v", string(body))
	if err := json.Unmarshal(body, &resp); err != nil {
		return TransferBatchApiResponse{}, err
	}

	if resp.Code != 0 {
		return TransferBatchApiResponse{}, errors.New(fmt.Sprintf("response code failed,code:%v,msg:%v", resp.Code, resp.Msg))
	}

	data := &TransferBatchApiResponse{}
	if err := json.Unmarshal([]byte(resp.Data), data); err != nil {
		return TransferBatchApiResponse{}, err
	}

	return TransferBatchApiResponse{
		Code: resp.Code,
		Msg:  resp.Msg,
	}, nil

}

// 打包参数(包括签名)
func packTransferBatchParam(param interface{}, secretKey string) (url.Values, error) {
	content, err := json.Marshal(param)
	if err != nil {
		return nil, err
	}

	contentStr := string(content)

	sign := GenerateSign(contentStr, secretKey)

	urlValue := url.Values{}
	urlValue.Set("data", contentStr)
	urlValue.Set("sign", sign)

	return urlValue, nil
}
