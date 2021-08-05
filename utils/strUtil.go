package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"regexp"

	"math"
	"strconv"
	"strings"
	"time"
)

func ToString(x interface{}) string {
	var v2 string

	switch v := x.(type) {
	case bool:
		if true == x {
			v2 = "1"
		} else {
			v2 = "0"
		}
	case int:
		v2 = strconv.Itoa(v)
	case int32:
		v2 = strconv.Itoa(int(v))
	case int64:
		v2 = strconv.FormatInt(v, 10)
	case float64:
		v2 = strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		v2 = v
	case []byte:
		v2 = string(v)
	case map[string]interface{}:
		v2b, _ := json.Marshal(v)
		v2 = string(v2b)
	case []map[string]interface{}:
		v2b, _ := json.Marshal(v)
		v2 = string(v2b)
	case []string:
		v2b, _ := json.Marshal(v)
		v2 = string(v2b)
	case nil:
		v2 = ""
	}

	return v2
}

func ToStringNoPoint(x interface{}) string {
	var v2 string

	switch v := x.(type) {
	case bool:
		if true == x {
			v2 = "1"
		} else {
			v2 = "0"
		}
	case int:
		v2 = strconv.Itoa(v)
	case int32:
		v2 = strconv.Itoa(int(v))
	case int64:
		v2 = strconv.FormatInt(v, 10)
	case float64:
		v2 = RemoveZero(v)
	case string:
		v2 = v
	case []byte:
		v2 = string(v)
	case map[string]interface{}:
		v2b, _ := json.Marshal(v)
		v2 = string(v2b)
	case *map[string]interface{}:
		v2b, _ := json.Marshal(v)
		v2 = string(v2b)
	case []map[string]interface{}:
		v2b, _ := json.Marshal(v)
		v2 = string(v2b)
	case []string:
		v2 = strings.Join(v, ",")
	case []float64:
		tmp := []string{}
		for _, f := range v {
			tmp = append(tmp, ToStringNoPoint(f))
		}
		v2 = ToJson(tmp)
	case []interface{}:
		tmp := []string{}
		for _, v3 := range v {
			tmp = append(tmp, ToStringNoPoint(v3))
		}
		v2 = ToJson(tmp)
	case time.Time:
		v2 = v.Format("2006-01-02 15:04:05")
	case time.Duration:
		v2 = v.String()
	case nil:
		v2 = ""
	}

	return strings.Trim(v2, " ")
}

func ToFloat64(x interface{}) float64 {
	var v2 float64

	switch v := x.(type) {
	case bool:
		if true == x {
			v2 = 1
		} else {
			v2 = 0
		}
	case int:
		v2 = float64(v * 1.0)
	case int32:
		v2 = float64(v * 1.0)
	case int64:
		v2 = float64(v * 1.0)
	case float32:
		v2 = float64(v)
	case float64:
		v2 = v
	case string:
		v2, _ = strconv.ParseFloat(v, 64)
	case []byte:
		v2, _ = strconv.ParseFloat(string(v), 64)
	case nil:
		v2 = 0.0
	}

	return v2
}

func ToFloat32(x interface{}) float32 {
	var v2 float32

	switch v := x.(type) {
	case bool:
		if true == x {
			v2 = 1
		} else {
			v2 = 0
		}
	case int:
		v2 = float32(v * 1.0)
	case int32:
		v2 = float32(v * 1.0)
	case int64:
		v2 = float32(v * 1.0)
	case float32:
		v2 = v
	case float64:
		v2 = float32(v)
	case string:
		vTmp, _ := strconv.ParseFloat(string(v), 64)
		v2 = float32(vTmp)
	case []byte:
		vTmp, _ := strconv.ParseFloat(string(v), 64)
		v2 = float32(vTmp)
	case nil:
		v2 = float32(0.0)
	}

	return v2
}

func ToInt(x interface{}) int {
	var v2 int

	switch v := x.(type) {
	case bool:
		if true == x {
			v2 = 1
		} else {
			v2 = 0
		}
	case int:
		v2 = v
	case int32:
		v2 = int(v)
	case int64:
		v2 = int(v)
	case float32:
		v2 = int(v)
	case float64:
		v2 = int(v)
	case string:
		v2, _ = strconv.Atoi(v)
	case []byte:
		v2, _ = strconv.Atoi(string(v))
	case nil:
		v2 = int(0)
	}

	return v2
}

