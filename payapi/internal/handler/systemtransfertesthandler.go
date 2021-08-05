package handler

import (
	"encoding/json"
	"net/http"
	"tpay_backend/payapi/internal/common"
	"tpay_backend/utils"

	"github.com/tal-tech/go-zero/core/logx"

	"tpay_backend/payapi/internal/logic"
	"tpay_backend/payapi/internal/svc"
	"tpay_backend/payapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func SystemTransferTestHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SystemTransferTestReq

		// 请求公共数据及签名验证
		merchant, err := CheckInternalRequestData(ctx, r, &req)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewSystemTransferTestLogic(r.Context(), ctx, merchant)
		resp, err := l.SystemTransferTest(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			jsonStr, jerr := json.Marshal(resp)
			if jerr != nil {
				logx.Errorf("返回参数json编码失败:%v,resp:%+v", jerr, resp)
			}
			common.OkJson(w, string(jsonStr), utils.GenerateSign(string(jsonStr), merchant.Md5Key))
		}
	}
}
