package order

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type PayOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) PayOrderListLogic {
	return PayOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PayOrderListLogic) PayOrderList(merchantId int64, req types.PayOrderListRequest) (*types.PayOrderListResponse, error) {
	f := model.FindPayOrderList{
		Page:                req.Page,
		PageSize:            req.PageSize,
		OrderNo:             req.OrderNo,
		MerchantOrderNo:     req.MerchantOrderNo,
		MerchantId:          merchantId,
		PlatformChannelId:   req.ChannelId,
		PlatformChannelName: req.ChannelName,
		StartCreateTime:     req.StartCreateTime,
		EndCreateTime:       req.EndCreateTime,
		OrderStatus:         req.OrderStatus,
		OrderType:           req.OrderType,
	}
	data, total, err := model.NewPayOrderModel(l.svcCtx.DbEngine).FindList(f)
	if err != nil {
		l.Errorf("查询商户收款订单失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var list []types.PayOrderList
	for _, v := range data {
		list = append(list, types.PayOrderList{
			OrderNo:         v.OrderNo,
			MerchantOrderNo: v.MerchantOrderNo,
			OrderAmount:     v.ReqAmount,
			Fee:             v.MerchantFee,
			IncreaseAmount:  v.IncreaseAmount,
			ChannelName:     v.PlatformChannelName,
			OrderStatus:     v.OrderStatus,
			CreateTime:      v.CreateTime,
			UpdateTime:      v.UpdateTime,
			PaymentAmount:   v.PaymentAmount,
		})
	}

	return &types.PayOrderListResponse{
		Total: total,
		List:  list,
	}, nil
}
