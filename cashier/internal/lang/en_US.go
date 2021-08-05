package lang

var LanguageEnUs = make(map[string]string)

func initLangEnUs() {
	LanguageEnUs["订单"] = "Order Number"
	LanguageEnUs["金额"] = "Payment amount"
	LanguageEnUs["银行"] = "Bank of Account"
	LanguageEnUs["卡号"] = "Bank Account Number"
	LanguageEnUs["姓名"] = "Account Name"

	LanguageEnUs["A001"] = "After the payment is received, we will confirm the receipt and delivery within 5 minutes"
	LanguageEnUs["剩余支付时间"] = "Remaining payment time"
	LanguageEnUs["复制"] = "Copy"
	LanguageEnUs["支行"] = "Branch"
	LanguageEnUs["复制成功"] = "Copy successfully"
	LanguageEnUs["请付款至以下账户"] = "Please pay to the following account"
	LanguageEnUs["收银台"] = "Cashier"
	LanguageEnUs["我已完成付款"] = "Payment completed"

	LanguageEnUs["请输入付款人姓名"] = "请输入付款人姓名"
	LanguageEnUs["下一步"] = "下一步"
	LanguageEnUs["缺少订单号"] = "缺少订单号"
	LanguageEnUs["查询订单失败"] = "查询订单失败"
	LanguageEnUs["付款人姓名已经存在"] = "付款人姓名已经存在"
	LanguageEnUs["更新付款人姓名失败"] = "更新付款人姓名失败"
	LanguageEnUs["提示：不按照正常支付，无法及时到账，修改金额支付，不到账概不负责"] = "Tips: Do not follow the normal payment, can not arrive in time, modify the amount of payment, not responsible for the account"

}
