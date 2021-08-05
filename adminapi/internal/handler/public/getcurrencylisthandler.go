package public

import (
	"net/http"
	"tpay_backend/adminapi/internal/common"
	logic2 "tpay_backend/adminapi/internal/logic/public"

	"github.com/tal-tech/go-zero/rest/httpx"
	"tpay_backend/adminapi/internal/svc"
)

func GetCurrencyListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic2.NewGetCurrencyListLogic(r.Context(), ctx)
		resp, err := l.GetCurrencyList()
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, resp)
		}
	}
}
