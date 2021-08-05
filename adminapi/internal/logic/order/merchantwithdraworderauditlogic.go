package order

import (
	"context"
	"gorm.io/gorm"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type MerchantWithdrawOrderAuditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantWithdrawOrderAuditLogic(ctx context.Context, svcCtx *svc.ServiceContext) MerchantWithdrawOrderAuditLogic {
	return MerchantWithdrawOrderAuditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantWithdrawOrderAuditLogic) MerchantWithdrawOrderAudit(req types.MerchantWithdrawOrderAuditRequest) error {
	order, err := model.NewMerchantWithdrawOrderModel(l.svcCtx.DbEngine).FindByOrderNo(req.OrderNo)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("商户提现订单[%v]不存在", req.OrderNo)
			return common.NewCodeError(common.OrderNotExist)
		} else {
			l.Errorf("查询商户提现订单[%v]失败, err=%v", req.OrderNo, err)
			return common.NewCodeError(common.SysDBSave)
		}
	}

	if order.OrderStatus != model.WithdrawOrderStatusPending {
		l.Errorf("商户提现订单[%v]已处理，audit_status=%v ", req.OrderNo, order.OrderStatus)
		return common.NewCodeError(common.OrderAlreadyProcessed)
	}

	switch req.AuditStatus {
	case model.WithdrawOrderStatusPass:
		// 修改审核状态
		err := model.NewMerchantWithdrawOrderModel(l.svcCtx.DbEngine).UpdateStatusToPass(order.Id, req.AuditRemark)
		if err != nil {
			l.Errorf("通过商户提现订单[%v]失败, err=%v", order.OrderNo, err)
			return common.NewCodeError(common.SysDBSave)
		}
	case model.WithdrawOrderStatusReject:
		txErr := l.svcCtx.DbEngine.Transaction(func(tx *gorm.DB) error {
			// 修改审核状态
			if err := model.NewMerchantWithdrawOrderModel(tx).UpdateStatusToReject(order.Id, req.AuditRemark); err != nil {
				l.Errorf("修改商户提现订单[%v]审核状态失败, err=%v", order.OrderNo, err)
				return err
			}

			// 驳回提现；加余额，减冻结金额
			walletLog := model.WalletLogExt{
				BusinessNo: order.OrderNo,
				Source:     model.AmountSourceWithdraw,
				Remark:     "提现订单-驳回",
			}
			if err := model.NewMerchantModel(tx).PlusBalanceUnfreezeTx(order.MerchantId, order.DecreaseAmount, walletLog); err != nil {
				l.Errorf("增加商户[%v]余额[%v]，减少冻结金额失败, err=%v", order.MerchantId, order.DecreaseAmount, err)
				return err
			}

			return nil
		})

		if txErr != nil {
			l.Errorf("驳回商户提现订单[%v]失败, err=%v", order.OrderNo, err)
			return common.NewCodeError(common.SysDBSave)
		}
	default:
		l.Errorf("不支持当前操作, req.AuditStatus=%v", req.AuditStatus)
		return common.NewCodeError(common.InvalidParam)
	}

	return nil
}
