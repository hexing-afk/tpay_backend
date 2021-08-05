package order

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetMerchantWithdrawOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMerchantWithdrawOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetMerchantWithdrawOrderListLogic {
	return GetMerchantWithdrawOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMerchantWithdrawOrderListLogic) GetMerchantWithdrawOrderList(req types.GetMerchantWithdrawOrderListRequest) (*types.GetMerchantWithdrawOrderListResponse, error) {
	f := model.FindWithdrawOrderList{
		Page:            req.Page,
		PageSize:        req.PageSize,
		StartCreateTime: req.StartCreateTime,
		EndCreateTime:   req.EndCreateTime,
		MerchantName:    req.MerchantName,
		OrderNo:         req.OrderNo,
		OrderStatus:     req.OrderStatus,
	}
	data, total, err := model.NewMerchantWithdrawOrderModel(l.svcCtx.DbEngine).FindList(f)
	if err != nil {
		l.Errorf("查询商户提现订单列表失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var list []types.MerchantWithdrawOrderList
	for _, v := range data {
		list = append(list, types.MerchantWithdrawOrderList{
			OrderNo:             v.OrderNo,
			MerchantName:        v.MerchantName,
			OrderAmount:         v.OrderAmount,
			MerchantFee:         v.MerchantFee,
			RealAmount:          v.RealAmount,
			Remark:              v.Remark,
			BankName:            v.BankName,
			PayeeName:           v.PayeeName,
			CardNumber:          v.CardNumber,
			BranchName:          v.BranchName,
			AuditRemark:         v.AuditRemark,
			CreateTime:          v.CreateTime,
			AuditTime:           v.AuditTime,
			OrderStatus:         v.OrderStatus,
			UpstreamOrderNo:     v.TransferOrderNo,
			UpstreamChannelName: v.ChannelName,
			Currency:            v.Currency,
			DeductionMethod:     v.DeductionMethod,
		})
	}

	return &types.GetMerchantWithdrawOrderListResponse{
		Total: total,
		List:  list,
	}, nil
}
