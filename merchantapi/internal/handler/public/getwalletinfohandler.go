package public

import (
	"net/http"
	"tpay_backend/merchantapi/internal/common"
	_func "tpay_backend/merchantapi/internal/handler/func"

	"github.com/tal-tech/go-zero/rest/httpx"
	_logic "tpay_backend/merchantapi/internal/logic/public"
	"tpay_backend/merchantapi/internal/svc"
)

func GetWalletInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId, errs := _func.GetLoginedUserIdRequestHeader(r)
		if errs != nil {
			httpx.Error(w, errs)
			return
		}

		l := _logic.NewGetWalletInfoLogic(r.Context(), ctx, userId)
		resp, err := l.GetWalletInfo()
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, resp)
		}
	}
}
