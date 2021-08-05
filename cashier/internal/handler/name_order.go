package handler

import (
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"net/http"
	"strings"
	//"time"
	"tpay_backend/cashier/internal/lang"
	"tpay_backend/cashier/internal/svc"
	//"tpay_backend/model"

	"github.com/gin-gonic/gin"
)

// 需要付款人-输入付款人页面
func NewNameOrderInputHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderNo := strings.TrimSpace(c.Query("o"))

		if orderNo == "" {
			c.String(http.StatusNotFound, default404Body)
			return
		}

		//var payOrder *model.PayOrder
		//var err error

		//if orderNo != "" {
		//	payOrder, err = model.NewPayOrderModel(svcCtx.DbEngine).FindOneByOrderNo(orderNo)
		//	if err != nil {
		//		fmt.Printf("查询订单失败:%v orderNo:%v \n", err, orderNo)
		//		c.String(http.StatusInternalServerError, default405Body)
		//		return
		//	}
		//}
		//fmt.Printf("payOrder:%+v\n", payOrder)

		// 已经有付款人姓名了
		//if strings.TrimSpace(payOrder.PayerName) != "" {
		//	c.Redirect(http.StatusFound, GetNameOrderDetailPath(orderNo))
		//	return
		//}

		// 获取当前语言
		currentLang, currentLangList := GetCurrentLang(c, svcCtx)
		fmt.Printf("currentLang:%v, currentLangList:%v\n", currentLang, currentLangList)

		//// 币种信息
		//currencyConfig, err := model.NewGlobalConfigModel(svcCtx.DbEngine).SystemCurrencyConfig()
		//if err != nil {
		//	fmt.Printf("查询币种信息失败:%v \n", err)
		//}
		//
		//// 金额是否需要除以100
		//reqAmount := float64(payOrder.ReqAmount)
		//if currencyConfig.IsDivideHundred == model.CurrencyDivideHundred {
		//	reqAmount /= 100
		//}

		c.HTML(http.StatusOK, "name_order_input.html", gin.H{
			//"title":     "输入付款人姓名",
			//"payOrder":  payOrder,
			//"maxTime":   payOrder.ExpiredTime - time.Now().Unix(),
			//"reqAmount": reqAmount,
			//
			//"langCookieName":        LangCookieName,
			//"cLang":                 currentLang,
			//"cLangList":             currentLangList,
			//"currencySymbol":        currencyConfig.Symbol,
			//"name_order_input_path": NameOrderInputPath,
			//"success_jump_url":      GetNameOrderDetailPath(payOrder.OrderNo),
		})
	}
}

// 需要付款人-输入付款人页面-提交付款人
func NewNameOrderInputDoHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderNo := strings.TrimSpace(c.PostForm("order_no"))
		payerName := strings.TrimSpace(c.PostForm("payer_name"))
		logx.Infof("req|order=[%v],name=[%v]", orderNo, payerName)

		// 获取当前语言
		currentLang, currentLangList := GetCurrentLang(c, svcCtx)
		fmt.Printf("currentLang:%v, currentLangList:%v\n", currentLang, currentLangList)

		var res struct {
			Code int64  `json:"code"`
			Msg  string `json:"msg"`
		}
		res.Code = 1 // 默认失败

		if orderNo == "" {
			res.Msg = lang.Lang(currentLang, "缺少订单号")
			c.JSON(http.StatusOK, res)
			return
		}

		if payerName == "" {
			res.Msg = lang.Lang(currentLang, "请输入付款人姓名")
			c.JSON(http.StatusOK, res)
			return
		}

		//var payOrder *model.PayOrder
		//var err error

		//payOrder, err = model.NewPayOrderModel(svcCtx.DbEngine).FindOneByOrderNo(orderNo)
		//if err != nil {
		//	fmt.Printf("查询订单失败:%v orderNo:%v \n", err, orderNo)
		//
		//	res.Msg = lang.Lang(currentLang, "查询订单失败")
		//	c.JSON(http.StatusOK, res)
		//	return
		//}

		//fmt.Printf("payOrder:%+v\n", payOrder)

		//// 已经有付款人姓名了
		//if strings.TrimSpace(payOrder.PayerName) != "" {
		//	res.Msg = lang.Lang(currentLang, "付款人姓名已经存在")
		//	c.JSON(http.StatusOK, res)
		//	return
		//}
		//
		//// 更新付款人姓名
		//logx.Infof("db|order=[%v],name=[%v]", payOrder.Id, payerName)
		//err = model.NewPayOrderModel(svcCtx.DbEngine).UpdatePayerName(payOrder.Id, payerName)
		//if err != nil {
		//	fmt.Printf("更新付款人姓名失败:%v payOrderId:%v, payerName:%v \n", err, payOrder.Id, payerName)
		//
		//	res.Msg = lang.Lang(currentLang, "更新付款人姓名失败")
		//	c.JSON(http.StatusOK, res)
		//	return
		//}

		res.Code = 0
		res.Msg = "success"
		c.JSON(http.StatusOK, res)
		return
	}
}

