package merchant

import (
	"net/http"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/adminapi/internal/types"

	logic "tpay_backend/adminapi/internal/logic/merchant"
	"tpay_backend/adminapi/internal/svc"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func DeleteMerchantChannelHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteMerchantChannelRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, common.NewCodeErrorWithMsg(common.VerifyParamFailed, err.Error()))
			return
		}

		l := logic.NewDeleteMerchantChannelLogic(r.Context(), ctx)
		err := l.DeleteMerchantChannel(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, nil)
		}
	}
}