func ToInt16(x interface{}) int16 {
	var v2 int16

	switch v := x.(type) {
	case bool:
		if true == x {
			v2 = 1
		} else {
			v2 = 0
		}
	case int:
		v2 = int16(v)
	case int16:
		v2 = int16(v)
	case int32:
		v2 = int16(v)
	case int64:
		v2 = int16(v)
	case float32:
		v2 = int16(v)
	case float64:
		v2 = int16(v)
	case string:
		t, _ := strconv.Atoi(v)
		v2 = int16(t)
	case []byte:
		t, _ := strconv.Atoi(string(v))
		v2 = int16(t)
	case nil:
		v2 = int16(0)
	}

	return v2
}

func ToUint16(x interface{}) uint16 {
	var v2 uint16

	switch v := x.(type) {
	case bool:
		if true == x {
			v2 = 1
		} else {
			v2 = 0
		}
	case int:
		v2 = uint16(v)
	case int16:
		v2 = uint16(v)
	case int32:
		v2 = uint16(v)
	case int64:
		v2 = uint16(v)
	case float32:
		v2 = uint16(v)
	case float64:
		v2 = uint16(v)
	case string:
		t, _ := strconv.Atoi(v)
		v2 = uint16(t)
	case []byte:
		t, _ := strconv.Atoi(string(v))
		v2 = uint16(t)
	case nil:
		v2 = uint16(0)
	}

	return v2
}

func ToInt64(x interface{}) int64 {
	var v2 int64

	switch v := x.(type) {
	case bool:
		if true == x {
			v2 = int64(1)
		} else {
			v2 = 0
		}
	case int:
		v2 = int64(v)
	case int32:
		v2 = int64(v)
	case int64:
		v2 = int64(v)
	case float32:
		v2 = int64(v)
	case float64:
		v2 = int64(v)
	case []byte:
		vTmp, _ := strconv.Atoi(string(v))
		v2 = int64(vTmp)
	case string:
		vTmp, _ := strconv.Atoi(v)
		v2 = int64(vTmp)
	case nil:
		v2 = int64(0)
	}

	return v2
}

func ToInt32(x interface{}) int32 {
	var v2 int32

	switch v := x.(type) {
	case bool:
		if true == x {
			v2 = int32(1)
		} else {
			v2 = int32(0)
		}
	case int:
		v2 = int32(v)
	case int32:
		v2 = int32(v)
	case int64:
		v2 = int32(v)
	case float32:
		v2 = int32(v)
	case float64:
		v2 = int32(v)
	case string:
		vTmp, _ := strconv.Atoi(v)
		v2 = int32(vTmp)
	case []byte:
		vTmp, _ := strconv.Atoi(string(v))
		v2 = int32(vTmp)
	case nil:
		v2 = int32(0)
	}

	return v2
}

func ToBool(x interface{}) bool {
	v2 := false
	switch v := x.(type) {
	case bool:
		v2 = v
	case int:
		v2 = v == 1
	case int32:
		v2 = v == 1
	case int64:
		v2 = v == 1
	case float32:
		v2 = v == 1
	case float64:
		v2 = v == 1
	case string:
		if string(v) == "true" || v == "1" || v == fmt.Sprint(0x01) {
			v2 = true
		}
	case []byte:
		if string(v) == "true" || string(v) == "1" || string(v) == fmt.Sprint(0x01) {
			v2 = true
		}
	case nil:
		v2 = false
	}

	return v2
}

func NewUUIDNoSplit() string {
	return strings.Replace(NewUUID(), "-", "", -1)
}

func ToStringList(r []interface{}) []string {
	var lst []string
	for _, v := range r {
		lst = append(lst, v.(string))
	}
	return lst
}

func ToUUid(x interface{}) string {
	var v2 string

	switch v := x.(type) {
	case bool:
		if true == x {
			v2 = "1"
		} else {
			v2 = "0"
		}
	case int:
		v2 = strconv.Itoa(v)
	case int32:
		v2 = strconv.Itoa(int(v))
	case int64:
		v2 = strconv.FormatInt(v, 10)
	case float64:
		v2 = strconv.FormatFloat(v, 'f', 2, 64)
	case string:
		v2 = v
	case []byte:
		v2 = string(v)
	case nil:
		v2 = ""
	}
	if v2 == "" {
		return "00000000-0000-0000-0000-000000000000"
	} else {
		return v2

	}
}

