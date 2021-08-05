package lang

var LanguageThTH = make(map[string]string)

func initLangThTH() {
	LanguageThTH["订单"] = "หมายเลขคำสั่งซื้อ"
	LanguageThTH["金额"] = "จำนวนเงินที่ชำระ"
	LanguageThTH["银行"] = "เปิดบัญชีธนาคาร"
	LanguageThTH["卡号"] = "บัญชีธนาคาร"
	LanguageThTH["姓名"] = "ชื่อบัญชี"

	LanguageThTH["A001"] = "หลังจากการชำระเงินมาถึงเราจะยืนยันการรับภายใน 5 นาที"
	LanguageThTH["剩余支付时间"] = "เวลาชำระเงินที่เหลือ"
	LanguageThTH["复制"] = "คัดลอก"
	LanguageThTH["支行"] = "支行"
	LanguageThTH["复制成功"] = "คัดลอกสำเร็จ"
	LanguageThTH["请付款至以下账户"] = "请付款至以下账户"
	LanguageThTH["收银台"] = "收银台"
	LanguageThTH["我已完成付款"] = "我已完成付款"

	LanguageThTH["请输入付款人姓名"] = "请输入付款人姓名"
	LanguageThTH["下一步"] = "下一步"
	LanguageThTH["缺少订单号"] = "缺少订单号"
	LanguageThTH["查询订单失败"] = "查询订单失败"
	LanguageThTH["付款人姓名已经存在"] = "付款人姓名已经存在"
	LanguageThTH["更新付款人姓名失败"] = "更新付款人姓名失败"
	LanguageThTH["提示：不按照正常支付，无法及时到账，修改金额支付，不到账概不负责"] = "提示：不按照正常支付，无法及时到账，修改金额支付，不到账概不负责"
}
