package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/tal-tech/go-zero/core/logx"
	"net/http"
	"strings"
	"time"
	"tpay_backend/cashier/internal/svc"
	"tpay_backend/model"
	"tpay_backend/upstream"
	"tpay_backend/utils"
)

const (
	LangCookieName = "lang"
)

// 订单详情
func NewOrderDetailHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderNo := strings.TrimSpace(c.Query("o"))

		if orderNo == "" {
			c.String(http.StatusNotFound, default404Body)
			return
		}

		payOrder, err := model.NewPayOrderModel(svcCtx.DbEngine).FindOneByOrderNo(orderNo)
		if err != nil {
			fmt.Printf("查询订单失败:%v orderNo:%v \n", err, orderNo)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}
		merc, err := model.NewMerchantModel(svcCtx.DbEngine).FindOneByMerchantNo(payOrder.MerchantNo)
		if err != nil {
			fmt.Printf("查询商户失败:%v orderNo:%v \n", err, payOrder.MerchantNo)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}
		banks, err := model.NewBankModel(svcCtx.DbEngine).FindManyLimit(payOrder.ReqAmount)
		if err != nil {
			fmt.Printf("查询银行列表失败:%v orderNo:%v \n", err, payOrder.MerchantNo)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}

		// 获取当前语言
		currentLang, currentLangList := GetCurrentLang(c, svcCtx)
		fmt.Printf("currentLang:%v, currentLangList:%v\n", currentLang, currentLangList)

		// 金额是否需要除以100
		amount := float64(payOrder.ReqAmount)
		amount /= 100

		unixTimeI := utils.ToInt64(payOrder.CreateTime)
		unixTimeT := time.Unix(unixTimeI, 0).In(time.Local)
		dateTime := unixTimeT.Format("2006-01-02")

		bankM := map[string]string{}
		for _, v := range banks {
			bankM[v.BankName] = v.BankName
		}

		c.HTML(http.StatusOK, "1.html", gin.H{
			"amount":      amount,
			"merc_name":   merc.Username,
			"order_no":    payOrder.OrderNo,
			"create_time": dateTime,
			"banks":       bankM,
		})
	}
}

type SetCardReq struct {
	CardNo  string `json:"card_no"`
	OrderNo string `json:"order_no"`
}

func GetRedisKey(prefix string, t interface{}) string {
	s := fmt.Sprintf("%sQPaySign%s", prefix, t)
	s = strings.Replace(s, "*", "", -1) // 去掉符号 *
	s = strings.Replace(s, ".", "", -1) // 去掉符号 .
	return s
}

func GetBusinessrecordnumberRedisKey(t interface{}, mercNo string) string {
	s := fmt.Sprintf("businessrecordnumber:%s_%s", t, mercNo)
	s = strings.Replace(s, "*", "", -1) // 去掉符号 *
	s = strings.Replace(s, ".", "", -1) // 去掉符号 .
	return s
}

func GetCardNo2CertRedisKey(cardNo string, mercNo string) string {
	s := fmt.Sprintf("cert:%s_%s", cardNo, mercNo)
	s = strings.Replace(s, "*", "", -1) // 去掉符号 *
	s = strings.Replace(s, ".", "", -1) // 去掉符号 .
	return s
}

