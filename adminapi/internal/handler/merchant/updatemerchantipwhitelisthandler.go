package merchant

import (
	"net/http"
	"tpay_backend/adminapi/internal/common"

	logic_ "tpay_backend/adminapi/internal/logic/merchant"
	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func UpdateMerchantIpWhitelistHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateMerchantIpWhitelistRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic_.NewUpdateMerchantIpWhitelistLogic(r.Context(), ctx)
		err := l.UpdateMerchantIpWhitelist(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, nil)
		}
	}
}
