package public

import (
	"net/http"
	"tpay_backend/adminapi/internal/common"

	logic_ "tpay_backend/adminapi/internal/logic/public"
	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func SaveOtherConfigHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SaveOtherConfigRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic_.NewSaveOtherConfigLogic(r.Context(), ctx)
		err := l.SaveOtherConfig(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, nil)
		}
	}
}
