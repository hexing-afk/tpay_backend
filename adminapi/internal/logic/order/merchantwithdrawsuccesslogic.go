package order

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type MerchantWithdrawSuccessLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMerchantWithdrawSuccessLogic(ctx context.Context, svcCtx *svc.ServiceContext) MerchantWithdrawSuccessLogic {
	return MerchantWithdrawSuccessLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MerchantWithdrawSuccessLogic) MerchantWithdrawSuccess(adminId int64, req types.MerchantWithdrawSuccessRequest) error {
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

	if order.OrderStatus != model.WithdrawOrderStatusPass {
		l.Errorf("提现订单[%v]当前状态[%v]不能直接修改为成功", req.OrderNo, order.OrderStatus)
		return common.NewCodeError(common.OrderNotOp)
	}

	txErr := l.svcCtx.DbEngine.Transaction(func(tx *gorm.DB) error {
		// 1.修改订单状态
		if err := model.NewMerchantWithdrawOrderModel(tx).UpdateStatusToSuccess(order.Id, req.Remark); err != nil {
			l.Errorf("修改提现订单[%v]失败, err=%v", order.OrderNo, err)
			return err
		}

		// 2.减少冻结金额
		log := model.WalletLogExt{
			BusinessNo: order.OrderNo,
			Source:     model.AmountSourceWithdraw,
			Remark:     "", // 提现成功 Withdraw successfully
		}
		if err := model.NewMerchantModel(tx).MinusFrozenAmount(order.MerchantId, order.DecreaseAmount, log); err != nil {
			l.Errorf("减商户[%v]冻结金额[%v]失败, err=%v", order.MerchantId, order.DecreaseAmount, err)
			return err
		}

		// 3.计算并记录平台收益
		data := &model.PlatformWalletLog{
			BusinessNo:  order.OrderNo,
			Source:      model.PlatformIncomeSourceWithdraw,
			MerchantFee: order.MerchantFee,
			UpstreamFee: 0,
			Income:      order.MerchantFee,
			Currency:    order.Currency,
		}
		if err := model.NewPlatformWalletLogModel(tx).Insert(data); err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		l.Errorf("修改提现订单[%v]为成功失败, err=%v", order.OrderNo, err)
	}

	data := &model.AdminWebLog{
		AdminId:     adminId,
		Description: fmt.Sprintf("手动将商户提现订单[%v]改为成功状态[%v]", order.OrderNo, model.WithdrawOrderStatusSuccess),
		Type:        model.LogTypeMerchant,
	}
	if err := model.NewAdminWebLogModel(l.svcCtx.DbEngine).Insert(data); err != nil {
		l.Errorf("记录管理员操作日志失败, log: %+v, err=%v", data, err)
	}
	return nil
}
