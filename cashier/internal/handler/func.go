package handler

import (
	"fmt"
	"net/http"
	"tpay_backend/cashier/internal/lang"
	"tpay_backend/cashier/internal/svc"
	"tpay_backend/utils"

	"github.com/gin-gonic/gin"
)

func GetCurrentLang(c *gin.Context, svcCtx *svc.ServiceContext) (currentLang string, currentLangList []string) {
	// 1.程序内的默认语言
	currentLang = lang.DefaultLang

	//configModel := model.NewGlobalConfigModel(svcCtx.DbEngine)
	//
	//// 2.获取配置的语言列表
	//if list, err := configModel.CashierLangList(); err != nil {
	//	fmt.Printf("获取语言列表配置失败:%v\n", err)
	//} else {
	//	for _, v := range list {
	//		// 配置的语言必须在程序允许的范围内
	//		if utils.InSlice(v, lang.LangList) {
	//			currentLangList = append(currentLangList, v)
	//		}
	//	}
	//}
	//
	//// 3.获取配置的默认语言
	//if defaultLang, err := configModel.CashierDefaultLang(); err != nil {
	//	fmt.Printf("获取默认语言配置失败:%v\n", err)
	//} else {
	//	// 2.配置的默认语言不为空并且在允许的范围内,则为默认语言
	//	if defaultLang != "" && utils.InSlice(defaultLang, currentLangList) {
	//		currentLang = defaultLang
	//	}
	//}

	// 4.cookie中有语言
	if cookieLang, err := c.Cookie(LangCookieName); err != nil {
		if err == http.ErrNoCookie {
			//fmt.Printf("cookie不存在\n")
		} else {
			fmt.Printf("获取cooke失败:%v\n", err)
		}
	} else {
		//fmt.Printf("lang_cookie:%v\n", cookieLang)
		if utils.InSlice(cookieLang, currentLangList) {
			currentLang = cookieLang
		}
	}

	return
}
