package utils

import (
	"bytes"
	"encoding/hex"
	"io"
	"log"
	"math"
	"math/rand"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"
)

/**
 * 初始化一下种子
 */
func RandomInit() {
	rand.Seed(time.Now().UnixNano())
}

func Random(min, max int64) int64 {
	if min == max {
		return min
	}

	//rand.Seed(time.Now().UnixNano())
	//fmt.Println("["+strconv.FormatInt(min, 10)+"]"+"["+strconv.FormatInt(max, 10)+"]")
	return rand.Int63n(max-min) + min
}

func RandomFloat64(min, max float64) float64 {
	// 浮点数比较
	if math.Abs(min-max) < 0.00000001 {
		return min
	}

	//rand.Seed(time.Now().UnixNano())
	//fmt.Println("["+strconv.FormatInt(min, 10)+"]"+"["+strconv.FormatInt(max, 10)+"]")
	return rand.Float64()*(max-min) + min
}

//
////RandomString 在数字、大写字母、小写字母范围内生成length位的随机字符串，所有字符出现的概率相同
//func RandomDigitStr(length int32) string {
//	// 48 ~ 57 数字
//	// 65 ~ 90 A ~ Z
//	// 97 ~ 122 a ~ z
//	// 一共62个字符，在0~61进行随机，小于10时，在数字范围随机，
//	// 小于36在大写范围内随机，其他在小写范围随机
//	rand.Seed(time.Now().UnixNano())
//	result := make([]string, 0, length)
//	var i int32
//	for i = 0; i < length; i++ {
//		t := rand.Intn(62)
//		if t < 10 {
//			result = append(result, strconv.Itoa(t))
//		} else if t < 36 {
//			result = append(result, string(rand.Intn(26)+65))
//		} else {
//			result = append(result, string(rand.Intn(26)+97))
//		}
//	}
//	return strings.Join(result, "")
//}

func RandomDigitStrOnlyNum(length int32) string {
	// 48 ~ 57 数字
	// 65 ~ 90 A ~ Z
	// 97 ~ 122 a ~ z
	// 一共62个字符，在0~61进行随机，小于10时，在数字范围随机，
	// 小于36在大写范围内随机，其他在小写范围随机
	rand.Seed(time.Now().UnixNano())
	result := make([]string, 0, length)
	var i int32
	for i = 0; i < length; i++ {
		t := rand.Intn(10)
		result = append(result, strconv.Itoa(t))
	}
	return strings.Join(result, "")
}

//func RandomDigitStrOnlyAlphabet(length int32) string {
//	// 97 ~ 122 a ~ z
//	rand.Seed(time.Now().UnixNano())
//	result := make([]string, 0, length)
//	var i int32
//	for i = 0; i < length; i++ {
//		result = append(result, string(rand.Intn(26)+97))
//	}
//	return strings.Join(result, "")
//}

//func RandomDigitStrOnlyAlphabetUpper(length int32) string {
//	// 65 ~ 90 A ~ Z
//	rand.Seed(time.Now().UnixNano())
//	result := make([]string, 0, length)
//	var i int32
//	for i = 0; i < length; i++ {
//		result = append(result, string(rand.Intn(26)+65))
//	}
//	return strings.Join(result, "")
//}

//
func RandomDigitHex(length int32) string {
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, 0, length)
	var i int32
	for i = 0; i < length; i++ {
		t := rand.Intn(255)
		result = append(result, byte(t))
	}
	return hex.EncodeToString(result)
}

// 只能在0~19位之间
func RandomDigit(digit int32) int64 {
	if digit <= 0 {
		return 0
	} else if digit > 19 {
		return 0
	}

	min := int64(math.Pow10(int(digit) - 1))
	max := int64(math.Pow10(int(digit))) - 1
	return Random(min, max)
}

// 生成随机字符串
func RandomStringLower(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytesSlice := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytesSlice[r.Intn(len(bytesSlice))])
	}
	return string(result)
}

/**
 * 打印异常
 */
func CheckNPrintError(err error) {
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func WriteHttpFile(filepath, httpFileKey string, otherHttpParams map[string]string) (data io.Reader, contentType string, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	body := bytes.Buffer{}
	writer := multipart.NewWriter(&body)
	defer writer.Close()
	part, err := writer.CreateFormFile(httpFileKey, filepath)
	if err != nil {
		return nil, "", err
	}
	_, err = io.Copy(part, file)
	log.Printf("%v\n", err)
	for key, val := range otherHttpParams {
		log.Printf("field=[%v]->[%v]\n", key, val)
		_ = writer.WriteField(key, val)
	}
	return &body, writer.FormDataContentType(), nil
}

func BetweenInt32(min, max, value int32) int32 {
	if min > max {
		t := min
		min = max
		max = t
	}
	if value >= min && value <= max {
		return value
	} else if value < min {
		return min
	}

	return max
}

func BetweenInt64(min, max, value int64) int64 {
	if min > max {
		t := min
		min = max
		max = t
	}
	if value >= min && value <= max {
		return value
	} else if value < min {
		return min
	}

	return max
}

func AtleastInt64(min, value int64) int64 {
	if value >= min {
		return value
	}

	return min
}

func FixDig(str string, n int, isLeft bool, spaceStr string) string {
	lenStr := len(str)
	str2 := str
	if lenStr < n {
		for i := 0; i < n-lenStr; i++ {
			if isLeft {
				str2 = spaceStr + str2
			} else {
				str2 += spaceStr
			}
		}
	}
	return str2
}

func Goooooooooooooooo(
	loopFunc func(...interface{}), // 每次循环执行这个
	args ...interface{}, // 给loopFunc的参数
) {
	go func() {
		defer func() {
			log.Printf("报错了,Goooooooooooooooo中断")
			recover()
		}()

		//log.Printf("runtime.NumGoroutine()=[%v]", runtime.NumGoroutine())
		loopFunc(args...)
	}()
}

func GooooooooooooooooDeadLoop(
	loopFunc func(...interface{}), // 每次循环执行这个
	args ...interface{}, // 给loopFunc的参数
) {
	go func() {
		defer func() {
			log.Printf("报错了,Goooooooooooooooo中断")
			recover()
			GooooooooooooooooDeadLoop(loopFunc, args...)
		}()

		//log.Printf("runtime.NumGoroutine()=[%v]", runtime.NumGoroutine())
		loopFunc(args...)
	}()
}

type Int64Slice []int64

func (c Int64Slice) Len() int {
	return len(c)
}
func (c Int64Slice) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c Int64Slice) Less(i, j int) bool {
	return c[i] < c[j]
}
