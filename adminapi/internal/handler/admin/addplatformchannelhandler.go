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

func AddPlatformChannelHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddPlatformChannelRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, common.NewCodeErrorWithMsg(common.VerifyParamFailed, err.Error()))
			return
		}

		if req.ChannelType != model.PlatformChannelTypeCollection && req.ChannelType != model.PlatformChannelTypeTransfer {
			httpx.Error(w, common.NewCodeError(common.InvalidParam))
		}

		l := logic.NewAddPlatformChannelLogic(r.Context(), ctx)
		err := l.AddPlatformChannel(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, nil)
		}
	}
}
