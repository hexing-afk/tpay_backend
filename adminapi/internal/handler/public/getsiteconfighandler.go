package public

import (
	"net/http"
	logic2 "tpay_backend/adminapi/internal/logic/public"

	"github.com/tal-tech/go-zero/rest/httpx"
	"tpay_backend/adminapi/internal/svc"
)

func GetSiteConfigHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic2.NewGetSiteConfigLogic(r.Context(), ctx)
		resp, err := l.GetSiteConfig()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
