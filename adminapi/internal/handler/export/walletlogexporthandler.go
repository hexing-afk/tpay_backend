package handler

import (
	"bytes"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/export"

	"tpay_backend/adminapi/internal/logic/export"
	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func WalletLogExportHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WalletLogExportRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, common.NewCodeErrorWithMsg(common.VerifyParamFailed, err.Error()))
			return
		}

		if req.Currency == "" {
			httpx.Error(w, common.NewCodeError(common.MissingParam))
			return
		}

		if req.StartCreateTime == 0 && req.EndCreateTime == 0 {
			httpx.Error(w, common.NewCodeError(common.CheckExportTime))
			return
		}

		l := logic.NewWalletLogExportLogic(r.Context(), ctx)
		resp, err := l.WalletLogExport(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			if resp.Total == 0 || len(resp.List) <= 0 {
				httpx.Error(w, common.NewCodeError(common.NotData))
				return
			}

			var logList []export.WalletLog
			for _, v := range resp.List {
				logList = append(logList, export.WalletLog{
					Id:           v.Id,
					BusinessNo:   v.BusinessNo,
					ChangeAmount: v.ChangeAmount,
					AfterBalance: v.AfterBalance,
					OpType:       v.OpType,
					OrderType:    v.OrderType,
					Remark:       v.Remark,
					CreateTime:   v.CreateTime,
				})
			}

			file := xlsx.NewFile()
			file, err = export.CreateWalletLogFile(file, &export.CreateWalletLogFileRequest{
				Sheet:           "Sheet1",
				Title:           "钱包流水导出表",
				Timezone:        ctx.Config.Timezone,
				IsDivideHundred: resp.IsDivideHundred,
				Content:         logList,
			})
			if err != nil {
				l.Errorf("创建文件失败, err=%v", err)
				httpx.Error(w, common.NewCodeError(common.ExportFail))
				return
			}

			var buffer bytes.Buffer
			if err := file.Write(&buffer); err != nil {
				l.Errorf("文件转换失败, err=%v", err)
				httpx.Error(w, common.NewCodeError(common.ExportFail))
				return
			}
			r := bytes.NewReader(buffer.Bytes())

			fSrc, err := ioutil.ReadAll(r)
			if err != nil {
				l.Errorf("文件内容读取失败, err=%v", err)
				httpx.Error(w, common.NewCodeError(common.ExportFail))
				return
			}

			fileName := "钱包流水导出表" + time.Now().Format("2006/01/02") + ".xlsx"
			// 防止中文乱码
			fileName = url.QueryEscape(fileName)
			//w.Header().Add("Content-Type", "application/octet-stream")
			//w.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
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
