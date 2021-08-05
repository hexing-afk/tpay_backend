package common

import "testing"

func TestCreateMerchantPassword(t *testing.T) {
	plainPassword := "123456"

	cipherPassword := CreateMerchantPassword(plainPassword)

	t.Logf("cipherPassword:%s", cipherPassword)
}
