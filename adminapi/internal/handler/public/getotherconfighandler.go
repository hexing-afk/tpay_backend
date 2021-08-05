package public

import (
	"net/http"
	"tpay_backend/adminapi/internal/common"

	"github.com/tal-tech/go-zero/rest/httpx"
	logic_ "tpay_backend/adminapi/internal/logic/public"
	"tpay_backend/adminapi/internal/svc"
)

func GetOtherConfigHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic_.NewGetOtherConfigLogic(r.Context(), ctx)
		resp, err := l.GetOtherConfig()
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, resp)
		}
	}
}
