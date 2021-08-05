package utils

import (
	"bytes"
	"fmt"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// 获取http请求的body
func GetHttpRequestBody(r *http.Request) []byte {
	// 获取body字节流数据
	body, _ := ioutil.ReadAll(r.Body)

	// 因为body是ReadCloser类型，读完就没有了，方便后续继续使用，需要将数据写回Body用于传递
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return body
}

// ClientIP 尽最大努力实现获取客户端 IP 的算法。
// 解析 X-Real-IP 和 X-Forwarded-For 以便于反向代理（nginx 或 haproxy）可以正常工作。
func ClientIP(r *http.Request) string {
	ip := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

// 修剪逗号分隔的数据(去掉多余的空格和空值)
// "a, b, c" => "a,b,c"
// "a, b, c," => "a,b,c"
func TrimCommaStr(str string) string {
	arr := strings.Split(str, ",")

	var list []string

	for _, s := range arr {
		s = strings.TrimSpace(s)
		if s != "" {
			list = append(list, s)
		}
	}

	return strings.Join(list, ",")
}

// 判断某个字符串是否在slice中, 类似php的in_array()函数
func InSlice(value string, slice []string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

/**
获取新uuid string
*/
func NewUUID() string {
	uuidB := uuid.NewV4()
	return uuidB.String()
}

// 按slice中的值的权重来选择元素
//
// 算法描述:
// 参数: {20, 30, 50}
// 1.计算总权重 = 20+30+50 = 100
// 2.分配权重区域: [1,21)[21,51)[51,101)
// 3.在[1,100]之间随机一个数(100来自总权重值),判断随机数在哪个范围区间,则返回命中区间的索引
// 假如随机数是89, 则命中在[51,101)范围之间,返回索引2(即值50的索引)
// 注意：为了方便计算，每个元素的权重值不能小于0
func PickByWeight(list []int64) int {
	if len(list) == 0 {
		panic("Parameter cannot be empty")
	}

	if len(list) == 1 { // 如果只有1个元素，直接返回索引0
		return 0
	}

	var total int64 // 总值

	for _, v := range list {
		if v < 0 {
			panic("The weight value cannot be less than 0")
		}
		total += v
	}

	// 随机挑选一个值
	hit := RandInt64(1, total)

	min := int64(1)
	for i, v := range list {
		max := min + v
		if min <= hit && hit < max { // 命中范围
			return i
		} else {
			min = max // 调整范围最小值
		}
	}

	panic(fmt.Sprintf("Failed to select weight value, slice:%v", list))
}

// 根据数字范围获取随机数,取值范围为[min,max]
func RandInt64(min, max int64) int64 {
	max += 1 // 因为Int63n取值范围为[0,n),加1是满足[0,n]
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min
}

// 生成时间唯一id
func GetDailyId() string {
	z, _ := time.LoadLocation("Asia/Shanghai")
	t0 := time.Now().In(z)
	t := t0.Format("20060102150405")
	i := RandInt64(10000, 99999)
	return fmt.Sprintf("%s%03d%05d", t, t0.Nanosecond()/1000%1000, i)
}

// 生成20位唯一id
func GetUniqueId() string {
	id := time.Now().UnixNano() / 1000
	return fmt.Sprintf("%v%v", id, RandInt64(1000, 9999))
}

// 生成指定长度的字符串
func RandString(l int) string {
	digitNumber := []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	}

	var str string
	for j := 0; j < l; j++ {
		a := digitNumber[rand.Intn(len(digitNumber))]
		str += a
	}

	return str
}

// 将参数排序并拼接成字符串
func ParamsMapToString(params map[string]interface{}, excludeField string) string {
	var pList = make([]string, 0)

	for key, value := range params {
		if strings.TrimSpace(key) == excludeField { // 忽略验签字段
			continue
		}

		// 将interface转换为字符串
		val := strings.TrimSpace(SingleTypeToStringNoExponent(value))
		if len(val) > 0 { // 忽略空值
			pList = append(pList, key+"="+val)
		}
	}

	// 按键排序
	sort.Strings(pList)

	// 使用&符号拼接
	return strings.Join(pList, "&")
}

// 单一类型转换为字符串 (浮点数和整数)不使用科学计数法
// 12.10 => 12.1 (浮点数尾部有0会被自动去掉)
func SingleTypeToStringNoExponent(x interface{}) string {
	var v2 string

	switch v := x.(type) {

	case int:
		v2 = strconv.FormatInt(int64(v), 10)
	case int32:
		v2 = strconv.FormatInt(int64(v), 10)
	case int64:
		v2 = strconv.FormatInt(v, 10)

	case uint32:
		v2 = strconv.FormatUint(uint64(v), 10)
	case uint64:
		v2 = strconv.FormatUint(v, 10)

	case float32:
		v2 = strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		v2 = strconv.FormatFloat(v, 'f', -1, 64)

	default:
		v2 = fmt.Sprintf("%v", v)
	}

	return v2
}

// 获取正在运行的函数名
func RunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

/*func Post(url string, data url.Values) ([]byte, error) {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}*/

func PostForm(url string, data url.Values) ([]byte, error) {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

func PostJson(url string, dataStr []byte) ([]byte, error) {
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(dataStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

// 去掉数组中重复的元素
func DistinctSliceOnIt64(arr []int64) (newArr []int64) {
	newArr = make([]int64, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}

// 对比原数组和新数组，
// 找出需要新增到原数组中的元素，和需要从原数组中删除的元素
func TemplateSliceInt64(oldS, newS []int64) (add, del []int64) {
	var addU, delU []int64
	source := newS
	for _, v := range oldS {
		// 原先有，现在没有
		idx := getKey(v, source)
		if idx < 0 {
			delU = append(delU, v)
		} else {
			source = remove(idx, source)
		}
	}
	addU = source

	return addU, delU
}

func getKey(value int64, lList []int64) int {
	for k, v := range lList {
		if v == value {
			return k
		}
	}
	return -1
}

func remove(idx int, lList []int64) []int64 {
	return append(lList[:idx], lList[(idx+1):]...)
}

func TimeSubToDay(start, end int64) int64 {
	subSecond := time.Unix(end, 0).Sub(time.Unix(start, 0)).Seconds()

	// 1天 = 86400秒
	day := decimal.NewFromFloat(subSecond).Div(decimal.NewFromInt(86400)).Floor().IntPart()

	return day
}
