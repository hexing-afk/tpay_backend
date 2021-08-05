package order

import (
	"net/http"
	logic "tpay_backend/adminapi/internal/logic/order"

	"tpay_backend/adminapi/internal/common"
	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func MerchantWithdrawOrderAuditHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MerchantWithdrawOrderAuditRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, common.NewCodeErrorWithMsg(common.VerifyParamFailed, err.Error()))
			return
		}

		l := logic.NewMerchantWithdrawOrderAuditLogic(r.Context(), ctx)
		err := l.MerchantWithdrawOrderAudit(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, nil)
		}
	}
}
