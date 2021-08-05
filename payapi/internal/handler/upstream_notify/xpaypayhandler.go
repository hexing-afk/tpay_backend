package upstream_notify

import (
	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"

	"tpay_backend/payapi/internal/common"
	logic "tpay_backend/payapi/internal/logic/upstream_notify"
	"tpay_backend/payapi/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

func XPayPayHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析接口数据
		var req logic.XPayNotifyRequest
		if err := httpx.Parse(r, &req); err != nil {
			logx.Errorf("参数解析失败:%v", err)
			common.OkString(w, "failed")
			return
		}

		logx.Infof("xpay-notify-pay:%+v", req)

		l := logic.NewXPayPayLogic(r.Context(), ctx)
		err := l.XPayPay(req)
		if err != nil {
			logx.Errorf("处理失败:%v", err)
			common.OkString(w, "failed")
		} else {
			common.OkString(w, "success")
		}
	}
}
