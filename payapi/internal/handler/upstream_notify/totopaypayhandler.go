package upstream_notify

import (
	"io/ioutil"
	"net/http"
	"tpay_backend/payapi/internal/common"
	"tpay_backend/payapi/internal/logic/upstream_notify"

	"github.com/tal-tech/go-zero/core/logx"

	"tpay_backend/payapi/internal/svc"
)

func TotopayPayHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 读取请求的body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logx.Errorf("读取body失败:%v", err)
			common.OkString(w, "failed")
			return
		}

		logx.Infof("totopay-notify-pay-body:%v", string(body))

		l := upstream_notify.NewTotopayPayLogic(r.Context(), ctx)

		if err := l.TotopayPay(body); err != nil {
			logx.Errorf("处理失败:%v", err)
			common.OkString(w, "failed")
		} else {
			common.OkString(w, "success")
		}
	}
}
