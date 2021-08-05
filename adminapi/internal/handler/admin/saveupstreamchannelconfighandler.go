package admin

import (
	"errors"
	"net/http"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	logic "tpay_backend/adminapi/internal/logic/admin"
	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func SaveUpstreamChannelConfigHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SaveUpstreamChannelConfigRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, common.NewCodeErrorWithMsg(common.VerifyParamFailed, err.Error()))
			return
		}

		if req.Rate < 0 || req.Rate > 100 {
			err := errors.New("Incorrect value of parameter 'rate' ")
			httpx.Error(w, common.NewCodeErrorWithMsg(common.VerifyParamFailed, err.Error()))
			return
		}

		if req.DeductionMethod != model.UpstreamChannelDeductionInner && req.DeductionMethod != model.UpstreamChannelDeductionOut {
			httpx.Error(w, common.NewCodeError(common.VerifyParamFailed))
			return
		}

		l := logic.NewSaveUpstreamChannelConfigLogic(r.Context(), ctx)
		err := l.SaveUpstreamChannelConfig(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, nil)
		}
	}
}
