package upstream

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"strings"
	"tpay_backend/utils"
)

const (
	TopPayQuickSignUrl = "/quickSign"

	// 订单是否已支付
	TopPayUnpaid = "成功"  // 成功
	TopPayPaid   = "失败"  // 失败
	TopPayRefund = "已退款" // 已退款
)

type TopPay struct {
	upMerchantNo string // 上游账号id
	config       TopPayConfig

	payUrl                string
	payOrderQueryUrl      string
	transferUrl           string
	transferOrderQueryUrl string
	queryBalanceUrl       string
	quickSignUrl          string
	quickSignConfirmUrl   string
	quickSignQueryUrl     string
	quickPayUrl           string
	quickPayConfirmUrl    string

	ServiceId string // 接口服务ID
	AppId     string // appid
}

type TopPayResponse struct {
	Code int64  `json:"code"` // 请求返回代码，值见数据字典
	Msg  string `json:"msg"`  // 出错消息，请求处理失败才会出现
	Data string `json:"data"` // 业务数据
	Sign string `json:"sign"` // 签名
}

func (o *TopPay) Pay(request *PayRequest) (*PayResponse, error) {
	return &PayResponse{
		PayUrl: o.config.CashierHost + "/1?o=" + request.OrderNo,
	}, nil
}

func (o *TopPay) PayOrderQuery(request *PayOrderQueryRequest) (*PayOrderQueryResponse, error) {
	panic("implement me")
}

func (o *TopPay) Transfer(req *TransferRequest) (*TransferResponse, error) {
	// 1.拼接参数
	reqParam := make(map[string]interface{})
	reqParam["transactionamount"] = req.Amount / 100.0 //	订单金额
	reqParam["businessnumber"] = req.OrderNo           // 	商户系统内部的订单号
	reqParam["backurl"] = req.NotifyUrl                // 	异步通知地址
	reqParam["subject"] = "商品"                         // 	原样返回字段

	reqParam["bankname"] = req.BankName               // 	收款银行名称
	reqParam["bankcardname"] = req.BankCardHolderName // 	银行卡持卡人姓名
	reqParam["bankcardnumber"] = req.BankCardNo       // 	银行卡号

	reqParam["appid"] = o.config.AppId

	dataByte, jErr := json.Marshal(reqParam)
	if jErr != nil {
		logx.Errorf("map转json失败, err=%v", jErr)
		return nil, jErr
	}

	// 2.发送请求	//调用java sdk
	funcName := utils.RunFuncName()

	logx.Infof(funcName+":request:%v", string(dataByte))
	//return nil, nil

	body, err := utils.PostJson(o.transferUrl, dataByte)
	logx.Infof(funcName+":response:%v", string(body))

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("response body is empty")
	}

	// 3.解析返回结果
	resp := map[string]interface{}{}
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		return nil, errors.New("response body is empty")
	}

	resp2 := map[string]interface{}{}
	if err := json.Unmarshal([]byte(resp["resp"].(string)), &resp2); err != nil {
		return nil, errors.New("response body is empty")
	}

	// 4.验证业务是否成功
	if resp2["data"] == nil || resp2["data"] == "" {
		return nil, errors.New("response body is empty")
	}

	//var dataMap map[string]interface{}
	//if err := json.Unmarshal([]byte(resp2["data"]), &dataMap); err != nil {
	//	return nil, errors.New("response body is empty")
	//}

	logx.Infof("status=%v", resp2["result"])
	var orderStatus int64
	switch resp2["result"] {
	case "success":
		orderStatus = PayPaying // 支付中
	default:
		orderStatus = PayFail // 支付失败
	}

	return &TransferResponse{
		UpstreamOrderNo: utils.ToStringNoPoint(resp2["businessrecordnumber"]),
		OrderStatus:     orderStatus,
	}, nil
}

func (o *TopPay) TransferOrderQuery(request *TransferOrderQueryRequest) (*TransferOrderQueryResponse, error) {
	panic("implement me")
}

func (o *TopPay) QueryBalance() (*QueryBalanceResponse, error) {
	panic("implement me")
}

func (o *TopPay) GenerateSign(data map[string]interface{}) string {
	panic("implement me")
}

func (o *TopPay) CheckSign(data map[string]interface{}) error {
	panic("implement me")
}

