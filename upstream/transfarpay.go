package upstream

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"io/ioutil"
	"software.sslmate.com/src/go-pkcs12"
	"strings"
	"time"
	"tpay_backend/utils"
)

const (
	// 请求PATH
	TransfarPayPayUrl                = "/api/create"                // 代收下单
	TransfarPayPayOrderQueryUrl      = "/api/getpayinfo?ticket="    // 代收订单查询
	TransfarPayTransferUrl           = "/api/withdrawal"            // 代付下单
	TransfarPayTransferOrderQueryUrl = "/api/getwithdrawal?ticket=" // 代付订单查询
	TransfarPayQueryBalance          = ""                           // 查询余额

	TransferPayQuickSignUrl = "tf56pay.gateway.quickSignConfirm" // 签约

	// 签名字段
	TransfarPaySignFiledName = "sign"

	// 订单是否已支付
	TransfarPayUnpaid = 0 // 未支付
	TransfarPayPaid   = 1 // 已支付

	// 订单是否已取消
	TransfarPayNotCancelled = 0 // 未取消
	TransfarPayCancelled    = 1 // 已取消
)

type TransfarPay struct {
	upMerchantNo string // 上游账号id
	config       TransfarPayConfig

	payUrl                string
	payOrderQueryUrl      string
	transferUrl           string
	transferOrderQueryUrl string
	queryBalanceUrl       string

	ServiceId string // 接口服务ID
	AppId     string // appid
}

type BankCardInfo struct {
	BankCardNo    string // 银行卡号
	BankCardName  string // 持卡人名
	BankCardIdNo  string // 持卡人身份证号
	BankCardPhone string // 持卡人手机
	// 信用卡必传
	BankCardCvv2 string // cvv2
	// 信用卡必传
	BankCardExpireDate string // 银行卡有效期

	ClientIp     string // 请求ip
	IsCreditCard bool   // 是否信用卡
}

type TransfarPayConfig struct {
	Host               string `json:"Host"`               // 请求的地址
	PfxFile            string `json:"PfxFile"`            // 商家通信秘钥
	PfxSecret          string `json:"PfxSecret"`          // 密钥密码
	PayNotifyPath      string `json:"PayNotifyPath"`      // 代收异步通知path路径
	TransferNotifyPath string `json:"TransferNotifyPath"` // 代付异步通知path路径
	AppId              string `json:"AppId"`              // 代付异步通知path路径
}

func CheckTransfarPayConfig(conf TransfarPayConfig) error {
	if strings.TrimSpace(conf.Host) == "" {
		return errors.New("TransfarPay.Host配置不能为空")
	}

	if strings.TrimSpace(conf.PfxFile) == "" {
		return errors.New("TransfarPay.PfxFile配置不能为空")
	}

	if strings.TrimSpace(conf.PfxSecret) == "" {
		return errors.New("TransfarPay.PfxSecret配置不能为空")
	}

	//if strings.TrimSpace(conf.PayNotifyPath) == "" {
	//	return errors.New("TransfarPay.PayNotifyPath配置不能为空")
	//}
	//
	//if strings.TrimSpace(conf.TransferNotifyPath) == "" {
	//	return errors.New("TransfarPay.TransferNotifyPath配置不能为空")
	//}

	// 更多判断...

	return nil
}

func NewTransfarPay(upMerchantNo string, jsonStrConfig string) (Upstream, error) {
	c := TransfarPayConfig{}

	if strings.TrimSpace(upMerchantNo) == "" {
		return nil, errors.New("upMerchantNo配置不能为空")
	}

	// 解析配置
	if err := json.Unmarshal([]byte(jsonStrConfig), &c); err != nil {
		return nil, err
	}

	// 检查配置
	if err := CheckTransfarPayConfig(c); err != nil {
		return nil, err
	}

	o := &TransfarPay{}
	o.config = c
	o.upMerchantNo = upMerchantNo

	o.payUrl = strings.TrimRight(c.Host, "/") + TransfarPayPayUrl
	o.payOrderQueryUrl = strings.TrimRight(c.Host, "/") + TransfarPayPayOrderQueryUrl
	o.transferUrl = strings.TrimRight(c.Host, "/") + TransfarPayTransferUrl
	o.transferOrderQueryUrl = strings.TrimRight(c.Host, "/") + TransfarPayTransferOrderQueryUrl
	o.queryBalanceUrl = strings.TrimRight(c.Host, "/") + TransfarPayQueryBalance

	return o, nil
}

// 获取上游配置
func (o *TransfarPay) GetUpstreamConfig() *ConfigResponse {
	return &ConfigResponse{
		PayNotifyPath:      o.config.PayNotifyPath,
		TransferNotifyPath: o.config.TransferNotifyPath,
		SecretKey:          o.config.PfxSecret,
	}
}

