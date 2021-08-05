package admin

import (
	"net/http"
	"tpay_backend/adminapi/internal/common"

	"github.com/tal-tech/go-zero/rest/httpx"
	logic_ "tpay_backend/adminapi/internal/logic/admin"
	"tpay_backend/adminapi/internal/svc"
)

func AreaListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic_.NewAreaListLogic(r.Context(), ctx)
		resp, err := l.AreaList()
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, resp)
		}
	}
}
