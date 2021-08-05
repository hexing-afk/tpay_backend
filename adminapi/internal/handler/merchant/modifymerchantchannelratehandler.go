package merchant

import (
	"github.com/tal-tech/go-zero/rest/httpx"
	"net/http"
	"tpay_backend/adminapi/internal/common"
	logic "tpay_backend/adminapi/internal/logic/merchant"
	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"
)

func ModifyMerchantChannelRateHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ModifyMerchantChannelRateRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, common.NewCodeErrorWithMsg(common.VerifyParamFailed, err.Error()))
			return
		}

		l := logic.NewModifyMerchantChannelRateLogic(r.Context(), ctx)
		err := l.ModifyMerchantChannelRate(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, nil)
		}
	}
}
