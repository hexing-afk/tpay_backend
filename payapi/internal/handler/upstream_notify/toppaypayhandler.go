package upstream_notify

import (
	"github.com/tal-tech/go-zero/core/logx"
	"io/ioutil"
	"net/http"
	"tpay_backend/payapi/internal/common"

	logic_ "tpay_backend/payapi/internal/logic/upstream_notify"
	"tpay_backend/payapi/internal/svc"
)

func ToppayPayHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 读取请求的body
		body, errs := ioutil.ReadAll(r.Body)
		if errs != nil {
			logx.Errorf("读取body失败:%v", errs)
			common.OkString(w, "failed")
			return
		}

		logx.Infof("toppay-notify-pay-body:%v", string(body))

		l := logic_.NewToppayPayLogic(r.Context(), ctx)
		err := l.ToppayPay(body)
		if err != nil {
			logx.Errorf("处理失败:%v", err)
			//common.OkString(w, "failed")
			common.OkString(w, `{
				"result":"success",
				"msg":"请求成功"
			}`)

		} else {
			common.OkString(w, `{
			"result":"error",
			"msg":"请求失败"
			}`)
		}
	}
}