// 订单详情-测试订单
func NewSetCardHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := SetCardReq{}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			fmt.Printf("参数错误:%v orderNo:%v \n", err, req)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}
		logx.Infof("req=%v", req)

		payOrder, err := model.NewPayOrderModel(svcCtx.DbEngine).FindOneByOrderNo(req.OrderNo)
		if err != nil {
			fmt.Printf("查询订单失败:%v orderNo:%v \n", err, req.OrderNo)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}
		fmt.Println(payOrder)

		certcode, err := svcCtx.Redis.Get(context.TODO(), GetCardNo2CertRedisKey(req.CardNo, payOrder.MerchantNo)).Result()
		if err != nil {
			logx.Errorf("查询redis 出错,key=%v", GetRedisKey("", req.CardNo))
		}

		logx.Infof("certcode=%v\n", certcode)

		if certcode != "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"url":  "/3?o=" + req.OrderNo + "&c=" + req.CardNo + "&businessrecordnumber=" + certcode,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 0,
				"url":  "/2?o=" + req.OrderNo + "&c=" + req.CardNo + "&businessrecordnumber=" + "",
			})
		}

		// 3.选择上游
		//upstreamObj, err := getUpstream(payOrder.UpstreamChannelId, svcCtx)
		//if err != nil {
		//	logx.Errorf("代收-选择上游失败:err:%v,req:%+v,merchant:%+v", err, req, payOrder.MerchantNo)
		//	c.String(http.StatusInternalServerError, default405Body)
		//	return
		//}
		//
		//logx.Infof("代收-选择上游成功:%+v", upstreamObj)

		//resp, err := upstreamObj.(*upstream.TopPay).QPaySignQuery(&upstream.QPaySignQueryRequest{
		//	BankCardNo: req.CardNo,
		//})
		//if err != nil {
		//	logx.Errorf("err=%v", err)
		//	//return
		//}

		//logx.Infof("结果：%+v", resp)
		//if resp != nil && resp.IsSigned {
		//	// 缓存卡的签约信息
		//	//logx.Infof("redis=%v\n%v", svcCtx.Redis, svcCtx.Config.Name)
		//	respB, _ := json.Marshal(resp)
		//	ret, err := svcCtx.Redis.SetEX(context.TODO(), GetRedisKey(svcCtx.Config.Name, req.CardNo), respB, time.Minute*10).Result()
		//	logx.Infof("ret=%v,err=%v", ret, err)
		//	c.JSON(http.StatusOK, gin.H{
		//		"code": 0,
		//		"url":  "/3?o=" + req.OrderNo + "&c=" + req.CardNo + "&businessrecordnumber=" + ret,
		//	})
		//} else {
		//	c.JSON(http.StatusOK, gin.H{
		//		"code": 0,
		//		"url":  "/2?o=" + req.OrderNo + "&c=" + req.CardNo + "&businessrecordnumber=" + "",
		//	})
		//}

	}
}

type RegSmsReq struct {
	OrderNo    string `json:"order_no"`
	CardHolder string `json:"card_holder"`
	IdNo       string `json:"id_no"`
	Phone      string `json:"phone"`
	CardNo     string `json:"card_no"`
}

