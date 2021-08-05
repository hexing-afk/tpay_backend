package common

import "testing"

func TestCreateAdminPassword(t *testing.T) {
	plainPassword := "123456"

	cipherPassword := CreateAdminPassword(plainPassword)

	t.Logf("cipherPassword:%s", cipherPassword)
}
