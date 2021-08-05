package admin

import (
	"net/http"
	"tpay_backend/adminapi/internal/common"

	logic "tpay_backend/adminapi/internal/logic/admin"
	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func GetPlatformChannelListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetPlatformChannelListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, common.NewCodeErrorWithMsg(common.VerifyParamFailed, err.Error()))
			return
		}

		l := logic.NewGetPlatformChannelListLogic(r.Context(), ctx)
		resp, err := l.GetPlatformChannelList(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, resp)
		}
	}
}
