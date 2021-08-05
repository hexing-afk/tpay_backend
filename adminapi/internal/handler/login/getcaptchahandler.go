package login

import (
	"net/http"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/adminapi/internal/logic/login"
	"tpay_backend/adminapi/internal/svc"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func GetCaptchaHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := login.NewGetCaptchaLogic(r.Context(), ctx)
		resp, err := l.GetCaptcha()
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, resp)
		}
	}
}
