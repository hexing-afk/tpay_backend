package order

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type MerchantRechargeOrderAuditLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantRechargeOrderAuditLogic(ctx context.Context, svcCtx *svc.ServiceContext) MerchantRechargeOrderAuditLogic {
	return MerchantRechargeOrderAuditLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantRechargeOrderAuditLogic) MerchantRechargeOrderAudit(req types.MerchantRechargeOrderAuditRequest) error {
	order, err := model.NewMerchantRechargeOrderModel(l.svcCtx.DbEngine).FindByOrderNo(req.OrderNo)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("商户充值订单[%v]不存在", req.OrderNo)
			return common.NewCodeError(common.OrderNotExist)
		} else {
			l.Errorf("查询商户充值订单[%v]失败, err=%v", req.OrderNo, err)
			return common.NewCodeError(common.SysDBSave)
		}
	}

	if order.OrderStatus != model.RechargeOrderStatusPending {
		l.Errorf("商户充值订单[%v]已经处理了, order_status=%v", req.OrderNo, order.OrderStatus)
		return common.NewCodeError(common.OrderAlreadyProcessed)
	}
	switch req.OrderStatus {
	case model.RechargeOrderStatusPass:
		txErr := l.svcCtx.DbEngine.Transaction(func(tx *gorm.DB) error {
			// 1.修改状态
			if err := model.NewMerchantRechargeOrderModel(tx).UpdateStatusToPass(req.OrderNo, req.DealWithRemark); err != nil {
				l.Errorf("修改商户充值订单[%v]状态为[%v]失败, err=%v", req.OrderNo, req.OrderStatus, err)
				return err
			}

			// 2.通过充值；加余额， 减冻结金额
			walletLog := model.WalletLogExt{
				BusinessNo: order.OrderNo,
				Source:     model.AmountSourceRecharge,
				Remark:     "Top up successfully", // 充值成功
			}
			if err := model.NewMerchantModel(tx).PlusBalance(order.MerchantId, order.OrderAmount, walletLog); err != nil {
				l.Errorf("增加商户[%v]余额[%v]失败, err=%v", order.MerchantId, order.OrderAmount, err)
				return err
			}

			// 3.增加平台收款卡的今日已收金额
			if err := model.NewPlatformBankCardModel(tx).PlusTodayReceived(order.PlatformBankCardId, order.OrderAmount); err != nil {
				l.Errorf("增加平台收款卡[%v]今日已收金额失败, err=%v", order.PlatformBankCardId, err)
				return err
			}

			return nil
		})
		if txErr != nil {
			l.Errorf("审核商户充值订单[%v]失败, err=%v", req.OrderNo, txErr)
			return common.NewCodeError(common.SysDBSave)
		}
	case model.RechargeOrderStatusReject:
		// 修改状态
		if err = model.NewMerchantRechargeOrderModel(l.svcCtx.DbEngine).UpdateStatusToReject(req.OrderNo, req.DealWithRemark); err != nil {
			l.Errorf("修改商户充值订单[%v]状态为[%v]失败, err=%v", req.OrderNo, req.OrderStatus, err)
			return err
		}

	default:
		l.Errorf("不支持当前操作, %v", req.OrderStatus)
		return errors.New(fmt.Sprintf("不支持当前操作, %v", req.OrderStatus))
	}

	return nil
}
