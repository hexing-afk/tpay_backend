package merchant

import (
	"net/http"
	"tpay_backend/adminapi/internal/common"
	_func "tpay_backend/adminapi/internal/handler/func"
	"tpay_backend/adminapi/internal/logic/merchant"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func ModifyMerchantBalanceHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ModifyMerchantBalanceRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, common.NewCodeErrorWithMsg(common.VerifyParamFailed, err.Error()))
			return
		}

		// 从redis获取用户登录id
		userId, err := _func.GetLoginedUserIdRequestHeader(r)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		l := merchant.NewModifyMerchantBalanceLogic(r.Context(), ctx)
		err = l.ModifyMerchantBalance(userId, req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, nil)
		}
	}
}