func RemoveZero(num float64) string {
	numStr := ToString(num)
	lenNumStr := len(numStr)
	j := lenNumStr
	isRemove := false
	for i := lenNumStr; i > 0; i-- {
		if numStr[i-1:i] == "." {
			isRemove = true
		}
		if numStr[i-1:i] == "." || numStr[i-1:i] == "0" {
			j--
		} else {
			break
		}
	}

	if len(numStr[:j]) <= 0 {
		return "0"
	}

	if isRemove {
		return numStr[:j]
	}
	return numStr
}

func DeSensitive(code string) string {
	lenCode := len(code)
	if lenCode < 7 {
		return code
	}

	return code[:3] + strings.Repeat("*", lenCode-7) + code[(lenCode-4):]
}

/**
 * 四舍五入
 */
func RoundingOff(f float64, n int) float64 {
	pow10_n := math.Pow10(n)
	return math.Trunc((f+0.5/pow10_n)*pow10_n) / pow10_n
}

/**
 * 去掉空格和开头的0
 * XXX 一定要注意,如果字串中间存在空格,空格后面的值会抹掉
 */
func CleanSpace(str0 string) string {
	var str1 string
	fGot := false
	for _, v := range str0 {
		v0 := string(v)
		if !fGot && (v0 == "0" || v0 == " ") {
			continue
		}
		if !fGot && v0 != "0" && v0 != " " {
			fGot = true
		}
		if fGot && v0 == " " {
			break
		}
		str1 += v0
	}
	return str1
}
func CleanSpaceTail(str0 string) string {
	if str0 == "" {
		return ""
	}

	var str1 string
	lenStr := len(str0)
	if lenStr <= 2 {
		if string(str0[0]) != " " {
			str1 += string(str0[0])
		}
		if string(str0[1]) != " " {
			str1 += string(str0[1])
		}
		return str1
	}
	i := 0
	for i = lenStr; i > 1; i-- {
		//log.Printf("str0=[%v:%v:%v]", str0, str0[i-1:i], i)
		v0 := string(str0[i-1 : i])
		if v0 != " " {
			break
		}
	}
	return str0[:i]
}

func ToStringNoPointNCleanSpace(var0 interface{}, defValue ...string) string {
	var1 := CleanSpace(ToStringNoPoint(var0))
	if var1 == "" && len(defValue) > 0 {
		return defValue[0]
	}
	return var1
}

func ToStringNoPointNCleanSpaceTail(var0 interface{}, defValue ...string) string {
	var1 := CleanSpaceTail(ToStringNoPoint(var0))
	if var1 == "" && len(defValue) > 0 {
		return defValue[0]
	}
	return var1
}

func IsJson(jsonStr string) bool {
	return json.Valid([]byte(jsonStr))
}

func Float64ToStringFixPoint(v float64, bit int) string {
	return strconv.FormatFloat(v, 'f', bit, 64)
}

func Float64ToPrice(v float64) string {
	return Float64ToStringFixPoint(v, 2)
}

func Int64ToPrice(v int64) string {
	return fmt.Sprintf("%d.%d", v/100, v%100)
}

func StringToPrice(v string) string {
	v2 := ToInt64(v)
	return fmt.Sprintf("%d.%02d", v2/100, v2%100)
}

func Float64Mul(a, b float64) decimal.Decimal {
	f1 := decimal.NewFromFloat(a)
	f2 := decimal.NewFromFloat(b)
	return f1.Mul(f2)
}

func RandomIp() string {
	ip1 := Random(1, 255)
	ip2 := Random(0, 255)
	ip3 := Random(0, 255)
	ip4 := Random(1, 255)
	return fmt.Sprintf("%d.%d.%d.%d", ip1, ip2, ip3, ip4)
}

func Json2Map(jsonB interface{}) map[string]interface{} {
	m := map[string]interface{}{}
	var jsonB2 []byte
	switch v := jsonB.(type) {
	case string:
		jsonB2 = []byte(v)
	case []byte:
		jsonB2 = v
	}
	err := json.Unmarshal(jsonB2, &m)
	if err != nil {
		fmt.Printf("err=[%v]\n", err)
		return nil
	}
	return m
}

