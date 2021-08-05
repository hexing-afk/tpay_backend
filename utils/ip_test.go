package utils

import (
	"encoding/json"
	"testing"
)

func TestRandomIp(t *testing.T) {
	a := GetFakeIp()
	t.Logf("a=%v", a)
}

func TestRandomIp2(t *testing.T) {
	a := "{\"code\":\"GP_00\",\"msg\":\"\",\"biz_code\":\"GPBIZ_00\",\"biz_msg\":\"\",\"data\":{\"sign_type\":\"RSA\",\"tf_sign\":\"1\",\"appid\":\"123\",\"businessnumber\":\"1\",\"businessrecordnumber\":\"1\",\"certcode\":\"2\",\"bankcardnumber\":\"2\",\"status\":\"1\"}}"
	m := map[string]interface{}{}
	json.Unmarshal([]byte(a), &m)
	t.Logf("a=%v", a)
}
