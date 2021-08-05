package handler

import (
	"tpay_backend/cashier/internal/svc"

	"github.com/gin-gonic/gin"
)

var default404Body = "404 page not found"
var default405Body = "405 method not allowed"

const (
	NameOrderInputPath  = "/name"
	NameOrderDetailPath = "/name/order"
)

func RegisterHandlers(engine *gin.Engine, svcCtx *svc.ServiceContext) {
	// 订单信息
	engine.GET("/1", NewOrderDetailHandler(svcCtx))
	// 订单信息-测试单
	engine.GET("/2", SmsSignHandler(svcCtx))
	// 订单信息-测试单
	engine.GET("/3", FinalPayHandler(svcCtx))
	// 订单信息-测试单
	engine.GET("/4", FinishHandler(svcCtx))
	// 订单信息-测试单
	engine.POST("/card", NewSetCardHandler(svcCtx))
	// 绑卡短信
	engine.POST("/reg_sms", RegSmsHandler(svcCtx))
	// 签约确认
	engine.POST("/reg_sms_confirm", RegSmsConfirmHandler(svcCtx))
	// 绑卡短信
	engine.POST("/pay_sms", QPaySmsHandler(svcCtx))
	// 支付确认
	engine.POST("/pay_sms_confirm", QPaySmsConfirmHandler(svcCtx))

	// 需要付款人-输入付款人页面
	engine.GET(NameOrderInputPath, NewNameOrderInputHandler(svcCtx))

	// 需要付款人-输入付款人页面-提交付款人
	engine.POST(NameOrderInputPath, NewNameOrderInputDoHandler(svcCtx))

	// 需要付款人-订单详情页面
	engine.GET(NameOrderDetailPath, NewNameOrderDetailHandler(svcCtx))
}
