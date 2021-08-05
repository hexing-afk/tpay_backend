package common

//--------------通用的错误定义小于1000---------------------
const (
	Success int = 0 // 成功

	MissingParam      int = 101 // 缺少必须参数
	InvalidParam      int = 102 // 无效的参数
	VerifyParamFailed int = 103 // 参数验证失败

	SystemUnknowErr   int = 201 // 未知错误
	SystemInternalErr int = 202 // 系统内部出错
	SystemBusy        int = 203 // 系统繁忙
)

//-------------业务参数请求大于等于1000----------------
const (
	SignFailed          int = 2001 // 校验签名失败
	DataFieldFormatErr  int = 2002 // data字段格式错误
	MerchantNotExist    int = 2003 // 商户不存在
	MerchantForbidden   int = 2004 // 商户已被禁用
	OrderNotExist       int = 2005 // 订单不存在
	IpWhitListForbidden int = 2006 // IP白名单受限
	DuplicateOrderNO    int = 2007 // 订单号重复
	ChannelUnusable     int = 2008 // 通道不可用
	OrderFailed         int = 2009 // 下单失败
	InsufficientBalance int = 2010 // 余额不足
	NotInTime           int = 2011 // 不在时间内
	NotInAmt            int = 2012 // 不在金额内
)

const (
	SuccessMsg string = "success" // 成功时的msg
)

// 状态码和错误消息映射map
var CodeMap map[int]string

func init() {
	CodeMap = make(map[int]string)
	CodeMap[Success] = SuccessMsg

	CodeMap[SystemUnknowErr] = "未知错误"
	CodeMap[SystemInternalErr] = "系统内部出错"
	CodeMap[SystemBusy] = "系统繁忙"

	CodeMap[MissingParam] = "缺少必须参数"
	CodeMap[InvalidParam] = "无效的参数"
	CodeMap[VerifyParamFailed] = "参数验证失败"

	CodeMap[SignFailed] = "校验签名失败"
	CodeMap[DataFieldFormatErr] = "data字段格式错误"
	CodeMap[MerchantNotExist] = "商户不存在"
	CodeMap[MerchantForbidden] = "商户已被禁用"
	CodeMap[OrderNotExist] = "订单不存在"
	CodeMap[IpWhitListForbidden] = "IP白名单受限"
	CodeMap[DuplicateOrderNO] = "订单号重复"
	CodeMap[ChannelUnusable] = "通道不可用"
	CodeMap[OrderFailed] = "下单失败"
	CodeMap[InsufficientBalance] = "余额不足"
	CodeMap[NotInTime] = "不在时间内"
	CodeMap[NotInAmt] = "不在金额范围内"

}

// 获取一个带code和msg的错误
func NewCodeError(code int) error {
	return &Response{Code: code, Msg: GetCodeMsg(code)}
}

// 获取一个带code和msg的错误
func NewCodeErrorWithMsg(code int, msg string) error {
	return &Response{Code: code, Msg: msg}
}

// 获取错误码对应的描述信息
func GetCodeMsg(code int) string {
	if msg, ok := CodeMap[code]; ok != false { // 有设置默认错误消息
		return msg
	}
	return ""
}