// 订单详情-测试订单
func RegSmsHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := RegSmsReq{}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			fmt.Printf("参数错误:%v orderNo:%v \n", err, req)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}
		logx.Infof("req=%v", req)

		req.IdNo = strings.TrimSpace(req.IdNo)
		req.Phone = strings.TrimSpace(req.Phone)
		req.CardHolder = strings.TrimSpace(req.CardHolder)
		req.CardNo = strings.TrimSpace(req.CardNo)
		if req.CardHolder == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 100,
				"msg":  "缺失持卡人姓名",
			})
			return
		}

		if req.IdNo == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 100,
				"msg":  "缺失证件号",
			})
			return
		}

		if req.Phone == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 100,
				"msg":  "缺失电话",
			})
			return
		}

		if req.CardNo == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 100,
				"msg":  "缺失卡号",
			})
			return
		}

		if req.OrderNo == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 100,
				"msg":  "订单号不全",
			})
			return
		}

		payOrder, err := model.NewPayOrderModel(svcCtx.DbEngine).FindOneByOrderNo(req.OrderNo)
		if err != nil {
			fmt.Printf("查询订单失败:%v orderNo:%v \n", err, req.OrderNo)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}
		fmt.Println(payOrder)

		// 3.选择上游
		upstreamObj, err := getUpstream(payOrder.UpstreamChannelId, svcCtx)
		if err != nil {
			logx.Errorf("代收-选择上游失败:err:%v,req:%+v,merchant:%+v", err, req, payOrder.MerchantNo)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}

		logx.Infof("upstreamConfig=%v", upstreamObj)
		upstreamConfig := upstreamObj.GetUpstreamConfig()
		upNotifyUrl, err := getUpstreamNotifyUrl(upstreamConfig.PayNotifyPath, svcCtx)
		if err != nil {
			logx.Errorf("获取上游异步回调地址失败,payNotifyPath:%v,err=%v", upstreamConfig.PayNotifyPath, err)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}

		logx.Infof("代收-选择上游成功:%+v", upstreamObj)
		resp, err := upstreamObj.(*upstream.TopPay).QPaySignSms(&upstream.QPaySignSmsRequest{
			Amount:    payOrder.ReqAmount,
			Currency:  payOrder.Currency,
			OrderNo:   payOrder.OrderNo,
			NotifyUrl: upNotifyUrl,
			ReturnUrl: payOrder.ReturnUrl,
			BankCardInfo: upstream.BankCardInfo{
				BankCardNo:    req.CardNo,
				BankCardName:  req.CardHolder,
				BankCardIdNo:  req.IdNo,
				BankCardPhone: req.Phone,
				IsCreditCard:  false,
			},
		})
		if err != nil {
			logx.Errorf("err=%v", err)
			if resp == nil {
				c.JSON(http.StatusOK, gin.H{
					"code": 100,
					"msg":  "签约失败",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code": 100,
					"msg":  resp.ErrMsg,
				})
			}
			return
		}
		logx.Infof("结果：%+v", resp)

		//存储签约订单号到redis中去
		ret, err := svcCtx.Redis.SetEX(context.TODO(), GetBusinessrecordnumberRedisKey(resp.Businessrecordnumber, payOrder.MerchantNo), req.CardNo, time.Hour*24*7).Result()
		logx.Infof("Businessrecordnumber|ret=%v, err=%v", ret, err)
		ret2, err := svcCtx.Redis.SetEX(context.TODO(), GetCardNo2CertRedisKey(req.CardNo, payOrder.MerchantNo), resp.Certcode, time.Hour*24*7).Result()
		logx.Infof("Certcode|ret2=%v, err=%v", ret2, err)

		c.JSON(http.StatusOK, gin.H{
			"code":                 0,
			"businessrecordnumber": resp.Businessrecordnumber,

			//如果有签约的话，以下是已签约的话就跳转到的页面
			"url":      "/3?o=" + req.OrderNo + "&c=" + req.CardNo,
			"certcode": resp.Certcode,
		})
	}
}

// 获取上游异步通知地址
func getUpstreamNotifyUrl(notifyPath string, svcCtx *svc.ServiceContext) (string, error) {
	// 获取payapi站点域名
	host, err := model.NewGlobalConfigModel(svcCtx.DbEngine).FindValueByKey(model.ConfigPayapiHostAddr)
	if err != nil {
		return "", errors.New(fmt.Sprintf("查询网站配置失败,key:%v,err=%v", model.ConfigPayapiHostAddr, err))
	}

	if strings.TrimSpace(host) == "" {
		return "", errors.New("系统没有配置异步回调地址")
	}

	if strings.TrimSpace(notifyPath) == "" {
		return "", errors.New("系统没有配置异步回调地址路径")
	}

	notifyPath = strings.TrimPrefix(notifyPath, "/")
	notifyUrl := strings.TrimRight(host, "/") + "/" + strings.TrimRight(notifyPath, "/")

	return notifyUrl, nil
}

type RegSmsConfirmReq struct {
	OrderNo              string `json:"order_no"`
	Businessrecordnumber string `json:"businessrecordnumber"`
	CardNo               string `json:"card_no"`
	VerifyCode           string `json:"verify_code"`
}

