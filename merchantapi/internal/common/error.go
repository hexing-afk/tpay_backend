package common

//--------------通用的错误定义小于1000---------------------
const (
	Success int = 0 // 成功

	MissingParam      int = 101 // 缺少必须参数
	InvalidParam      int = 102 // 无效的参数
	VerifyParamFailed int = 103 // 参数验证失败
	UploadFail        int = 104 // 上传失败

	SysUnKnow         int = 201 // 未知错误
	SysDBErr          int = 202 // 数据库操作失败(系统内部错误)
	SysDBGet          int = 203 // 获取数据失败
	SysDBAdd          int = 204 // 添加数据失败
	SysDBUpdate       int = 205 // 更新数据失败
	SysDBDelete       int = 206 // 删除数据失败
	SystemInternalErr int = 207 // 系统内部错误

	UserNotLogin int = 400 // 用户未登录

)

//-------------业务参数请求大于等于1000----------------
const (
	GetLoginCaptchaFailed    int = 1001 // 获取验证码失败
	LoginCaptchaNotMatch     int = 1002 // 验证码错误
	UserLoginFailed          int = 1003 // 用户登录失败
	UserJwtVerifyFailed      int = 1004 // 登录jwt验证失败
	AccountRepeat            int = 1005 // 账号重复
	UserNotExist             int = 1006 // 用户不存在
	GetLoginUserInfoFailed   int = 1007 // 获取登录用户信息失败
	RefreshTokenFailed       int = 1008 // 刷新token失败
	NoLoginToken             int = 1009 // 缺少登录token
	UserLogoutFailed         int = 1010 // 退出登录失败
	LoginTokenParseFailed    int = 1011 // 解析登录token失败
	ReSetMd5KeyFailed        int = 1012 // 刷新Md5失败
	ReSetLoginPwdFailed      int = 1013 // 修改登录密码失败
	LoginPwdFailed           int = 1014 // 登录密码错误
	UpdatePayPwdFailed       int = 1015 // 修改支付密码失败
	BankCardRepeatAdd        int = 1016 // 银行卡已经绑定
	AccountDisable           int = 1017 // 账号已被禁用
	NotWithdrawConfig        int = 1018 // 没有提现配置
	PayPasswordErr           int = 1019 // 支付密码错误
	BankCardNotExist         int = 1020 // 银行卡不存在
	BankCardUnusable         int = 1021 // 银行卡不可用
	AmountOutOfLimit         int = 1022 // 金额超出限制
	InsufficientBalance      int = 1023 // 余额不足
	RechargeFailed           int = 1024 // 申请充值失败
	HankCardNotExist         int = 1025 // 银行卡不存在
	ApplyFail                int = 1026 // 申请失败
	BankCardMaxAmountLacking int = 1027 // 银行卡收款额度不足
	AmountFail               int = 1028 // 金额错误
	NotChannelAvailable      int = 1029 // 没有可用通道
	OrderNotExist            int = 1030 // 订单不存在
	OrderNotOp               int = 1031 // 订单不能执行此操作
	OrderMissingNotifyUrl    int = 1032 // 订单缺少异步通知地址
	OrderNotifyFail          int = 1033 // 订单通知失败
	ExportFail               int = 1034 // 导出失败
	NotData                  int = 1035 // 没有数据
	ExportLimit31Day         int = 1036 // 最长导出31天的数据
	CheckExportTime          int = 1037 // 请选择导出时间
	NotSupportUploadFile     int = 1038 // 不支持的上传文件格式
	FileNotExist             int = 1039 // 文件不存在
	TransferFail             int = 1040 // 付款失败
)

const (
	SuccessMsg string = "success" // 成功时的msg
)

// 状态码和错误消息映射map
var CodeMap map[int]string

func init() {
	CodeMap = make(map[int]string)
	CodeMap[Success] = SuccessMsg

	CodeMap[MissingParam] = "缺少必须参数"
	CodeMap[InvalidParam] = "无效的参数"
	CodeMap[VerifyParamFailed] = "参数验证失败"
	CodeMap[UploadFail] = "上传失败"

	CodeMap[SysUnKnow] = "未知错误"
	CodeMap[SysDBGet] = "获取数据失败"
	CodeMap[SysDBAdd] = "添加数据失败"
	CodeMap[SysDBUpdate] = "更新数据失败"
	CodeMap[SysDBDelete] = "删除数据失败"
	CodeMap[SysDBErr] = "数据库操作失败(系统内部错误)"
	CodeMap[SystemInternalErr] = "系统内部错误"

	CodeMap[GetLoginCaptchaFailed] = "获取验证码失败"
	CodeMap[LoginCaptchaNotMatch] = " 验证码错误"
	CodeMap[UserLoginFailed] = "用户登录失败"
	CodeMap[UserJwtVerifyFailed] = "登录jwt验证失败"
	CodeMap[AccountRepeat] = "账号重复"
	CodeMap[UserNotExist] = "用户不存在"
	CodeMap[GetLoginUserInfoFailed] = "获取登录用户信息失败"

	CodeMap[RefreshTokenFailed] = "刷新token失败"
	CodeMap[NoLoginToken] = "缺少登录token"
	CodeMap[UserNotLogin] = "用户未登录"
	CodeMap[UserLogoutFailed] = "退出登录失败"
	CodeMap[LoginTokenParseFailed] = "解析登录token失败"
	CodeMap[ReSetMd5KeyFailed] = "刷新Md5失败"
	CodeMap[ReSetLoginPwdFailed] = "修改登录密码失败"
	CodeMap[LoginPwdFailed] = "登录密码错误"
	CodeMap[UpdatePayPwdFailed] = "修改支付密码失败"
	CodeMap[BankCardRepeatAdd] = "银行卡已经绑定"
	CodeMap[NotWithdrawConfig] = "没有提现配置"
	CodeMap[PayPasswordErr] = "支付密码错误"
	CodeMap[BankCardNotExist] = "银行卡不存在"
	CodeMap[BankCardUnusable] = "银行卡不可用"
	CodeMap[AmountOutOfLimit] = "金额超出限制"
	CodeMap[InsufficientBalance] = "余额不足"
	CodeMap[RechargeFailed] = "申请充值失败"
	CodeMap[HankCardNotExist] = "银行卡不存在"
	CodeMap[ApplyFail] = "申请失败"
	CodeMap[AccountDisable] = "账号已被禁用"
	CodeMap[BankCardMaxAmountLacking] = "银行卡收款额度不足"
	CodeMap[AmountFail] = "金额错误"
	CodeMap[NotChannelAvailable] = "没有可用通道"
	CodeMap[OrderNotExist] = "订单不存在"
	CodeMap[OrderNotOp] = "订单不能执行此操作"
	CodeMap[OrderMissingNotifyUrl] = "订单缺少异步通知地址"
	CodeMap[OrderNotifyFail] = "订单通知失败"
	CodeMap[ExportFail] = "导出失败"
	CodeMap[NotData] = "没有数据"
	CodeMap[ExportLimit31Day] = "最长导出31天的数据"
	CodeMap[CheckExportTime] = "请选择导出时间"
	CodeMap[NotSupportUploadFile] = "不支持的上传文件格式"
	CodeMap[FileNotExist] = "文件不存在"
	CodeMap[TransferFail] = "付款失败"

	//CodeMap[] = ""
	//CodeMap[] = ""
	//CodeMap[] = ""

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