func Json2StrList(jsonB interface{}) []string {
	m := []string{}
	var jsonB2 []byte
	switch v := jsonB.(type) {
	case string:
		jsonB2 = []byte(v)
	case []byte:
		jsonB2 = v
	}
	if len(jsonB2) <= 0 {
		return m
	}
	err := json.Unmarshal(jsonB2, &m)
	if err != nil {
		fmt.Printf("err=[%v]\n", err)
		return nil
	}
	return m
}

func ToCapExt(a, sep string) (c string) {
	b := strings.Split(a, sep)
	for _, v := range b {
		c = fmt.Sprintf("%s%s%s", c, strings.ToUpper(v[:1]), v[1:])
	}
	return c
}

func ToCapExtNSkip(a, sep string, isCapLower bool) (c string) {
	b := strings.Split(a, sep)
	for _, v := range b {
		c = fmt.Sprintf("%s%s%s", c, strings.ToUpper(v[:1]), v[1:])
	}

	var c2 string
	if isCapLower {
		c2 = strings.ToLower(string(c[:1])) + c[1:]
	} else {
		c2 = c
	}

	return c2
}

func ToJson(m interface{}) string {
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("toJsonErr=[%v]", err)
		return ""
	}
	return string(b)
}

func ToJsonNotChange(m interface{}) string {
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Printf("toJsonErr=[%v]", err)
		return ""
	}
	return strings.Replace(string(b), "\\u0026", "&", -1)
}

func UnderLineToCamel(field string) string {
	reg := regexp.MustCompile("_")
	found := reg.FindAllStringIndex(field, -1)
	var extStr string
	lenFound := len(found)
	if lenFound > 0 {
		for k2, v2 := range found {
			if k2+1 < lenFound {
				extStr += strings.ToUpper(string(field[v2[1]:v2[1]+1])) + field[v2[1]+1:found[k2+1][0]]
			} else {
				extStr += strings.ToUpper(string(field[v2[1]:v2[1]+1])) + field[v2[1]+1:]
			}
		}
		return strings.ToUpper(string(field[0])) + field[1:found[0][0]] + extStr
	} else {
		field = strings.Replace(field, `"`, "", -1)
		field = strings.Replace(field, "`", "", -1)
		return strings.ToUpper(string(field[0])) + field[1:]
	}
}

func UnderLineToCamel1stLower(field string) string {
	reg := regexp.MustCompile("_")
	found := reg.FindAllStringIndex(field, -1)
	var extStr string
	lenFound := len(found)
	if lenFound > 0 {
		for k2, v2 := range found {
			if k2+1 < lenFound {
				extStr += strings.ToUpper(string(field[v2[1]:v2[1]+1])) + field[v2[1]+1:found[k2+1][0]]
			} else {
				extStr += strings.ToUpper(string(field[v2[1]:v2[1]+1])) + field[v2[1]+1:]
			}
		}
		return strings.ToLower(string(field[0])) + field[1:found[0][0]] + extStr
	} else {
		field = strings.Replace(field, `"`, "", -1)
		field = strings.Replace(field, "`", "", -1)
		return strings.ToLower(string(field[0])) + field[1:]
	}
}

func ForShort(msg string, n int) string {
	if len(msg) < n {
		return msg
	}
	return msg[:n]
}

func StringToFloat64(s string) float64 {
	fromString, _ := decimal.NewFromString(s)
	f, _ := fromString.Float64()
	return f
}

func Float64ToString(f float64) string {
	return decimal.NewFromFloat(f).String()
}

func FLoatCmp(a, b interface{}) (int, error) {
	switch v := a.(type) {
	case float32:
		switch v2 := b.(type) {
		case float32:
			return decimal.NewFromFloat32(v).Cmp(decimal.NewFromFloat32(v2)), nil
		case float64:
			return decimal.NewFromFloat32(v).Cmp(decimal.NewFromFloat(v2)), nil
		}
	case float64:
		switch v2 := b.(type) {
		case float32:
			return decimal.NewFromFloat(v).Cmp(decimal.NewFromFloat32(v2)), nil
		case float64:
			return decimal.NewFromFloat(v).Cmp(decimal.NewFromFloat(v2)), nil
		}
	default:
		return 1, errors.New("not support type")
	}
	return 1, errors.New("not support type")
}