// 签约确认
func RegSmsConfirmHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := RegSmsConfirmReq{}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			fmt.Printf("参数错误:%v orderNo:%v \n", err, req)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}
		logx.Infof("req=%v", req)

		req.VerifyCode = strings.TrimSpace(req.VerifyCode)
		req.Businessrecordnumber = strings.TrimSpace(req.Businessrecordnumber)
		req.OrderNo = strings.TrimSpace(req.OrderNo)
		req.CardNo = strings.TrimSpace(req.CardNo)
		if req.VerifyCode == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 100,
				"msg":  "缺失验证码",
			})
			return
		}

		if req.CardNo == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 100,
				"msg":  "缺失卡号",
			})
			return
		}

		if req.OrderNo == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 100,
				"msg":  "订单号不全",
			})
			return
		}
		if req.Businessrecordnumber == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 100,
				"msg":  "签约订单号不全",
			})
			return
		}

		payOrder, err := model.NewPayOrderModel(svcCtx.DbEngine).FindOneByOrderNo(req.OrderNo)
		if err != nil {
			fmt.Printf("查询订单失败:%v orderNo:%v \n", err, req.OrderNo)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}
		//fmt.Println(payOrder)

		// 3.选择上游
		upstreamObj, err := getUpstream(payOrder.UpstreamChannelId, svcCtx)
		if err != nil {
			logx.Errorf("代收-选择上游失败:err:%v,req:%+v,merchant:%+v", err, req, payOrder.MerchantNo)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}

		//确认签约
		logx.Infof("代收-选择上游成功:%+v", upstreamObj)
		resp, err := upstreamObj.(*upstream.TopPay).QPaySignConfirm(&upstream.QPaySignConfirmRequest{
			Businessrecordnumber: req.Businessrecordnumber,
			Verifycode:           req.VerifyCode,
		})
		if err != nil {
			logx.Errorf("err=%v", err)
			if resp == nil {
				c.JSON(http.StatusOK, gin.H{
					"code": 100,
					"msg":  "签约失败",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code": 100,
					"msg":  resp.ErrMsg,
				})
			}
			return
		}
		logx.Infof("结果：%+v", resp)

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"url":  "/3?o=" + req.OrderNo + "&c=" + req.CardNo,
		})
	}
}

type QPaySmsConfirmHandlerReq struct {
	OrderNo    string `json:"order_no"`
	CardNo     string `json:"card_no"`
	VerifyCode string `json:"verify_code"`
}

// 签约确认
func QPaySmsConfirmHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {

		req := RegSmsConfirmReq{}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			fmt.Printf("参数错误:%v orderNo:%v \n", err, req)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}

		logx.Infof("req=%v", req)

		req.VerifyCode = strings.TrimSpace(req.VerifyCode)
		req.OrderNo = strings.TrimSpace(req.OrderNo)
		req.CardNo = strings.TrimSpace(req.CardNo)
		req.Businessrecordnumber = strings.TrimSpace(req.Businessrecordnumber)
		if req.VerifyCode == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 100,
				"msg":  "缺失验证码",
			})
			return
		}

		if req.CardNo == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 100,
				"msg":  "缺失卡号",
			})
			return
		}

		if req.OrderNo == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 100,
				"msg":  "订单号不全",
			})
			return
		}

		if req.Businessrecordnumber == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 100,
				"msg":  "签约订单号不全",
			})
			return
		}

		payOrder, err := model.NewPayOrderModel(svcCtx.DbEngine).FindOneByOrderNo(req.OrderNo)
		if err != nil {
			fmt.Printf("查询订单失败:%v orderNo:%v \n", err, req.OrderNo)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}
		fmt.Println(payOrder)

		// 3.选择上游
		upstreamObj, err := getUpstream(payOrder.UpstreamChannelId, svcCtx)
		if err != nil {
			logx.Errorf("代收-选择上游失败:err:%v,req:%+v,merchant:%+v", err, req, payOrder.MerchantNo)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}

		logx.Infof("代收-选择上游成功:%+v", upstreamObj)

		//5.确认上游订单存在才能进行支付
		resp, err := upstreamObj.(*upstream.TopPay).QPayConfirm(&upstream.QPayConfirmRequest{
			Businessrecordnumber: req.Businessrecordnumber,
			Verifycode:           req.VerifyCode,
		})
		if err != nil {
			logx.Errorf("err=%v", err)
			if resp == nil {
				c.JSON(http.StatusOK, gin.H{
					"code": 100,
					"msg":  "支付失败",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code": 100,
					"msg":  resp.ErrMsg,
				})
			}
			return
		}
		logx.Infof("结果：%+v", resp)

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"url":  "/4?o=" + req.OrderNo + "&c=" + req.CardNo,
		})
	}
}