//
//type BankCardInfo struct {
//	BankCardNo    string // 银行卡号
//	BankCardName  string // 持卡人名
//	BankCardIdNo  string // 持卡人身份证号
//	BankCardPhone string // 持卡人手机
//	// 信用卡必传
//	BankCardCvv2 string // cvv2
//	// 信用卡必传
//	BankCardExpireDate string // 银行卡有效期
//
//	ClientIp     string // 请求ip
//	IsCreditCard bool   // 是否信用卡
//}

type TopPayConfig struct {
	Host               string `json:"Host"`               // 请求的地址
	PayNotifyPath      string `json:"PayNotifyPath"`      // 请求的地址
	TransferNotifyPath string `json:"TransferNotifyPath"` // 请求的地址
	AppId              string `json:"AppId"`              // 请求的地址
	CashierHost        string `json:"CashierHost"`        // 请求的地址
}

func CheckTopPayConfig(conf TopPayConfig) error {
	if strings.TrimSpace(conf.Host) == "" {
		return errors.New("TopPay.Host配置不能为空")
	}

	// 更多判断...

	return nil
}

func NewTopPay(upMerchantNo string, jsonStrConfig string) (Upstream, error) {
	c := TopPayConfig{}

	// 解析配置
	if err := json.Unmarshal([]byte(jsonStrConfig), &c); err != nil {
		return nil, err
	}

	// 检查配置
	if err := CheckTopPayConfig(c); err != nil {
		return nil, err
	}

	o := &TopPay{}
	o.config = c
	o.upMerchantNo = upMerchantNo

	o.quickSignUrl = strings.TrimRight(c.Host, "/") + TopPayQuickSignUrl
	o.quickSignQueryUrl = strings.TrimRight(c.Host, "/") + "/signQuery"
	o.quickPayUrl = strings.TrimRight(c.Host, "/") + "/quickPay"
	o.quickPayConfirmUrl = strings.TrimRight(c.Host, "/") + "/quickPayConfirm"
	o.quickSignConfirmUrl = strings.TrimRight(c.Host, "/") + "/quickSignConfirm"
	o.transferUrl = strings.TrimRight(c.Host, "/") + "/transfer"

	return o, nil
}

// 获取上游配置
func (o *TopPay) GetUpstreamConfig() *ConfigResponse {
	return &ConfigResponse{
		Host:               o.config.Host,
		PayNotifyPath:      o.config.PayNotifyPath,
		TransferNotifyPath: o.config.TransferNotifyPath,
		SecretKey:          "",
		AppId:              o.config.AppId,
	}
}

