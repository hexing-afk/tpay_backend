package admin

import (
	"net/http"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	logic "tpay_backend/adminapi/internal/logic/admin"
	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func AddUpstreamChannelHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddUpstreamChannelRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, common.NewCodeErrorWithMsg(common.VerifyParamFailed, err.Error()))
			return
		}

		if req.ChannelType != model.UpstreamChannelTypeCollection && req.ChannelType != model.UpstreamChannelTypeTransfer {
			httpx.Error(w, common.NewCodeError(common.InvalidParam))
		}

		l := logic.NewAddUpstreamChannelLogic(r.Context(), ctx)
		err := l.AddUpstreamChannel(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, nil)
		}
	}
}
