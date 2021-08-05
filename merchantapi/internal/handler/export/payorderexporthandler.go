package handler

import (
	"bytes"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"tpay_backend/export"
	"tpay_backend/merchantapi/internal/common"
	_func "tpay_backend/merchantapi/internal/handler/func"
	"tpay_backend/merchantapi/internal/logic/export"
	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/rest/httpx"
)

func PayOrderExportHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PayOrderExportRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, common.NewCodeErrorWithMsg(common.VerifyParamFailed, err.Error()))
			return
		}

		if req.StartCreateTime == 0 && req.EndCreateTime == 0 {
			httpx.Error(w, common.NewCodeError(common.CheckExportTime))
			return
		}

		userId, errs := _func.GetLoginedUserIdRequestHeader(r)
		if errs != nil {
			httpx.Error(w, errs)
			return
		}

		l := logic.NewPayOrderExportLogic(r.Context(), ctx)
		resp, err := l.PayOrderExport(userId, req)
		if err != nil {
			httpx.Error(w, err)
		} else {

			if len(resp.List) <= 0 || resp.Total == 0 {
				httpx.Error(w, common.NewCodeError(common.NotData))
				return
			}

			var orderList []export.PayOrder
			for _, v := range resp.List {
				orderList = append(orderList, export.PayOrder{
					Id:              v.Id,
					MerchantName:    v.MerchantName,
					OrderNo:         v.OrderNo,
					MerchantOrderNo: v.MerchantOrderNo,
					ReqAmount:       v.ReqAmount,
					PaymentAmount:   v.PaymentAmount,
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
			file, err = export.CreatePayOrderFile(file, &export.CreatePayOrderFileRequest{
				Sheet:           "Sheet1",
				Title:           "代收订单导出表",
				Timezone:        ctx.Config.Timezone,
				IsDivideHundred: resp.IsDivideHundred,
				Total: &export.PayOrderTotal{
					Total:               resp.Total,
					TotalReqAmount:      resp.TotalReqAmount,
					TotalPayAmount:      resp.TotalPayAmount,
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
				l.Errorf("文件读取失败, err=%v", err)
				httpx.Error(w, common.NewCodeError(common.ExportFail))
				return
			}

			fileName := "代收订单导出表-" + time.Now().Format("2006/01/02") + ".xlsx"
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
