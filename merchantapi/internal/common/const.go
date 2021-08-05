package common

import "time"

const (
	JwtUserId = "user_id"
)

const (
	LoginRedisKeyPrefix = "merchantapi:login:" // 登录存入redis的键名前缀
	LoginRedisExpire    = time.Hour * 2        // redis登录有效期
	LoginTokenHeaderKey = "x-login-token"      // 登录头部token

	RequestHeaderLoginedUserId = "login-user-id"
	RequestHeaderLoginedToken  = "login-token-id"

	PasswordAecEncryptKey = "BdPEHs10Ku1g7rMijlGvAQ+="
)

const (
	CaptchaPrefix = "captcha:merchantapi:" // redis键前缀
	CaptchaExpire = time.Second * 600      // 10分钟
)
