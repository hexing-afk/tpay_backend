package public

import (
	"net/http"
	"tpay_backend/merchantapi/internal/common"
	_func "tpay_backend/merchantapi/internal/handler/func"

	_logic "tpay_backend/merchantapi/internal/logic/public"
	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func UpdateBaseInfoHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateBaseInfoReq
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

		l := _logic.NewUpdateBaseInfoLogic(r.Context(), ctx, userId)
		err := l.UpdateBaseInfo(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, nil)
		}
	}
}