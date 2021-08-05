package logic

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"
	"tpay_backend/utils"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type TransferOrderExportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransferOrderExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) TransferOrderExportLogic {
	return TransferOrderExportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransferOrderExportLogic) TransferOrderExport(merchantId int64, req types.TransferOrderExportRequest) (*types.TransferOrderExportResponse, error) {
	day := utils.TimeSubToDay(req.StartCreateTime, req.EndCreateTime)
	if day < 0 || day > 31 {
		l.Errorf("导出时间超出限制(31天), start:%v, end:%v", req.StartCreateTime, req.EndCreateTime)
		return nil, common.NewCodeError(common.ExportLimit31Day)
	}

	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneById(merchantId)
	if err != nil {
		l.Errorf("查询商户失败, err=%v", err)
		return nil, common.NewCodeError(common.ExportFail)
	}

	f := model.FindTransferExportData{
		MerchantId:      merchantId,
		OrderNo:         req.OrderNo,
		MerchantOrderNo: req.MerchantOrderNo,
		StartCreateTime: req.StartCreateTime,
		EndCreateTime:   req.EndCreateTime,
		OrderStatus:     req.OrderStatus,
		OrderType:       model.TransferModePro,
		OrderSourceList: []int64{model.TransferOrderSourceInterface},
	}

	data, err := model.NewTransferOrderModel(l.svcCtx.DbEngine).FindExportData(f)
	if err != nil {
		l.Errorf("查询代付订单数据失败, err=%v", err)
		return nil, common.NewCodeError(common.ExportFail)
	}

	if data == nil || data.Total == 0 {
		l.Errorf("没有数据")
		return nil, common.NewCodeError(common.NotData)
	}

	currency, err := model.NewCurrencyModel(l.svcCtx.DbEngine).FindByCurrency(merchant.Currency)
	if err != nil {
		l.Errorf("查询币种失败, err=%v", err)
		return nil, common.NewCodeError(common.ExportFail)
	}

	var list []types.TransferOrderExportData
	for _, v := range data.OrderList {
		var status string
		switch v.OrderStatus {
		case model.TransferOrderStatusPending:
			status = "待支付"
		case model.TransferOrderStatusPaid:
			status = "已支付"
		case model.TransferOrderStatusFail:
			status = "支付失败"
		}

		list = append(list, types.TransferOrderExportData{
			Id:              v.Id,
			MerchantName:    v.MerchantName,
			OrderNo:         v.OrderNo,
			MerchantOrderNo: v.MerchantOrderNo,
			ReqAmount:       v.ReqAmount,
			Rate:            v.MerchantRate,
			SingleFee:       v.MerchantSingleFee,
			Fee:             v.MerchantFee,
			IncreaseAmount:  v.PayeeRealAmount,
			ChannelName:     v.PlatformChannelName,
			OrderStatus:     status,
			CreateTime:      v.CreateTime,
			UpdateTime:      v.UpdateTime,
		})
	}

	return &types.TransferOrderExportResponse{
		Total:               data.Total,
		TotalReqAmount:      data.TotalReqAmount,
		TotalFee:            data.TotalFee,
		TotalIncreaseAmount: data.TotalIncreaseAmount,
		IsDivideHundred:     currency.IsDivideHundred == model.DivideHundred,
		List:                list,
	}, nil
}