func (o *TransfarPay) Pay(req *PayRequest) (*PayResponse, error) {
	// 1.拼接参数
	params := make(map[string]interface{})
	params["service_id"] = o.ServiceId                           // String	M	商户订单号，平台不会判断orderid是否重复，该字段平台主要用作签名
	params["appid"] = o.AppId                                    // String	M	商户订单号，平台不会判断orderid是否重复，该字段平台主要用作签名
	params["tf_timestamp"] = time.Now().Format("20060102150405") // int	M  	金额
	params["sign_type"] = "RSA"                                  // String	M	商户编号，商户ID平台提供
	params["terminal"] = "Android"                               // String	M 	商户订单类型

	params["backurl"] = req.NotifyUrl      // String	M  	异步通知回调地址
	params["subject"] = "装修服务"             // String	M  	异步通知回调地址
	params["businesstype"] = "商业服务费用"      // String	M  	异步通知回调地址
	params["kind"] = "商业服务费用"              // String	M  	异步通知回调地址
	params["description"] = "商业服务费用"       // String	M  	异步通知回调地址
	params["businessnumber"] = req.OrderNo // String	O  	商户备注，平台通知或商户查询时原文返回
	params["billamount"] = req.Amount      // String	O  	商户备注，平台通知或商户查询时原文返回
	params["toaccountnumber"] = ""         // String	O  	商户备注，平台通知或商户查询时原文返回

	params["merchantuserid"] = o.upMerchantNo // String	M  	签名

	params["tf_sign"] = o.GenerateSign(params) // String	M 	商户订单类型
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
		Success  int64  `json:"success"`  // 请求是否成功，1 成功 0 失败
		Message  string `json:"message"`  // 出错消息，请求处理失败才会出现
		Ticket   string `json:"ticket"`   // 订单访问票据或标识 商户系统应该保存，用于平台通知时分辨订单或查询时作为主参数
		OrderId  string `json:"orderid"`  // 商户订单号
		UserId   string `json:"userid"`   // 商户编号
		PageUrl  string `json:"pageurl"`  // 商户已经直接在前台跳转到该地址，打开平台提供的支付页面
		SerialNo string `json:"serialno"` // 需要最终用户在支付是必须输入的备注信息。需在支付页面展示给用户。4个数字，从1000到9999
		BMount   string `json:"bmount"`   // 需要最终用户在支付是必须输入的尾部金额信息信息，需在支付页面展示给用户。
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	// 4.验证业务是否成功
	if resp.Success != 1 {
		return nil, errors.New(fmt.Sprintf("response code failed,code:%v,msg:%v", resp.Success, resp.Message))
	}

	// 5.返回结果
	return &PayResponse{
		UpstreamOrderNo: resp.Ticket,
		PayUrl:          resp.PageUrl,
	}, nil
}

func (o *TransfarPay) PayOrderQuery(req *PayOrderQueryRequest) (*PayOrderQueryResponse, error) {
	// 1.拼接参数
	url := o.payOrderQueryUrl + req.UpstreamOrderNo

	funcName := utils.RunFuncName()

	// 2.发送请求
	logx.Infof(funcName+":request:%v", url)
	body, err := utils.Get(url)
	logx.Infof(funcName+":response:%v", string(body))

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("response body is empty")
	}

	var resp struct {
		Success   int64  `json:"success"`   // 请求是否成功 1、成功；0、失败
		Message   string `json:"message"`   // 出错消息，请求处理失败才会出现
		Ticket    string `json:"ticket"`    // 访问票据
		IsPay     int64  `json:"ispay"`     // 是否支付，0 没有支付 1 已经支付
		PayCode   string `json:"paycode"`   // 支付代码	支付网关返回编码
		PayAmount int64  `json:"payamount"` // 支付金额 支付网关返回的实际金额，业务逻辑中应使用此金额作为入金金额而非定单金额
		PayTime   string `json:"msg"`       // 支付时间	字符串类型格式为： 2000-01-01 23:34:56
		PayUser   string `json:"status"`    // 支付用户
		Sign      string `json:"sign"`      // 签名
		Amount    int64  `json:"amount"`    // 创建订单时的金额，原样返回
		Note      string `json:"note"`      // 创建订单时的备注，原样返回
		UserId    string `json:"userid"`    // 商户编号
		OrderId   string `json:"orderid"`   // 商户订单号
		PayType   string `json:"type"`      // 支付类型
		SerialNo  string `json:"serialno"`  // 支付备注 需要最终用户在支付是必须输入的备注信息。需在支付页面展示给用户。4个数字，从1000到9999
		BMount    string `json:"bmount"`    // 尾部金额 需要最终用户在支付是必须输入的尾部金额信息信息，需在支付页面展示给用户。
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	// 4.验证业务是否成功
	if resp.Success != 1 {
		return nil, errors.New(fmt.Sprintf("response code failed,code:%v,msg:%v", resp.Success, resp.Message))
	}

	// 5.验签
	checkSignMap := make(map[string]interface{})
	checkSignMap["orderid"] = resp.OrderId
	checkSignMap["amount"] = resp.PayAmount
	checkSignMap["sign"] = resp.Sign

	if err := o.CheckSign(checkSignMap); err != nil {
		return nil, errors.New(fmt.Sprintf("sign verification failed, err:%v, body:%v", err, string(body)))
	}

	var orderStatus int64
	switch resp.IsPay {
	case TransfarPayUnpaid:
		orderStatus = PayPaying
	case TransfarPayPaid:
		orderStatus = PaySuccess
	default:
		orderStatus = PayPaying
	}

	return &PayOrderQueryResponse{
		OrderNo:         resp.OrderId,
		UpstreamOrderNo: resp.Ticket,
		OrderStatus:     orderStatus,
	}, nil
}

