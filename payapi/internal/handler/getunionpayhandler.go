package handler

import (
	"net/http"

	"tpay_backend/payapi/internal/logic"
	"tpay_backend/payapi/internal/svc"
	"tpay_backend/payapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func GetUnionpayHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUnionpayReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewGetUnionpayLogic(r.Context(), ctx)
		resp, err := l.GetUnionpay(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