// 订单详情-测试订单
func SmsSignHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderNo := strings.TrimSpace(c.Query("o"))

		if orderNo == "" {
			c.String(http.StatusNotFound, default404Body)
			return
		}

		payOrder, err := model.NewPayOrderModel(svcCtx.DbEngine).FindOneByOrderNo(orderNo)
		if err != nil {
			fmt.Printf("查询订单失败:%v orderNo:%v \n", err, orderNo)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}
		merc, err := model.NewMerchantModel(svcCtx.DbEngine).FindOneByMerchantNo(payOrder.MerchantNo)
		if err != nil {
			fmt.Printf("查询商户失败:%v orderNo:%v \n", err, payOrder.MerchantNo)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}

		// 获取当前语言
		currentLang, currentLangList := GetCurrentLang(c, svcCtx)
		fmt.Printf("currentLang:%v, currentLangList:%v\n", currentLang, currentLangList)

		// 金额是否需要除以100
		amount := float64(payOrder.ReqAmount)
		amount /= 100

		unixTimeI := utils.ToInt64(payOrder.CreateTime)
		unixTimeT := time.Unix(unixTimeI, 0).In(time.Local)
		dateTime := unixTimeT.Format("2006-01-02")

		c.HTML(http.StatusOK, "2.html", gin.H{
			"amount":      amount,
			"merc_name":   merc.Username,
			"order_no":    payOrder.OrderNo,
			"create_time": dateTime,
		})
	}
}

// 最后付款页面
func FinalPayHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderNo := strings.TrimSpace(c.Query("o"))

		if orderNo == "" {
			c.String(http.StatusNotFound, default404Body)
			return
		}

		payOrder, err := model.NewPayOrderModel(svcCtx.DbEngine).FindOneByOrderNo(orderNo)
		if err != nil {
			fmt.Printf("查询订单失败:%v orderNo:%v \n", err, orderNo)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}
		merc, err := model.NewMerchantModel(svcCtx.DbEngine).FindOneByMerchantNo(payOrder.MerchantNo)
		if err != nil {
			fmt.Printf("查询商户失败:%v orderNo:%v \n", err, payOrder.MerchantNo)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}

		// 获取当前语言
		currentLang, currentLangList := GetCurrentLang(c, svcCtx)
		fmt.Printf("currentLang:%v, currentLangList:%v\n", currentLang, currentLangList)

		// 金额是否需要除以100
		amount := float64(payOrder.ReqAmount)
		amount /= 100

		unixTimeI := utils.ToInt64(payOrder.CreateTime)
		unixTimeT := time.Unix(unixTimeI, 0).In(time.Local)
		dateTime := unixTimeT.Format("2006-01-02")

		c.HTML(http.StatusOK, "3.html", gin.H{
			"amount":      amount,
			"merc_name":   merc.Username,
			"order_no":    payOrder.OrderNo,
			"create_time": dateTime,
		})
	}
}

