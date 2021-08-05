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

func TransferOrderExportHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TransferOrderExportRequest
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

		l := logic.NewTransferOrderExportLogic(r.Context(), ctx)
		resp, err := l.TransferOrderExport(req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			if resp.Total == 0 || len(resp.List) <= 0 {
				httpx.Error(w, common.NewCodeError(common.NotData))
				return
			}

			var orderList []export.TransferOrder
			for _, v := range resp.List {
				orderList = append(orderList, export.TransferOrder{
					Id:              v.Id,
					MerchantName:    v.MerchantName,
					OrderNo:         v.OrderNo,
					MerchantOrderNo: v.MerchantOrderNo,
					ReqAmount:       v.ReqAmount,
					Rate:            v.Rate,
					SingleFee:       v.SingleFee,
					Fee:             v.Fee,
					IncreaseAmount:  v.IncreaseAmount,
					ChannelName:     v.ChannelName,
					OrderStatus:     v.OrderStatus,
					CreateTime:      v.CreateTime,
					UpdateTime:      v.UpdateTime,
				})
			}

			file := xlsx.NewFile()
			file, err = export.CreateTransferOrderFile(file, &export.CreateTransferOrderFileRequest{
				Sheet:           "Sheet1",
				Title:           "代付订单导出表",
				Timezone:        ctx.Config.Timezone,
				IsDivideHundred: resp.IsDivideHundred,
				Total: &export.TransferOrderTotal{
					Total:               resp.Total,
					TotalReqAmount:      resp.TotalReqAmount,
					TotalFee:            resp.TotalFee,
					TotalIncreaseAmount: resp.TotalIncreaseAmount,
				},
				Content: orderList,
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

			fileName := "代付订单导出表-" + time.Now().Format("2006/01/02") + ".xlsx"
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

			//common.OkJson(w, resp)
			httpx.Ok(w)
		}
	}
}