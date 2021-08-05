package order

import (
	"net/http"
	"tpay_backend/adminapi/internal/common"

	logic_ "tpay_backend/adminapi/internal/logic/order"
	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func GetPayOrderDetailHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetPayOrderDetailRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic_.NewGetPayOrderDetailLogic(r.Context(), ctx)
		resp, err := l.GetPayOrderDetail(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, resp)
		}
	}
}
