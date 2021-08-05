package common

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"tpay_backend/utils"
)

// 创建商户账号密码
func CreateMerchantPassword(plainPassword string) string {
	//key1 := "Ytt6m$*9b"
	//key2 := "bQpd@k.3p8x"

	key1 := "Yo6m$*9b"
	key2 := "Qpd@k.3p8"
	return utils.Sha256(key1 + plainPassword + key2)
}

// 创建商户账号支付密码
func CreateMerchantPayPassword(plainPayPassword string) string {
	key1 := "Yay6m$*9b"
	key2 := "bQpd@k.3p8x"

	return utils.Sha256(key1 + plainPayPassword + key2)
}

//将密码解密为明文
func DecryptPassword(passWord string) (string, error) {
	plainPassword, err := utils.AesDecrypt(passWord, PasswordAecEncryptKey)
	if err != nil {
		return "", err
	}
	return plainPassword, nil
}

const (
	LoginTokenSignKey = "e10adc3949ba59abbe56e057f20f883e"
)

/**
返回前端的 x-login-token格式 --> user_id:token:sign
*/
func LoginTokenGenerate(userId int64, token string) string {
	sign := LoginTokenMakeSign(userId, token)
	return fmt.Sprintf("%v.%v.%v", token, userId, sign)
}

func LoginTokenParse(loginToken string) (int64, string, error) {
	// 格式: token.user_id.sign
	//tokenArr[0] token
	//tokenArr[1] user_id
	//tokenArr[2] sign

	tokenArr := strings.Split(loginToken, ".")
	if len(tokenArr) != 3 {
		return 0, "", errors.New("loginToken参数格式错误")
	}

	if strings.TrimSpace(tokenArr[0]) == "" {
		return 0, "", errors.New("loginToken 错误, token为空")
	}

	// 解析userid是否为int64
	userId, err := strconv.ParseInt(tokenArr[1], 10, 64)
	if err != nil {
		return 0, "", errors.New(fmt.Sprintf("loginToken 错误, userId转int64发生错误err[%v]", err))
	}

	if strings.TrimSpace(tokenArr[2]) == "" {
		return 0, "", errors.New("loginToken 错误, sign签名为空")
	}

	// 验证签名
	if LoginTokenMakeSign(userId, tokenArr[0]) != tokenArr[2] {
		return 0, "", errors.New("sign签名错误")
	}

	return userId, tokenArr[0], nil
}

func LoginTokenMakeSign(userId int64, token string) string {
	return utils.Md5(fmt.Sprintf("%s%s%d", token, LoginTokenSignKey, userId))
}
