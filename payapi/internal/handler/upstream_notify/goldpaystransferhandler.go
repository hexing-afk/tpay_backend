package upstream_notify

import (
	"github.com/tal-tech/go-zero/core/logx"
	"io/ioutil"
	"net/http"
	"tpay_backend/payapi/internal/common"

	"tpay_backend/payapi/internal/logic/upstream_notify"
	"tpay_backend/payapi/internal/svc"
)

func GoldPaysTransferHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 读取请求的body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logx.Errorf("读取body失败:%v", err)
			common.OkString(w, "failed")
			return
		}

		logx.Infof("goldpays-notify-transfer-body:%v", string(body))

		l := upstream_notify.NewGoldPaysTransferLogic(r.Context(), ctx)
		if err := l.GoldPaysTransfer(body); err != nil {
			logx.Errorf("处理失败:%v", err)
			common.OkString(w, "failed")
		} else {
			common.OkString(w, "success")
		}
	}
}
