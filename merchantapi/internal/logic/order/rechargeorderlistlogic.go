package order

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type RechargeOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRechargeOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) RechargeOrderListLogic {
	return RechargeOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RechargeOrderListLogic) RechargeOrderList(merchantId int64, req types.RechargeOrderListRequest) (*types.RechargeOrderListResponse, error) {
	f := model.FindRechargeOrderList{
		Page:            req.Page,
		PageSize:        req.PageSize,
		MerchantId:      merchantId,
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

	var list []types.RechargeOrderList
	for _, v := range data {
		list = append(list, types.RechargeOrderList{
			OrderNo:     v.OrderNo,
			OrderAmount: v.OrderAmount,
			BankName:    v.BankName,
			PayeeName:   v.PayeeName,
			CardNumber:  v.CardNumber,
			BranchName:  v.BranchName,
			OrderStatus: v.OrderStatus,
			AuditRemark: v.AuditRemark,
			CreateTime:  v.CreateTime,
			UpdateTime:  v.FinishTime,
			Currency:    v.Currency,
			Remark:      v.RechargeRemark,
		})
	}

	return &types.RechargeOrderListResponse{
		Total: total,
		List:  list,
	}, nil
}
