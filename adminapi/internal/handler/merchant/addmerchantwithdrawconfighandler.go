package merchant

import (
	"net/http"
	"tpay_backend/adminapi/internal/common"
	logic "tpay_backend/adminapi/internal/logic/merchant"
	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func AddMerchantWithdrawConfigHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddMerchantWithdrawConfigRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, common.NewCodeErrorWithMsg(common.VerifyParamFailed, err.Error()))
			return
		}

		if req.SingleMaxAmount < 0 || req.SingleMinAmount < 0 || req.SingleMaxAmount < req.SingleMinAmount {
			httpx.Error(w, common.NewCodeError(common.InvalidParam))
			return
		}

		l := logic.NewAddMerchantWithdrawConfigLogic(r.Context(), ctx)
		resp, err := l.AddMerchantWithdrawConfig(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, resp)
		}
	}
}
