package upstream

//-------------------上游配置-----------------------
type ConfigResponse struct {
	Host               string
	PayNotifyPath      string
	TransferNotifyPath string
	SecretKey          string
	AppId              string
}

const (
	PayPaying  int64 = 1 // 支付中
	PaySuccess int64 = 2 // 成功
	PayFail    int64 = 3 // 失败
)

//-------------------代收---------------------------
type PayRequest struct {
	Amount      int64  // 请求金额
	Currency    string // 币种
	OrderNo     string // 平台订单号
	NotifyUrl   string // 异步通知地址
	ReturnUrl   string // 跳转地址
	ProductType string // 交易方式
	Subject     string // 商品的标题/交易标题/订单标题/订单关键字等

	CustomName   string // 客户姓名
	CustomMobile string // 客户电话
	CustomEmail  string // 客户邮箱
	Attach       string // 原样返回字段

	BankCardInfo BankCardInfo
}

type PayResponse struct {
	UpstreamOrderNo string // 上游订单号
	PayUrl          string // 支付连接
}

//-------------------快捷签约---------------------------
type QPaySignConfirmRequest struct {
	Businessrecordnumber string
	Verifycode           string
	Clientip             string
}

type QPaySignConfirmResponse struct {
	Businessnumber       string
	Businessrecordnumber string
	ErrMsg               string
}

//-------------------快捷签约---------------------------
type QPaySignSmsRequest struct {
	Amount      int64  // 请求金额
	Currency    string // 币种
	OrderNo     string // 平台订单号
	NotifyUrl   string // 异步通知地址
	ReturnUrl   string // 跳转地址
	ProductType string // 交易方式
	Subject     string // 商品的标题/交易标题/订单标题/订单关键字等

	CustomName   string // 客户姓名
	CustomMobile string // 客户电话
	CustomEmail  string // 客户邮箱
	Attach       string // 原样返回字段

	BankCardInfo BankCardInfo
}

type QPaySignSmsResponse struct {
	Businessnumber       string
	Businessrecordnumber string
	Certcode             string
	ErrMsg               string
}

//-------------------快捷签约---------------------------
type QPaySignQueryRequest struct {
	BankCardNo string
}

type QPaySignQueryResponse struct {
	IsSigned             bool //
	Businessnumber       string
	Businessrecordnumber string
	Certcode             string
}

//-------------------快捷签约---------------------------
type QPayRequest struct {
	Backurl         string
	Subject         string
	Businesstype    string
	Kind            string
	Description     string
	Businessnumber  string
	Billamount      string
	Toaccountnumber string
	Certcode        string
	Clientip        string
	Merchantuserid  string
}

type QPayResponse struct {
	Businessnumber       string
	Businessrecordnumber string
	ErrMsg               string
}

//-------------------快捷签约---------------------------
type QPayConfirmRequest struct {
	Businessrecordnumber string
	Verifycode           string
	Clientip             string
}

type QPayConfirmResponse struct {
	Businessnumber       string
	Businessrecordnumber string
	ErrMsg               string
}

//-------------------代收订单查询---------------------------
type PayOrderQueryRequest struct {
	OrderNo         string
	UpstreamOrderNo string
}

type PayOrderQueryResponse struct {
	OrderNo         string
	UpstreamOrderNo string
	OrderStatus     int64
}

const (
	TransferPaying  int64 = 1 // 支付中
	TransferSuccess int64 = 2 // 成功
	TransferFail    int64 = 3 // 失败
)

//-------------------代付---------------------------
type TransferRequest struct {
	Amount      int64
	Currency    string
	OrderNo     string
	NotifyUrl   string
	ReturnUrl   string
	ProductType string
	Attach      string
	Remark      string

	BankName           string // 银行名称
	BankBranchName     string // 银行支行名称(开户行)
	BankCardNo         string // 银行卡号
	BankCode           string // 银行代码
	BankCardHolderName string // 银行卡持卡人姓名
	CardHolderMobile   string // 持卡人电话
	CardHolderEmail    string // 持卡人邮箱
}

type TransferResponse struct {
	UpstreamOrderNo string
	OrderStatus     int64
}

//-------------------代付订单查询---------------------------
type TransferOrderQueryRequest struct {
	OrderNo         string
	UpstreamOrderNo string
}

type TransferOrderQueryResponse struct {
	OrderNo         string
	UpstreamOrderNo string
	OrderStatus     int64
}

//-------------------查询余额---------------------------
type QueryBalanceResponse struct {
	Balance           float64 // 余额
	PayOutBalance     float64 // 代付余额
	PayAmountLimit    float64 // 今日剩余代收额度
	PayoutAmountLimit float64 // 今日剩余代付额度
	Currency          string  // 币种
}

//------------------代付异步通知-------------------------
type TransferNotifyRequest struct {
	ReqData map[string]interface{}
}

type TransferNotifyResponse struct {
	OrderNo         string
	UpstreamOrderNo string
	OrderStatus     int64
}

type Upstream interface {
	// 获取上游配置
	GetUpstreamConfig() *ConfigResponse

	// 代收
	Pay(*PayRequest) (*PayResponse, error)

	// 代收订单查询
	PayOrderQuery(*PayOrderQueryRequest) (*PayOrderQueryResponse, error)

	// 代付
	Transfer(*TransferRequest) (*TransferResponse, error)

	// 代付订单查询
	TransferOrderQuery(*TransferOrderQueryRequest) (*TransferOrderQueryResponse, error)

	// 查询余额
	QueryBalance() (*QueryBalanceResponse, error)

	// 签名
	GenerateSign(data map[string]interface{}) string

	// 校验签名
	CheckSign(data map[string]interface{}) error
}
