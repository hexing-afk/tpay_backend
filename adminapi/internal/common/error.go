package common

//--------------通用的错误定义小于1000---------------------
const (
	Success int = 0 // 成功

	MissingParam      int = 101 // 缺少必须参数
	InvalidParam      int = 102 // 无效的参数
	VerifyParamFailed int = 103 // 参数验证失败
	UploadFail        int = 104 // 上传失败

	SystemUnknownErr  int = 201 // 未知错误
	SystemInternalErr int = 202 // 系统内部错误
	SystemBusy        int = 203 // 系统繁忙
	SysDBErr          int = 204 // 数据库操作失败(系统内部错误)
	SysDBGet          int = 205 // 获取数据失败
	SysDBAdd          int = 206 // 添加数据失败
	SysDBUpdate       int = 207 // 更新数据失败
	SysDBDelete       int = 208 // 删除数据失败
	SysDBSave         int = 209 // 保存数据失败

	UserNotLogin int = 400 // 用户未登录

)

//-------------业务参数请求大于等于1000----------------
const (
	GetLoginCaptchaFailed             int = 1001 // 获取验证码失败
	LoginCaptchaNotMatch              int = 1002 // 验证码错误
	UserLoginFailed                   int = 1003 // 用户登录失败
	AccountRepeat                     int = 1004 // 账号重复
	UserNotExist                      int = 1005 // 用户不存在
	GetLoginPasswordError             int = 1006 // 密码错误
	InsufficientBalance               int = 1007 // 余额不足
	UpdateContentNotExist             int = 1008 // 修改的内容不存在
	FindContentNotExist               int = 1009 // 查询的内容不存在
	ChannelDisable                    int = 1010 // 通道已禁用
	ChannelHaveLinkedChannelNotUpdate int = 1011 // 通道有关联关系存在，不能修改
	ChannelHaveLinkedChannelNotDelete int = 1012 // 通道有关联关系存在，不能删除
	NoLoginToken                      int = 1013 // 缺少登录token
	AccountDisable                    int = 1014 // 账号已被禁用
	UserLogoutFailed                  int = 1015 // 退出登录失败
	InvalidUpstreamChannel            int = 1016 // 无效的上游通道
	ChannelRepetition                 int = 1017 // 通道重复
	LoginTokenParseFailed             int = 1018 // 解析登录token失败
	ChannelNotExist                   int = 1019 // 通道不存在
	ChannelNotUpstreamChannel         int = 1020 // 通道未关联上游通道
	ConfigAlreadyExist                int = 1021 // 配置已经存在
	OrderNotExist                     int = 1022 // 订单不存在
	OrderAlreadyProcessed             int = 1023 // 订单已处理
	MerchantAlreadyDisable            int = 1024 // 商户已被禁用
	OrderNotOp                        int = 1025 // 订单不能执行此操作
	PlatformBankCardRepetition        int = 1026 // 收款卡重复
	OrderDispatch                     int = 1027 // 订单处理中，等待上游返回结果
	OrderDispatchSuccess              int = 1028 // 订单已经派单成功，不可重复派单
	OrderCompleted                    int = 1029 // 订单已完成
	UpstreamChannelUnusable           int = 1030 // 通道不可用
	MerchantNotExist                  int = 1031 // 商户不存在
	PayChannelIsInnerDeduction        int = 1032 // 代收通道只能是内扣
	PayOutChannelIsOutDeduction       int = 1033 // 代付通道只能是外扣
	AreaNotExist                      int = 1034 // 地区不存在
	MerchantAreaNotSame               int = 1035 // 商户地区与通道地区不一致
	OrderMissingNotifyUrl             int = 1036 // 订单缺少异步通知地址
	OrderNotifyFail                   int = 1037 // 订单通知失败
	NotData                           int = 1038 // 没有数据
	ExportFail                        int = 1039 // 导出失败
	ExportLimit31Day                  int = 1040 // 最长导出31天的数据
	CheckExportTime                   int = 1041 // 请选择导出时间
	EndTimeOverStartTime              int = 1042 // 结束时间早于开始时间
	AmountOver                        int = 1043 // 结束金额小于开始金额
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

	CodeMap[SystemUnknownErr] = "未知错误"
	CodeMap[SystemInternalErr] = "系统内部错误"
	CodeMap[SystemBusy] = "系统繁忙"

	CodeMap[SysDBGet] = "获取数据失败"
	CodeMap[SysDBAdd] = "添加数据失败"
	CodeMap[SysDBUpdate] = "更新数据失败"
	CodeMap[SysDBDelete] = "删除数据失败"
	CodeMap[SysDBErr] = "系统内部错误"
	CodeMap[SysDBSave] = "保存数据失败"

	CodeMap[GetLoginCaptchaFailed] = "获取验证码失败"
	CodeMap[LoginCaptchaNotMatch] = " 验证码错误"
	CodeMap[UserLoginFailed] = "用户登录失败"

	CodeMap[AccountRepeat] = "账号重复"
	CodeMap[UserNotExist] = "用户不存在"
	CodeMap[GetLoginPasswordError] = "密码错误"
	CodeMap[InsufficientBalance] = "余额不足"
	CodeMap[UpdateContentNotExist] = "修改的内容不存在"
	CodeMap[FindContentNotExist] = "查询的内容不存在"
	CodeMap[ChannelDisable] = "通道已禁用"
	CodeMap[ChannelHaveLinkedChannelNotUpdate] = "通道有关联关系存在，不能修改"
	CodeMap[ChannelHaveLinkedChannelNotDelete] = "通道有关联关系存在，不能删除"
	CodeMap[NoLoginToken] = "缺少登录token"
	CodeMap[UserNotLogin] = "用户未登录"
	CodeMap[UserLogoutFailed] = "退出登录失败"
	CodeMap[InvalidUpstreamChannel] = "无效的上游通道"
	CodeMap[ChannelRepetition] = "通道重复"
	CodeMap[LoginTokenParseFailed] = "解析登录token失败"
	CodeMap[ChannelNotExist] = "通道不存在"
	CodeMap[ChannelNotUpstreamChannel] = "当前通道未关联上游通道"
	CodeMap[ConfigAlreadyExist] = "配置已经存在"
	CodeMap[OrderNotExist] = "订单不存在"
	CodeMap[OrderAlreadyProcessed] = "订单已处理"
	CodeMap[MerchantAlreadyDisable] = "商户已被禁用"
	CodeMap[OrderNotOp] = "订单不能执行此操作"
	CodeMap[PlatformBankCardRepetition] = "收款卡重复"
	CodeMap[AccountDisable] = "账号已被禁用"
	CodeMap[OrderDispatch] = "订单处理中，等待上游返回结果"
	CodeMap[OrderDispatchSuccess] = "订单已经派单成功，不可重复派单"
	CodeMap[OrderCompleted] = "订单已完成"
	CodeMap[UpstreamChannelUnusable] = "通道不可用"
	CodeMap[MerchantNotExist] = "商户不存在"
	CodeMap[PayChannelIsInnerDeduction] = "代收通道只能是内扣"
	CodeMap[PayOutChannelIsOutDeduction] = "代付通道只能是外扣"
	CodeMap[AreaNotExist] = "地区不存在"
	CodeMap[MerchantAreaNotSame] = "商户地区与通道地区不一致"
	CodeMap[OrderMissingNotifyUrl] = "订单缺少异步通知地址"
	CodeMap[OrderNotifyFail] = "订单通知失败"
	CodeMap[NotData] = "没有数据"
	CodeMap[ExportFail] = "导出失败"
	CodeMap[ExportLimit31Day] = "最长导出31天的数据"
	CodeMap[CheckExportTime] = "请选择导出时间"
	CodeMap[EndTimeOverStartTime] = "结束时间早于开始时间"
	CodeMap[AmountOver] = "结束金额小于开始金额"
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
