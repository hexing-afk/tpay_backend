package admin

import (
	"net/http"
	"tpay_backend/adminapi/internal/common"
	_func "tpay_backend/adminapi/internal/handler/func"
	logic "tpay_backend/adminapi/internal/logic/admin"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func ResetAdminPwdHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ResetAdminPwdRequest
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

		l := logic.NewResetAdminPwdLogic(r.Context(), ctx)
		err = l.ResetAdminPwd(userId, req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, nil)
		}
	}
}
