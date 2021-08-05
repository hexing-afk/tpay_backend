package utils

import "testing"

func TestGenerateTOTPSecret(t *testing.T) {
	accountName := "lch"
	secret, err := GenerateTOTPSecret(accountName)

	t.Logf("secret:%v", secret)
	t.Logf("err:%v", err)
}

func TestVerifyTOTPPasscode(t *testing.T) {
	secret := "2SFLPQSJFBZGJVMREDBGCEVANAHJWBNV"
	passcode := "523365"

	result := VerifyTOTPPasscode(passcode, secret)
	t.Logf("result:%v", result)
}
