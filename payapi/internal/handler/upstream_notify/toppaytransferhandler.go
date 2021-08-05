package upstream_notify

import (
	"io/ioutil"
	"net/http"

	"tpay_backend/payapi/internal/common"
	logic "tpay_backend/payapi/internal/logic/upstream_notify"
	"tpay_backend/payapi/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

func ToppayTransferHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析接口数据
		body, errs := ioutil.ReadAll(r.Body)
		if errs != nil {
			logx.Errorf("读取body失败:%v", errs)
			common.OkString(w, "failed")
			return
		}

		logx.Infof("toppay-notify-pay-body:%v", string(body))

		l := logic.NewToppayTransferLogic(r.Context(), ctx)
		err := l.ToppayTransfer(body)
		if err != nil {
			logx.Errorf("处理失败:%v", err)
			common.OkString(w, "failed")
		} else {
			common.OkString(w, "success")
		}
	}
}