// 最后付款页面
func FinishHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		orderNo := strings.TrimSpace(c.Query("o"))

		if orderNo == "" {
			c.String(http.StatusNotFound, default404Body)
			return
		}

		payOrder, err := model.NewPayOrderModel(svcCtx.DbEngine).FindOneByOrderNo(orderNo)
		if err != nil {
			fmt.Printf("查询订单失败:%v orderNo:%v \n", err, orderNo)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}
		merc, err := model.NewMerchantModel(svcCtx.DbEngine).FindOneByMerchantNo(payOrder.MerchantNo)
		if err != nil {
			fmt.Printf("查询商户失败:%v orderNo:%v \n", err, payOrder.MerchantNo)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}

		// 获取当前语言
		currentLang, currentLangList := GetCurrentLang(c, svcCtx)
		fmt.Printf("currentLang:%v, currentLangList:%v\n", currentLang, currentLangList)

		// 金额是否需要除以100
		amount := float64(payOrder.ReqAmount)
		amount /= 100

		unixTimeI := utils.ToInt64(payOrder.CreateTime)
		unixTimeT := time.Unix(unixTimeI, 0).In(time.Local)
		dateTime := unixTimeT.Format("2006-01-02")

		c.HTML(http.StatusOK, "4.html", gin.H{
			"amount":      amount,
			"merc_name":   merc.Username,
			"order_no":    payOrder.OrderNo,
			"create_time": dateTime,
		})
	}
}

type QPaySmsReq struct {
	OrderNo    string `json:"order_no"`
	CardHolder string `json:"card_holder"`
	IdNo       string `json:"id_no"`
	Phone      string `json:"phone"`
	CardNo     string `json:"card_no"`
}

// 订单详情-测试订单
func QPaySmsHandler(svcCtx *svc.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := QPaySmsReq{}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			fmt.Printf("参数错误:%v orderNo:%v \n", err, req)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}
		logx.Infof("req=%v", req)

		req.Phone = strings.TrimSpace(req.Phone)
		req.OrderNo = strings.TrimSpace(req.OrderNo)
		req.CardNo = strings.TrimSpace(req.CardNo)
		if req.Phone == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 100,
				"msg":  "缺失电话",
			})
			return
		}

		if req.CardNo == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 100,
				"msg":  "缺失卡号",
			})
			return
		}

		if req.OrderNo == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 100,
				"msg":  "订单号不全",
			})
			return
		}

		payOrder, err := model.NewPayOrderModel(svcCtx.DbEngine).FindOneByOrderNo(req.OrderNo)
		if err != nil {
			fmt.Printf("查询订单失败:%v orderNo:%v \n", err, req.OrderNo)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}
		fmt.Println(payOrder)

		// 3.选择上游
		upstreamObj, err := getUpstream(payOrder.UpstreamChannelId, svcCtx)
		if err != nil {
			logx.Errorf("代收-选择上游失败:err:%v,req:%+v,merchant:%+v", err, req, payOrder.MerchantNo)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}

		logx.Infof("upstreamConfig=%v", upstreamObj)
		upstreamConfig := upstreamObj.GetUpstreamConfig()
		upNotifyUrl, err := getUpstreamNotifyUrl(upstreamConfig.PayNotifyPath, svcCtx)
		if err != nil {
			logx.Errorf("获取上游异步回调地址失败,payNotifyPath:%v,err=%v", upstreamConfig.PayNotifyPath, err)
			c.String(http.StatusInternalServerError, default405Body)
			return
		}

		fmt.Printf("-----------------------------GetRedisKey = %v--------------------------",
			GetRedisKey("", req.CardNo))

		certcode, err := svcCtx.Redis.Get(context.TODO(), GetCardNo2CertRedisKey(req.CardNo, payOrder.MerchantNo)).Result()
		if err != nil {
			logx.Errorf("查询redis 出错,key=%v", GetRedisKey("", req.CardNo))
		}

		logx.Infof("certcode=%v\n", certcode)

		currency, err := model.NewCurrencyModel(svcCtx.DbEngine).FindByCurrency(payOrder.Currency)
		if err != nil {
			logx.Errorf("币种有误, err:%v", err)
			return
		}
		billamount := fmt.Sprintf("%d", payOrder.ReqAmount)

		if currency.IsDivideHundred == model.DivideHundred {

			billamountDe, err := decimal.NewFromString(billamount)
			if err != nil {
				logx.Errorf("金额有误, err:%v", err)
				return
			}

			billamount = billamountDe.Div(decimal.NewFromInt(100)).String()
		}

		logx.Infof("代收-选择上游成功:%+v", upstreamObj)
		resp, err := upstreamObj.(*upstream.TopPay).QPay(&upstream.QPayRequest{
			Backurl:        upNotifyUrl,
			Subject:        "普通客户维护收费",
			Businesstype:   "其他商家消费",
			Kind:           "电子小票",
			Description:    "普通客户维护收费",
			Businessnumber: req.OrderNo,
			Billamount:     billamount,
			//Toaccountnumber: "8800013385833", //收款方会员账号
			Certcode: certcode,
			Clientip: utils.GetFakeIp(),
			//Merchantuserid:  utils.ToStringNoPoint(utils.Random(10000, 99999999)),//在商户平台中用户id
		})
		if err != nil {
			logx.Errorf("err=%v", err)
			if resp == nil {
				c.JSON(http.StatusOK, gin.H{
					"code": 100,
					"msg":  "签约失败",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"code": 100,
					"msg":  resp.ErrMsg,
				})
			}
			return
		}
		logx.Infof("结果：%+v", resp)

		c.JSON(http.StatusOK, gin.H{
			"code":                 0,
			"businessrecordnumber": resp.Businessrecordnumber,
		})
	}
}

