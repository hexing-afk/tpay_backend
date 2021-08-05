package merchant

import (
	"net/http"
	"tpay_backend/adminapi/internal/common"
	_func "tpay_backend/adminapi/internal/handler/func"

	logic_ "tpay_backend/adminapi/internal/logic/merchant"
	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func UpdateMerchantPwdHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateMerchantPwdRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		// 从redis获取用户登录id
		userId, errs := _func.GetLoginedUserIdRequestHeader(r)
		if errs != nil {
			httpx.Error(w, errs)
			return
		}

		l := logic_.NewUpdateMerchantPwdLogic(r.Context(), ctx)
		err := l.UpdateMerchantPwd(userId, req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, nil)
		}
	}
}
