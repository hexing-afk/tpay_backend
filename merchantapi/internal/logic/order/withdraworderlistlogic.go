package order

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

const (
	WithdrawOrderStatusPending    = 1 // 待处理
	WithdrawOrderStatusProcessing = 2 // 处理中
	WithdrawOrderStatusFail       = 3 // 失败
	WithdrawOrderStatusSuccess    = 4 // 成功
)

type WithdrawOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWithdrawOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) WithdrawOrderListLogic {
	return WithdrawOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WithdrawOrderListLogic) WithdrawOrderList(merchantId int64, req types.WithdrawOrderListRequest) (*types.WithdrawOrderListResponse, error) {
	orderStatuses := l.withdrawOrderStatusInnerChange(req.OrderStatus)

	f := model.FindWithdrawOrderList{
		Page:            req.Page,
		PageSize:        req.PageSize,
		StartCreateTime: req.StartCreateTime,
		EndCreateTime:   req.EndCreateTime,
		MerchantId:      merchantId,
		OrderNo:         req.OrderNo,
		OrderStatuses:   orderStatuses,
	}
	data, total, err := model.NewMerchantWithdrawOrderModel(l.svcCtx.DbEngine).FindList(f)
	if err != nil {
		l.Errorf("查询商户[%v]提现订单列表失败, err=%v", merchantId, err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var list []types.WithdrawOrderList
	for _, v := range data {
		list = append(list, types.WithdrawOrderList{
			OrderNo:     v.OrderNo,
			OrderStatus: l.withdrawOrderStatusOutChange(v.OrderStatus),
			OrderAmount: v.OrderAmount,
			Fee:         v.MerchantFee,
			RealAmount:  v.RealAmount,
			BankName:    v.BankName,
			PayeeNmae:   v.PayeeName,
			CardNumber:  v.CardNumber,
			BranchName:  v.BranchName,
			AuditRemark: v.AuditRemark,
			AuditTime:   v.AuditTime,
			CreateTime:  v.CreateTime,
			Currency:    v.Currency,
			Remark:      v.Remark,
		})
	}

	return &types.WithdrawOrderListResponse{
		Total: total,
		List:  list,
	}, nil
}

// 前端订单状态转数据库订单状态
func (l *WithdrawOrderListLogic) withdrawOrderStatusInnerChange(status int64) []int64 {
	switch status {
	case WithdrawOrderStatusPending: // 待处理
		// 待处理
		return []int64{model.WithdrawOrderStatusPending}
	case WithdrawOrderStatusProcessing: // 处理中
		// 通过审核，派单中
		return []int64{model.WithdrawOrderStatusPass, model.WithdrawOrderStatusAllot}
	case WithdrawOrderStatusFail: // 失败
		// 驳回， 派单失败
		return []int64{model.WithdrawOrderStatusReject, model.WithdrawOrderStatusAllotFail}
	case WithdrawOrderStatusSuccess: // 成功
		// 派单成功， 成功
		return []int64{model.WithdrawOrderStatusAllotSuccess, model.WithdrawOrderStatusSuccess}
	default:
		return nil
	}
}

// 数据库订单状态转前端订单状态
func (l *WithdrawOrderListLogic) withdrawOrderStatusOutChange(status int64) int64 {
	switch status {
	case model.WithdrawOrderStatusPending: // 待处理
		// 待处理
		return WithdrawOrderStatusPending
	case model.WithdrawOrderStatusPass: // 通过审核
		// 处理中
		return WithdrawOrderStatusProcessing
	case model.WithdrawOrderStatusAllot: // 派单中
		// 处理中
		return WithdrawOrderStatusProcessing
	case model.WithdrawOrderStatusReject: // 驳回
		// 失败
		return WithdrawOrderStatusFail
	case model.WithdrawOrderStatusAllotFail: // 派单失败
		// 失败
		return WithdrawOrderStatusFail
	case model.WithdrawOrderStatusAllotSuccess: // 派单成功
		// 成功
		return WithdrawOrderStatusSuccess
	case model.WithdrawOrderStatusSuccess: // 成功
		// 成功
		return WithdrawOrderStatusSuccess
	default:
		return 0
	}
}
