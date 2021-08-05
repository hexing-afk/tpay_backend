package admin

import (
	"net/http"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/adminapi/internal/logic/admin"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func GetPlanformChannelSettingHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetPlanformChannelSettingRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := admin.NewGetPlanformChannelSettingLogic(r.Context(), ctx)
		resp, err := l.GetPlanformChannelSetting(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, resp)
		}
	}
}
