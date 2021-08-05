package merchant

import (
	"net/http"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/adminapi/internal/logic/merchant"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func GetMerchantBankCardListHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetMerchantBankCardListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, common.NewCodeErrorWithMsg(common.VerifyParamFailed, err.Error()))
			return
		}

		l := merchant.NewGetMerchantBankCardListLogic(r.Context(), ctx)
		resp, err := l.GetMerchantBankCardList(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, resp)
		}
	}
}