// 需要付款人-订单详情页面
func NewNameOrderDetailHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderNo := strings.TrimSpace(c.Query("o"))

		if orderNo == "" {
			c.String(http.StatusNotFound, default404Body)
			return
		}

		//payOrder, err := model.NewPayOrderModel(svcCtx.DbEngine).FindOneByOrderNo(orderNo)
		//if err != nil {
		//	fmt.Printf("查询订单失败:%v orderNo:%v \n", err, orderNo)
		//	c.String(http.StatusInternalServerError, default405Body)
		//	return
		//}

		//// 没有付款人姓名
		//if strings.TrimSpace(payOrder.PayerName) == "" {
		//	c.Redirect(http.StatusFound, GetNameOrderInputPath(orderNo))
		//	return
		//}
		//
		//// 收款卡信息改为从订单信息获取
		//// 为了兼容如果订单信息中收款卡号为空就还是要查一次，这种数据的收款卡信息可能不一致
		//carddealerCard := new(model.CarddealerCard)
		//if payOrder.CarddealerCardNumber != "" && payOrder.CarddealerBankName != "" {
		//	carddealerCard.CardHolder = payOrder.CarddealerCardHolder
		//	carddealerCard.CardNumber = payOrder.CarddealerCardNumber
		//	carddealerCard.BankName = payOrder.CarddealerBankName
		//	carddealerCard.CardOrganization = payOrder.CarddealerCardOrganization
		//} else {
		//	if payOrder.CardId > 0 {
		//		carddealerCard, err = model.NewCarddealerCardModel(svcCtx.DbEngine).FindOneById(payOrder.CardId)
		//		if err != nil {
		//			fmt.Printf("查询卡信息失败:%v CardId:%v \n", err, payOrder.CardId)
		//			c.String(http.StatusInternalServerError, default405Body)
		//			return
		//		}
		//	}
		//}

		//// 获取当前语言
		//currentLang, currentLangList := GetCurrentLang(c, svcCtx)
		//fmt.Printf("currentLang:%v, currentLangList:%v\n", currentLang, currentLangList)
		//
		//// 币种信息
		//currencyConfig, err := model.NewGlobalConfigModel(svcCtx.DbEngine).SystemCurrencyConfig()
		//if err != nil {
		//	fmt.Printf("查询币种信息失败:%v \n", err)
		//}
		//
		//// 金额是否需要除以100
		//reqAmount := float64(payOrder.ReqAmount)
		//if currencyConfig.IsDivideHundred == model.CurrencyDivideHundred {
		//	reqAmount /= 100
		//}

		c.HTML(http.StatusOK, "name_order_detail.html", gin.H{
			//"title":     "Order Info",
			//"payOrder":  payOrder,
			//"card":      carddealerCard,
			//"maxTime":   payOrder.ExpiredTime - time.Now().Unix(),
			//"reqAmount": reqAmount,
			//
			//"langCookieName": LangCookieName,
			//"cLang":          currentLang,
			//"cLangList":      currentLangList,
			//"currencySymbol": currencyConfig.Symbol,
		})
	}
}

func GetNameOrderInputPath(orderNo string) string {
	return fmt.Sprintf("%s?o=%s", NameOrderInputPath, orderNo)
}

func GetNameOrderDetailPath(orderNo string) string {
	return fmt.Sprintf("%s?o=%s", NameOrderDetailPath, orderNo)
}
