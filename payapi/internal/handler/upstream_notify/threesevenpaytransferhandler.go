package upstream_notify

import (
	"github.com/tal-tech/go-zero/core/logx"
	"io/ioutil"
	"net/http"
	"tpay_backend/payapi/internal/common"

	logic "tpay_backend/payapi/internal/logic/upstream_notify"
	"tpay_backend/payapi/internal/svc"
)

func ThreeSevenPayTransferHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 读取请求的body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logx.Errorf("读取body失败:%v", err)
			common.OkString(w, "failed")
			return
		}

		logx.Infof("zf777pay-notify-transfer-body:%v", string(body))

		l := logic.NewThreeSevenPayTransferLogic(r.Context(), ctx)
		err = l.ThreeSevenPayTransfer(body)
		if err != nil {
			logx.Errorf("处理失败:%v", err)
			//common.OkString(w, "failed")
			common.OkString(w, "success")
		} else {
			common.OkString(w, "success")
		}
	}
}
