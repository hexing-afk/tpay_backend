package public

import (
	"net/http"
	"tpay_backend/merchantapi/internal/common"
	_func "tpay_backend/merchantapi/internal/handler/func"
	logic "tpay_backend/merchantapi/internal/logic/public"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
	"tpay_backend/merchantapi/internal/svc"
)

func GetChannelNameListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetChannelNameListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		userId, errs := _func.GetLoginedUserIdRequestHeader(r)
		if errs != nil {
			httpx.Error(w, errs)
			return
		}

		l := logic.NewGetChannelNameListLogic(r.Context(), ctx)
		resp, err := l.GetChannelNameList(userId, req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, resp)
		}
	}
}
