package public

import (
	"net/http"
	"net/url"
	"tpay_backend/merchantapi/internal/common"

	"github.com/tal-tech/go-zero/rest/httpx"
	logic_ "tpay_backend/merchantapi/internal/logic/public"
	"tpay_backend/merchantapi/internal/svc"
)

func BatchTransferFileHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic_.NewBatchTransferFileLogic(r.Context(), ctx)
		fSrc, err := l.BatchTransferFile()
		if err != nil {
			httpx.Error(w, err)
		} else {

			fileName := "批量转账模板.xlsx"
			// 防止中文乱码
			fileName = url.QueryEscape(fileName)
			w.Header().Add("Content-Disposition", "attachment; filename=\""+fileName+"\"")

			if _, err := w.Write(fSrc); err != nil {
				l.Errorf("文件写入response失败, err=%v", err)
				httpx.Error(w, common.NewCodeError(common.ExportFail))
				return
			}

			httpx.Ok(w)
		}

	}
}
