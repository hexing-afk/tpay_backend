package order

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetPayOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPayOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPayOrderListLogic {
	return GetPayOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPayOrderListLogic) GetPayOrderList(req types.GetPayOrderListRequest) (*types.GetPayOrderListResponse, error) {
	f := model.FindPayOrderList{
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
	}
	data, total, err := model.NewPayOrderModel(l.svcCtx.DbEngine).FindList(f)
	if err != nil {
		l.Errorf("查询订单列表失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var list []types.PayOrderList
	for _, v := range data {
		list = append(list, types.PayOrderList{
			OrderNo:         v.OrderNo,
			MerchantOrderNo: v.MerchantOrderNo,
			UpstreamOrderNo: v.UpstreamOrderNo,
			MerchantName:    v.MerchantName,
			Currency:        v.Currency,
			OrderAmount:     v.ReqAmount,
			MerchantFee:     v.MerchantFee,
			IncreaseAmount:  v.IncreaseAmount,
			UpstreamName:    v.UpstreamName,
			ChannelName:     v.PlatformChannelName,
			OrderStatus:     v.OrderStatus,
			CreateTime:      v.CreateTime,
			UpdateTime:      v.UpdateTime,
			PaymentAmount:   v.PaymentAmount,
		})
	}

	return &types.GetPayOrderListResponse{
		Total: total,
		List:  list,
	}, nil
}
