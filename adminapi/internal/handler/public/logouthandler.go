package public

import (
	"net/http"
	"tpay_backend/adminapi/internal/common"
	_func "tpay_backend/adminapi/internal/handler/func"

	logic "tpay_backend/adminapi/internal/logic/public"
	"tpay_backend/adminapi/internal/svc"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func LogoutHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(common.RequestHeaderLoginedToken)

		userId, err := _func.GetLoginedUserIdRequestHeader(r)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewLogoutLogic(r.Context(), ctx, userId)
		resp, err := l.Logout(token)
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, resp)
		}
	}
}
