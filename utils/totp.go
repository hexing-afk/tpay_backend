package utils

import (
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// 生成TOTP秘钥
func GenerateTOTPSecret(accountName string) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "tpay",
		AccountName: accountName,
		Period:      30,
		Algorithm:   otp.AlgorithmSHA512,
	})

	if err != nil {
		return "", err
	}

	// key.Url格式:
	// otpauth://totp/Example.com:alice@example.com?algorithm=SHA512&digits=6&issuer=Example.com&period=30&secret=XA4HKCWD5AQOQR5HK6INZ3NM36LIGA6C

	// key.Secret格式: XA4HKCWD5AQOQR5HK6INZ3NM36LIGA6C
	return key.Secret(), nil
}

// 验证TOTP密码code
func VerifyTOTPPasscode(passcode, secret string) bool {
	// return totp.Validate(passcode, secret)

	result, _ := totp.ValidateCustom(
		passcode,
		secret,
		time.Now().UTC(),
		totp.ValidateOpts{
			Period:    30,
			Skew:      0,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		},
	)

	return result
}