func (o *TransfarPay) Transfer(req *TransferRequest) (*TransferResponse, error) {
	//// 1.拼接参数
	//payload := make(map[string]interface{})
	//payload["cardname"] = req.BankCardHolderName // 收款人姓名
	//payload["cardno"] = req.BankCardNo           // 收款卡号
	//payload["bankname"] = req.BankName           // 银行名称
	//payload["bankid"] = req.BankCode             // 银行编号
	//payload["province"] = ""                     // 银行所在省（非必传）
	//payload["city"] = ""                         // 银行所在市（非必传）
	//payload["branchname"] = ""                   // 支行名称（非必传）
	//payload["ifsc"] = ""                         // IFSC CODE（用于印度）
	//
	//payloadByte, jErr := json.Marshal(payload)
	//if jErr != nil {
	//	return nil, jErr
	//}
	//
	//params := make(map[string]interface{})
	//params["orderid"] = req.OrderNo         // String	M	商户订单号，平台不会判断orderid是否重复，该字段平台主要用作签名
	//params["amount"] = req.Amount           // int	M  	金额
	//params["sign"] = o.GenerateSign(params) // String	M  	签名
	//params["userid"] = o.upMerchantNo       // String	M	商户编号，商户ID平台提供
	//params["type"] = req.ProductType        // String	M 	商户订单类型
	//params["notifyurl"] = req.NotifyUrl     // String	M  	异步通知回调地址
	//params["returnurl"] = req.ReturnUrl     // String	O  	最终用户支付后，平台支付页面跳转到商户的地址
	//params["note"] = req.Attach             // String	O  	商户备注，平台通知或商户查询时原文返回
	//params["payload"] = string(payloadByte) // String	O  	订单其他数据，json格式字符串
	//
	//dataByte, jErr := json.Marshal(params)
	//if jErr != nil {
	//	return nil, jErr
	//}
	//
	//funcName := utils.RunFuncName()
	//
	//// 2.发送请求
	//logx.Infof(funcName+":request:%v", string(dataByte))
	//body, err := utils.PostJson(o.transferUrl, dataByte)
	//logx.Infof(funcName+":response:%v", string(body))
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//if len(body) == 0 {
	//	return nil, errors.New("response body is empty")
	//}
	//
	//// 3.解析返回结果
	//var resp struct {
	//	Success  int64  `json:"success"`  // 请求是否成功，1 成功 0 失败
	//	Message  string `json:"message"`  // 出错消息，请求处理失败才会出现
	//	Ticket   string `json:"ticket"`   // 订单访问票据或标识 商户系统应该保存，用于平台通知时分辨订单或查询时作为主参数
	//	OrderId  string `json:"orderid"`  // 商户订单号
	//	UserId   string `json:"userid"`   // 商户编号
	//	PageUrl  string `json:"pageurl"`  // 商户已经直接在前台跳转到该地址，打开平台提供的支付页面
	//	SerialNo string `json:"serialno"` // 需要最终用户在支付是必须输入的备注信息。需在支付页面展示给用户。4个数字，从1000到9999
	//	BMount   string `json:"bmount"`   // 需要最终用户在支付是必须输入的尾部金额信息信息，需在支付页面展示给用户。
	//}
	//
	//if err := json.Unmarshal(body, &resp); err != nil {
	//	return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	//}
	//
	//// 4.验证业务是否成功
	//if resp.Success != 1 {
	//	return nil, errors.New(fmt.Sprintf("response code failed,code:%v,msg:%v", resp.Success, resp.Message))
	//}

	//return &TransferResponse{
	//	UpstreamOrderNo: resp.Ticket,
	//	OrderStatus:     TransferPaying,
	//}, nil
	return nil, nil
}

