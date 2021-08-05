package order

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetMerchantRechargeOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMerchantRechargeOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetMerchantRechargeOrderListLogic {
	return GetMerchantRechargeOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMerchantRechargeOrderListLogic) GetMerchantRechargeOrderList(req types.GetMerchantRechargeOrderListRequest) (*types.GetMerchantRechargeOrderListResponse, error) {
	f := model.FindRechargeOrderList{
		Page:            req.Page,
		PageSize:        req.PageSize,
		MerchantName:    req.MerchantName,
		OrderNo:         req.OrderNo,
		OrderStatus:     req.OrderStatus,
		StartCreateTime: req.StartCreateTime,
		EndCreateTime:   req.EndCreateTime,
	}
	data, total, err := model.NewMerchantRechargeOrderModel(l.svcCtx.DbEngine).FindList(f)
	if err != nil {
		l.Errorf("查询商户充值订单列表失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var list []types.MerchantRechargeOrderList
	for _, v := range data {
		list = append(list, types.MerchantRechargeOrderList{
			OrderNo:        v.OrderNo,
			MerchantName:   v.MerchantName,
			OrderAmount:    v.OrderAmount,
			RechargeRemark: v.RechargeRemark,
			BankName:       v.BankName,
			PayeeName:      v.PayeeName,
			CardNumber:     v.CardNumber,
			OrderStatus:    v.OrderStatus,
			DealWithRemark: v.AuditRemark,
			CreateTime:     v.CreateTime,
			FinishTime:     v.FinishTime,
			Currency:       v.Currency,
		})
	}

	return &types.GetMerchantRechargeOrderListResponse{
		Total: total,
		List:  list,
	}, nil
}