func (o *TopPay) QPaySignConfirm(req *QPaySignConfirmRequest) (*QPaySignConfirmResponse, error) {
	// 1.拼接参数
	params := make(map[string]interface{})
	params["businessrecordnumber"] = req.Businessrecordnumber
	params["verifycode"] = req.Verifycode  // String	M  	异步通知回调地址
	params["clientip"] = utils.GetFakeIp() // String	M  	异步通知回调地址

	params["appid"] = o.config.AppId

	dataByte, jErr := json.Marshal(params)
	if jErr != nil {
		return nil, jErr
	}

	funcName := utils.RunFuncName()

	// 2.发送请求
	logx.Infof(funcName+":request:%v", string(dataByte))
	//return nil, nil

	body, err := utils.PostJson(o.quickSignConfirmUrl, dataByte)
	logx.Infof(funcName+":response:%v", string(body))

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("response body is empty")
	}

	// 3.解析返回结果
	resp := map[string]interface{}{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	if err := json.Unmarshal([]byte(resp["resp"].(string)), &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	// 4.验证业务是否成功
	if resp["data"] == nil {
		return nil, errors.New(fmt.Sprintf("response code failed,msg:%v", resp))
	}
	//data := resp["data"].(map[string]interface{})

	var dataMap map[string]interface{}
	if err := json.Unmarshal([]byte(resp["data"].(string)), &dataMap); err != nil {
		return &QPaySignConfirmResponse{
			ErrMsg: fmt.Sprintf("%v|%v", resp["msg"], resp["biz_msg"]),
		}, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	if _, ok := dataMap["businessrecordnumber"]; !ok {
		dataMap["businessrecordnumber"] = ""
	}
	if _, ok := dataMap["businessnumber"]; !ok {
		dataMap["businessnumber"] = ""
	}

	// 5.返回结果
	return &QPaySignConfirmResponse{
		Businessnumber:       dataMap["businessnumber"].(string),
		Businessrecordnumber: dataMap["businessrecordnumber"].(string),
	}, nil
}

func (o *TopPay) QPaySignSms(req *QPaySignSmsRequest) (*QPaySignSmsResponse, error) {
	// 1.拼接参数
	params := make(map[string]interface{})
	params["businessnumber"] = req.OrderNo
	params["bankcardnumber"] = req.BankCardInfo.BankCardNo      // String	M  	异步通知回调地址
	params["bankcardname"] = req.BankCardInfo.BankCardName      // String	M  	异步通知回调地址
	params["certificatenumber"] = req.BankCardInfo.BankCardIdNo // String	M  	异步通知回调地址
	params["bankmobilenumber"] = req.BankCardInfo.BankCardPhone // String	M  	异步通知回调地址
	params["backUrl"] = req.NotifyUrl                           // String	M  	异步通知回调地址
	params["clientip"] = utils.GetFakeIp()                      // String	M  	异步通知回调地址

	params["appid"] = o.config.AppId

	dataByte, jErr := json.Marshal(params)
	if jErr != nil {
		return nil, jErr
	}

	funcName := utils.RunFuncName()

	// 2.发送请求
	logx.Infof(funcName+":request:%v", string(dataByte))
	//return nil, nil

	body, err := utils.PostJson(o.quickSignUrl, dataByte)
	logx.Infof(funcName+":response:%v", string(body))

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("response body is empty")
	}

	// 3.解析返回结果
	resp := map[string]interface{}{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	if err := json.Unmarshal([]byte(resp["resp"].(string)), &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	// 4.验证业务是否成功
	if resp["data"] == nil {
		return nil, errors.New(fmt.Sprintf("response code failed,msg:%v", resp))
	}

	//data := resp["data"].(map[string]interface{})

	var dataMap map[string]interface{}
	if err := json.Unmarshal([]byte(resp["data"].(string)), &dataMap); err != nil {
		return &QPaySignSmsResponse{
			ErrMsg: fmt.Sprintf("%v|%v", resp["msg"], resp["biz_msg"]),
		}, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	if _, ok := dataMap["businessrecordnumber"]; !ok {
		dataMap["businessrecordnumber"] = ""
	}
	if _, ok := dataMap["businessnumber"]; !ok {
		dataMap["businessnumber"] = ""
	}
	if _, ok := dataMap["certcode"]; !ok {
		dataMap["certcode"] = ""
	}

	fmt.Printf("businessrecordnumber=%v", dataMap["businessrecordnumber"].(string))

	// 5.返回结果
	return &QPaySignSmsResponse{
		Businessnumber:       utils.ToStringNoPoint(dataMap["businessnumber"]),
		Businessrecordnumber: utils.ToStringNoPoint(dataMap["businessrecordnumber"]),
		Certcode:             utils.ToStringNoPoint(dataMap["certcode"]),
		ErrMsg:               fmt.Sprintf("%v|%v", resp["msg"], resp["biz_msg"]),
	}, nil
}

// 查询是否已经签约
func (o *TopPay) QPaySignQuery(req *QPaySignQueryRequest) (*QPaySignQueryResponse, error) {
	// 1.拼接参数
	params := make(map[string]interface{})
	params["businessnumber"] = req.BankCardNo

	params["appid"] = o.config.AppId

	dataByte, jErr := json.Marshal(params)
	if jErr != nil {
		return nil, jErr
	}

	funcName := utils.RunFuncName()

	// 2.发送请求
	logx.Infof(funcName+":request:%v", string(dataByte))
	//return nil, nil

	body, err := utils.PostJson(o.quickSignQueryUrl, dataByte)
	logx.Infof(funcName+":response:%v", string(body))

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("response body is empty")
	}

	// 3.解析返回结果
	resp := map[string]interface{}{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	if err := json.Unmarshal([]byte(resp["resp"].(string)), &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	// 4.验证业务是否成功
	if resp["data"] == nil {
		return nil, errors.New(fmt.Sprintf("response code failed,msg:%v", resp))
	}
	//data := resp["data"].(map[string]interface{})

	var dataMap map[string]interface{}
	if err := json.Unmarshal([]byte(resp["data"].(string)), &dataMap); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	if _, ok := dataMap["businessrecordnumber"]; !ok {
		dataMap["businessrecordnumber"] = ""
	}
	if _, ok := dataMap["businessnumber"]; !ok {
		dataMap["businessnumber"] = ""
	}
	if _, ok := dataMap["certcode"]; !ok {
		dataMap["certcode"] = ""
	}

	// 5.返回结果
	return &QPaySignQueryResponse{
		IsSigned:             dataMap["certcode"].(string) != "",
		Businessnumber:       dataMap["businessnumber"].(string),
		Businessrecordnumber: dataMap["businessrecordnumber"].(string),
		Certcode:             dataMap["certcode"].(string),
	}, nil
}

// 下单
func (o *TopPay) QPay(req *QPayRequest) (*QPayResponse, error) {
	// 1.拼接参数
	params := make(map[string]interface{})
	params["backurl"] = req.Backurl
	params["subject"] = req.Subject
	params["businesstype"] = req.Businesstype
	params["kind"] = req.Kind
	params["description"] = req.Description
	params["businessnumber"] = req.Businessnumber
	params["billamount"] = req.Billamount
	params["toaccountnumber"] = req.Toaccountnumber
	params["certcode"] = req.Certcode
	params["clientip"] = utils.GetFakeIp()
	params["merchantuserid"] = req.Merchantuserid

	params["appid"] = o.config.AppId

	dataByte, jErr := json.Marshal(params)
	if jErr != nil {
		return nil, jErr
	}

	funcName := utils.RunFuncName()

	// 2.发送请求
	logx.Infof(funcName+":request:%v", string(dataByte))
	//return nil, nil

	body, err := utils.PostJson(o.quickPayUrl, dataByte)
	logx.Infof(funcName+":response:%v", string(body))

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("response body is empty")
	}

	// 3.解析返回结果
	resp := map[string]interface{}{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	if err := json.Unmarshal([]byte(resp["resp"].(string)), &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	// 4.验证业务是否成功
	if resp["data"] == nil {
		return nil, errors.New(fmt.Sprintf("response code failed,msg:%v", resp))
	}
	//data := resp["data"].(map[string]interface{})

	var dataMap map[string]interface{}
	if err := json.Unmarshal([]byte(resp["data"].(string)), &dataMap); err != nil {
		return &QPayResponse{
			ErrMsg: fmt.Sprintf("%v|%v", resp["msg"], resp["biz_msg"]),
		}, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	if _, ok := dataMap["businessrecordnumber"]; !ok {
		dataMap["businessrecordnumber"] = ""
	}
	if _, ok := dataMap["businessnumber"]; !ok {
		dataMap["businessnumber"] = ""
	}

	// 5.返回结果
	return &QPayResponse{
		Businessnumber:       dataMap["businessnumber"].(string),
		Businessrecordnumber: dataMap["businessrecordnumber"].(string),
	}, nil
}

// 下单确认支付
func (o *TopPay) QPayConfirm(req *QPayConfirmRequest) (*QPayConfirmResponse, error) {
	// 1.拼接参数
	params := make(map[string]interface{})
	params["businessrecordnumber"] = req.Businessrecordnumber
	params["verifycode"] = req.Verifycode
	params["clientip"] = utils.GetFakeIp()

	params["appid"] = o.config.AppId

	dataByte, jErr := json.Marshal(params)
	if jErr != nil {
		return nil, jErr
	}

	funcName := utils.RunFuncName()

	// 2.发送请求
	logx.Infof(funcName+":request:%v", string(dataByte))
	//return nil, nil

	body, err := utils.PostJson(o.quickPayConfirmUrl, dataByte)
	logx.Infof(funcName+":response:%v", string(body))

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("response body is empty")
	}

	// 3.解析返回结果
	resp := map[string]interface{}{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	if err := json.Unmarshal([]byte(resp["resp"].(string)), &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	// 4.验证业务是否成功
	if resp["data"] == nil {
		return nil, errors.New(fmt.Sprintf("response code failed,msg:%v", resp))
	}

	//data := resp["data"].(map[string]interface{})

	var dataMap map[string]interface{}
	if err := json.Unmarshal([]byte(resp["data"].(string)), &dataMap); err != nil {
		return &QPayConfirmResponse{
			ErrMsg: fmt.Sprintf("%v|%v", resp["msg"], resp["biz_msg"]),
		}, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	if _, ok := dataMap["businessrecordnumber"]; !ok {
		dataMap["businessrecordnumber"] = ""
	}
	if _, ok := dataMap["businessnumber"]; !ok {
		dataMap["businessnumber"] = ""
	}

	// 5.返回结果
	return &QPayConfirmResponse{
		Businessnumber:       dataMap["businessnumber"].(string),
		Businessrecordnumber: dataMap["businessrecordnumber"].(string),
	}, nil
}
