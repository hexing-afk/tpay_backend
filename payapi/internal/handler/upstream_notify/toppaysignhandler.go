package upstream_notify

import (
	"github.com/tal-tech/go-zero/core/logx"
	"io/ioutil"
	"net/http"
	"tpay_backend/payapi/internal/common"

	"github.com/tal-tech/go-zero/rest/httpx"
	logic_ "tpay_backend/payapi/internal/logic/upstream_notify"
	"tpay_backend/payapi/internal/svc"
)

func ToppaySignHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// 读取请求的body
		body, errs := ioutil.ReadAll(r.Body)
		if errs != nil {
			logx.Errorf("读取body失败:%v", errs)
			common.OkString(w, "failed")
			return
		}

		logx.Infof("toppay-notify-sign-body:%v", string(body))

		l := logic_.NewToppaySignLogic(r.Context(), ctx)
		err := l.ToppaySign(body)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
