package public

import (
	"github.com/tal-tech/go-zero/core/logx"
	"net/http"
	"tpay_backend/merchantapi/internal/common"
	_func "tpay_backend/merchantapi/internal/handler/func"

	"github.com/tal-tech/go-zero/rest/httpx"
	logic_ "tpay_backend/merchantapi/internal/logic/public"
	"tpay_backend/merchantapi/internal/svc"
)

func UploadFileHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//得到上传的文件
		_, fileHeader, err := r.FormFile("file")
		if err != nil {
			logx.Errorf("获取文件信息出错，err=[%v]", err)
			httpx.Error(w, common.NewCodeError(common.UploadFail))
			return
		}

		userId, errs := _func.GetLoginedUserIdRequestHeader(r)
		if errs != nil {
			httpx.Error(w, errs)
			return
		}

		l := logic_.NewUploadFileLogic(r.Context(), ctx)
		resp, err := l.UploadFile(userId, fileHeader)
		if err != nil {
			httpx.Error(w, err)
		} else {
			common.OkJson(w, resp)
		}
	}
}
