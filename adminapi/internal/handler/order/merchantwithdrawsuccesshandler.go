package order

import (
	"net/http"
	"tpay_backend/adminapi/internal/common"
	_func "tpay_backend/adminapi/internal/handler/func"

	logic "tpay_backend/adminapi/internal/logic/order"
	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func MerchantWithdrawSuccessHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MerchantWithdrawSuccessRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, common.NewCodeErrorWithMsg(common.VerifyParamFailed, err.Error()))
			return
		}

		userId, err := _func.GetLoginedUserIdRequestHeader(r)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewMerchantWithdrawSuccessLogic(r.Context(), ctx)
		err = l.MerchantWithdrawSuccess(userId, req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, nil)
		}
	}
}
