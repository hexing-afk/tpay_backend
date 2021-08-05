package login

import (
	"bytes"
	"github.com/tal-tech/go-zero/rest/httpx"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"tpay_backend/export"
	"tpay_backend/merchantapi/internal/common"
	logic "tpay_backend/merchantapi/internal/logic/login"
	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"
)

func PayOrderExportTHandler(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PayOrderExportRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		var userId int64 = 61

		l := logic.NewPayOrderExportTLogic(r.Context(), ctx)
		resp, err := l.PayOrderExportT(userId, req)
		if err != nil {
			httpx.Error(w, err)
		} else {

			if len(resp.List) <= 0 {
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
				httpx.Error(w, common.NewCodeError(common.SystemInternalErr))
				return
			}

			var buffer bytes.Buffer
			if err := file.Write(&buffer); err != nil {
				l.Errorf("文件转换失败, err=%v", err)
				httpx.Error(w, common.NewCodeError(common.SystemInternalErr))
				return
			}
			r := bytes.NewReader(buffer.Bytes())

			fSrc, err := ioutil.ReadAll(r)
			if err != nil {
				l.Errorf("文件内容读取失败, err=%v", err)
				httpx.Error(w, common.NewCodeError(common.SystemInternalErr))
				return
			}

			fileName := "代收订单导出表" + time.Now().Format("2006/01/02") + ".xlsx"
			// 防止中文乱码
			fileName = url.QueryEscape(fileName)
			//w.Header().Add("Content-Type", "application/octet-stream")
			//w.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
			w.Header().Add("content-disposition", "attachment; filename=\""+fileName+"\"")

			if _, err := w.Write(fSrc); err != nil {
				l.Errorf("文件写入response失败, err=%v", err)
				httpx.Error(w, common.NewCodeError(common.SystemInternalErr))
				return
			}

			//common.OkJson(w, fSrc)
			httpx.Ok(w)
		}
	}
}
