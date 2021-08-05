package logic

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"
	"tpay_backend/utils"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type PayOrderExportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayOrderExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) PayOrderExportLogic {
	return PayOrderExportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PayOrderExportLogic) PayOrderExport(req types.PayOrderExportRequest) (*types.PayOrderExportResponse, error) {
	day := utils.TimeSubToDay(req.StartCreateTime, req.EndCreateTime)
	if day < 0 || day > 31 {
		l.Errorf("导出时间超出限制(31天), start:%v, end:%v", req.StartCreateTime, req.EndCreateTime)
		return nil, common.NewCodeError(common.ExportLimit31Day)
	}

	currency, err := model.NewCurrencyModel(l.svcCtx.DbEngine).FindByCurrency(req.Currency)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("币种不存在, err=%v", err)
			return nil, common.NewCodeError(common.InvalidParam)
		} else {
			l.Errorf("查询币种失败, err=%v", err)
			return nil, common.NewCodeError(common.SystemInternalErr)
		}
	}

	f := model.FindExportData{
		OrderNo:           req.OrderNo,
		MerchantOrderNo:   req.MerchantOrderNo,
		UpstreamOrderNo:   req.UpstreamOrderNo,
		MerchantName:      req.MerchantName,
		PlatformChannelId: req.ChannelId,
		StartCreateTime:   req.StartCreateTime,
		EndCreateTime:     req.EndCreateTime,
		OrderStatus:       req.OrderStatus,
		Currency:          req.Currency,
		OrderType:         model.PayModePro,
	}

	data, err := model.NewPayOrderModel(l.svcCtx.DbEngine).FindExportData(f)
	if err != nil {
		l.Errorf("查询代收订单数据失败, err=%v", err)
		return nil, common.NewCodeError(common.SystemInternalErr)
	}

	if data == nil || data.Total == 0 {
		l.Errorf("没有数据")
		return nil, common.NewCodeError(common.NotData)
	}

	var list []types.PayOrderExportData
	for _, v := range data.OrderList {
		var status string
		switch v.OrderStatus {
		case model.PayOrderStatusPending:
			status = "待支付"
		case model.PayOrderStatusPaid:
			status = "已支付"
		case model.PayOrderStatusFail:
			status = "支付失败"
		}
		list = append(list, types.PayOrderExportData{
			Id:              v.Id,
			MerchantName:    v.MerchantName,
			OrderNo:         v.OrderNo,
			MerchantOrderNo: v.MerchantOrderNo,
			ReqAmount:       v.ReqAmount,
			PaymentAmount:   v.PaymentAmount,
			Rate:            v.MerchantRate,
			SingleFee:       v.MerchantSingleFee,
			Fee:             v.MerchantFee,
			IncreaseAmount:  v.IncreaseAmount,
			ChannelName:     v.PlatformChannelName,
			OrderStatus:     status,
			CreateTime:      v.CreateTime,
			UpdateTime:      v.UpdateTime,
		})
	}

	return &types.PayOrderExportResponse{
		Total:               data.Total,
		TotalReqAmount:      data.TotalPayAmount,
		TotalPayAmount:      data.TotalPayAmount,
		TotalFee:            data.TotalMerchantFee,
		TotalIncreaseAmount: data.TotalIncreaseAmount,
		IsDivideHundred:     currency.IsDivideHundred == model.DivideHundred,
		List:                list,
	}, nil
}
