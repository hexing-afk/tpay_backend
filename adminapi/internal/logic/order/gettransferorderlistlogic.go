package order

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetTransferOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTransferOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetTransferOrderListLogic {
	return GetTransferOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTransferOrderListLogic) GetTransferOrderList(req types.GetTransferOrderListRequest) (*types.GetTransferOrderListResponse, error) {
	f := model.FindTransferOrderList{
		Page:              req.Page,
		PageSize:          req.PageSize,
		OrderNo:           req.OrderNo,
		MerchantOrderNo:   req.MerchantOrderNo,
		UpstreamOrderNo:   req.UpstreamOrderNo,
		MerchantName:      req.MerchantName,
		Currency:          req.Currency,
		PlatformChannelId: req.ChannelId,
		StartCreateTime:   req.StartCreateTime,
		EndCreateTime:     req.EndCreateTime,
		OrderStatus:       req.OrderStatus,
		OrderType:         req.OrderType,
		//OrderSource:       model.TransferOrderSourceInterface,
	}
	data, total, err := model.NewTransferOrderModel(l.svcCtx.DbEngine).FindList(f)
	if err != nil {
		l.Errorf("查询代付订单失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var list []types.TransferOrderList
	for _, v := range data {
		list = append(list, types.TransferOrderList{
			OrderNo:         v.OrderNo,
			MerchantOrderNo: v.MerchantOrderNo,
			UpstreamOrderNo: v.UpstreamOrderNo,
			MerchantName:    v.MerchantName,
			Currency:        v.Currency,
			OrderAmount:     v.ReqAmount,
			MerchantFee:     v.MerchantFee,
			PayeeRealAmount: v.PayeeRealAmount,
			UpstreamName:    v.UpstreamName,
			ChannelName:     v.PlatformChannelName,
			OrderStatus:     v.OrderStatus,
			CreateTime:      v.CreateTime,
			UpdateTime:      v.UpdateTime,
			OrderSource:     v.OrderSource,
		})
	}

	return &types.GetTransferOrderListResponse{
		Total: total,
		List:  list,
	}, nil
}
