package common

import "time"

const (
	LoginRedisKeyPrefix = "adminapi:login:" // 登录存入redis的键名前缀
	LoginRedisExpire    = time.Hour * 2     // redis登录有效期
	LoginTokenHeaderKey = "x-login-token"   // 登录头部token

	RequestHeaderLoginedUserId = "login-user-id"
	RequestHeaderLoginedToken  = "login-token-id"

	MerchantLoginRedisKeyPrefix = "merchantapi:login:" // 商家登录存入redis的键名前缀

	PasswordAecEncryptKey = "U2Fs-GVkX/9W+GVrAxcL2in1"
)

const (
	AdminDefaultPassword       = "123456"  // 管理员默认密码
	MerchantDefaultPassword    = "q123456" // 商户默认密码
	MerchantDefaultPayPassword = "123456"  // 商户默认支付密码
)

const (
	CaptchaPrefix = "captcha:adminapi:" // redis键前缀
	CaptchaExpire = time.Second * 600   // 10分钟
)