func (o *TransfarPay) TransferOrderQuery(req *TransferOrderQueryRequest) (*TransferOrderQueryResponse, error) {
	// 1.拼接参数
	url := o.transferOrderQueryUrl + req.UpstreamOrderNo

	funcName := utils.RunFuncName()

	// 2.发送请求
	logx.Infof(funcName+":request:%v", url)
	body, err := utils.Get(url)
	logx.Infof(funcName+":response:%v", string(body))

	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return nil, errors.New("response body is empty")
	}

	var resp struct {
		Success   int64  `json:"success"`   // 请求是否成功 1、成功；0、失败
		Message   string `json:"message"`   // 出错消息，请求处理失败才会出现
		Ticket    string `json:"ticket"`    // 访问票据
		IsPay     int64  `json:"ispay"`     // 是否支付，0 没有支付 1 已经支付
		PayCode   string `json:"paycode"`   // 支付代码	支付网关返回编码
		PayAmount int64  `json:"payamount"` // 支付金额 支付网关返回的实际金额，业务逻辑中应使用此金额作为入金金额而非定单金额
		PayTime   string `json:"msg"`       // 支付时间	字符串类型格式为： 2000-01-01 23:34:56
		PayUser   string `json:"status"`    // 支付用户
		Sign      string `json:"sign"`      // 签名
		Amount    int64  `json:"amount"`    // 创建订单时的金额，原样返回
		Note      string `json:"note"`      // 创建订单时的备注，原样返回
		UserId    string `json:"userid"`    // 商户编号
		OrderId   string `json:"orderid"`   // 商户订单号
		PayType   string `json:"type"`      // 支付类型
		SerialNo  string `json:"serialno"`  // 支付备注 需要最终用户在支付是必须输入的备注信息。需在支付页面展示给用户。4个数字，从1000到9999
		BMount    string `json:"bmount"`    // 尾部金额 需要最终用户在支付是必须输入的尾部金额信息信息，需在支付页面展示给用户。
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, errors.New(fmt.Sprintf("parse json body failed, err:%v, body:%v", err, string(body)))
	}

	// 4.验证业务是否成功
	if resp.Success != 1 {
		return nil, errors.New(fmt.Sprintf("response code failed,code:%v,msg:%v", resp.Success, resp.Message))
	}

	// 5.验签
	checkSignMap := make(map[string]interface{})
	checkSignMap["orderid"] = resp.OrderId
	checkSignMap["amount"] = resp.PayAmount
	checkSignMap["sign"] = resp.Sign

	if err := o.CheckSign(checkSignMap); err != nil {
		return nil, errors.New(fmt.Sprintf("sign verification failed, err:%v, body:%v", err, string(body)))
	}

	var orderStatus int64
	switch resp.IsPay {
	case TransfarPayUnpaid:
		orderStatus = TransferPaying
	case TransfarPayPaid:
		orderStatus = PaySuccess
	default:
		orderStatus = TransferPaying
	}

	return &TransferOrderQueryResponse{
		OrderNo:         resp.OrderId,
		UpstreamOrderNo: resp.Ticket,
		OrderStatus:     orderStatus,
	}, nil
}

func (o *TransfarPay) QueryBalance() (*QueryBalanceResponse, error) {
	return &QueryBalanceResponse{
		Balance:           0,
		PayOutBalance:     0,
		PayAmountLimit:    0,
		PayoutAmountLimit: 0,
		Currency:          "",
	}, nil
}

// 生成签名
func (o *TransfarPay) GenerateSign(data map[string]interface{}) string {
	// 将参数排序并拼接成字符串
	dataStr := utils.ParamsMapToString(data, "tf_sign")
	dataStr = fmt.Sprintf("%s&", dataStr)

	priv, _, err := getKeys(o.config.PfxFile, o.config.PfxSecret)
	if err != nil {
		logx.Errorf("err=%v", err)
		return ""
	}

	message := []byte("message to be signed")
	hashed := sha1.Sum(message)

	signature, err := rsa.SignPKCS1v15(rand.Reader, priv, crypto.SHA1, hashed[:])
	if err != nil {
		logx.Errorf("err=%v", err)
		return ""
	}

	signatureBase64 := base64.StdEncoding.EncodeToString(signature)
	return signatureBase64
}

// 校验签名
func (o *TransfarPay) CheckSign(data map[string]interface{}) error {
	sign, exist := data[TransfarPaySignFiledName]
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

func getKeys(privateKeyName, privatePassword string) (*rsa.PrivateKey, *x509.Certificate, error) {
	// 因为pfx证书公钥和密钥是成对的，所以要先转成pem.Block
	b, err := ioutil.ReadFile(privateKeyName)
	if err != nil {
		return nil, nil, err
	}
	privKey, pubkey, _, err := pkcs12.DecodeChain(b, privatePassword)
	if err != nil {
		return nil, nil, err
	}
	return privKey.(*rsa.PrivateKey), pubkey, nil
}