func getTopPayObj() upstream.Upstream {
	jsonStr := `{
		"Host": "http://127.0.0.1:8080"
	}`
	//jsonStr := `{
	//	"Host": "http://127.0.0.1:23331"
	//}`
	obj, err := upstream.NewTopPay("", jsonStr)
	if err != nil {
		logx.Errorf("获取对象失败:%v", err)
	}

	return obj
}

// 获取上游
func getUpstream(upstreamChannelId int64, ctx *svc.ServiceContext) (upstream.Upstream, error) {
	var err error
	// 1.获取对应上游
	up, err := model.NewUpstreamChannelModel(ctx.DbEngine).FindUpstreamByChannelId(upstreamChannelId)
	if err != nil {
		if err == model.ErrRecordNotFound {
			return nil, errors.New(fmt.Sprintf("未找到对应的上游:upstreamChannelId:%v", upstreamChannelId))
		} else {
			return nil, errors.New(fmt.Sprintf("查询上游信息失败:err:%v,upstreamChannelId:%v", err, upstreamChannelId))
		}
	}

	obj, err := getUpstreamObject(up)

	return obj, err
}

func getUpstreamObject(up *model.Upstream) (upstream.Upstream, error) {
	var upstreamObj upstream.Upstream
	var err error

	// 2.初始化上游
	switch up.UpstreamCode {
	case upstream.UpstreamCodeTotopay: // totopay
		upstreamObj, err = upstream.NewTotopay(up.UpstreamMerchantNo, up.CallConfig)
	case upstream.UpstreamCodeGoldPays:
		upstreamObj, err = upstream.NewGoldPays(up.UpstreamMerchantNo, up.CallConfig)
	case upstream.UpstreamCodeZf777Pay:
		upstreamObj, err = upstream.NewThreeSevenPay(up.UpstreamMerchantNo, up.CallConfig)
	case upstream.UpstreamCodeXPay:
		upstreamObj, err = upstream.NewXPay(up.UpstreamMerchantNo, up.CallConfig)
	case upstream.UpstreamCodeToppay:
		upstreamObj, err = upstream.NewTopPay(up.UpstreamMerchantNo, up.CallConfig)
	}

	if err != nil {
		return nil, errors.New(fmt.Sprintf("初始化上游失败err:%v,upstream:%v,config:%v", err, up.UpstreamName, up.CallConfig))
	}

	if upstreamObj == nil {
		return nil, errors.New(fmt.Sprintf("上游(%v)未配置", up.UpstreamName))
	}

	return upstreamObj, nil
}
