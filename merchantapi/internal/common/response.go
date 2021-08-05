package common

import (
	"net/http"

	"github.com/tal-tech/go-zero/rest/httpx"
)

// 接口返回结构体
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (e *Response) Error() string {
	return e.Msg
}

// 错误处理函数
func ErrorHandler(err error) (int, interface{}) {
	// 默认返回错误类型
	res := Response{
		Code: SysUnKnow,
		Msg:  err.Error(),
	}

	if val, ok := err.(*Response); ok { // 如果逻辑层返回的错误是*Response类型
		res.Code = val.Code
		res.Msg = val.Msg
	}

	return http.StatusOK, res
}

// 成功时返回json数据
func OkJson(w http.ResponseWriter, v interface{}) {
	httpx.OkJson(w, Response{
		Code: Success,
		Msg:  SuccessMsg,
		Data: v,
	})
}
