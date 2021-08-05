package public

import (
	"context"
	"gorm.io/gorm"
	"strings"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"
	"tpay_backend/utils"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type WithdrawLogic struct {
	logx.Logger
	ctx        context.Context
	svcCtx     *svc.ServiceContext
	merchantId int64
}

func NewWithdrawLogic(ctx context.Context, svcCtx *svc.ServiceContext, merchantId int64) WithdrawLogic {
	return WithdrawLogic{
		Logger:     logx.WithContext(ctx),
		ctx:        ctx,
		svcCtx:     svcCtx,
		merchantId: merchantId,
	}
}

func (l *WithdrawLogic) Withdraw(req types.WithdrawRequest) error {
	logx.Infof("req：%+v", req)

	req.CardNumber = strings.ReplaceAll(req.CardNumber, " ", "")

	// 1.查商户信息
	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneById(l.merchantId)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("商户[%v]不存在", l.merchantId)
			return common.NewCodeError(common.NotWithdrawConfig)
		} else {
			l.Errorf("查询商户[%v]失败, err=%v", l.merchantId, err)
			return common.NewCodeError(common.ApplyFail)
		}
	}

	// google验证码是否关闭
	totpIsClose, err := model.NewGlobalConfigModel(l.svcCtx.DbEngine).TotpIsClose()
	if err != nil {
		l.Errorf("查询totp配置失败:%v", err)
		return common.NewCodeError(common.SysDBErr)
	}

	if !totpIsClose { // 没有关闭
		// 验证TOTP密码code
		if !utils.VerifyTOTPPasscode(req.TotpCode, merchant.TotpSecret) {
			return common.NewCodeError(common.LoginCaptchaNotMatch)
		}
	}

	plainPassword, err := common.DecryptPassword(req.PayPassword)
	if err != nil {
		l.Errorf("密码解密发生错误,err:%v, password:%v", err, req.PayPassword)
		return common.NewCodeError(common.SysDBErr)
	}

	if merchant.PayPassword != common.CreateMerchantPayPassword(plainPassword) {
		l.Errorf("商户[%v]支付密码错误", merchant)
		return common.NewCodeError(common.PayPasswordErr)
	}

	if merchant.Balance <= 0 {
		l.Errorf("商户[%v]余额[%v]不足", l.merchantId, merchant.Balance)
		return common.NewCodeError(common.InsufficientBalance)
	}

	// 2.查商户提现配置信息
	config, err := model.NewMerchantWithdrawConfigModel(l.svcCtx.DbEngine).FindOneByMerchantId(l.merchantId)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("商户[%v]没有提现配置", l.merchantId)
			return common.NewCodeError(common.NotWithdrawConfig)
		} else {
			l.Errorf("查询商户[%v]提现配置失败, err=%v", l.merchantId, err)
			return common.NewCodeError(common.ApplyFail)
		}
	}

	if config.SingleMinAmount > req.Amount || config.SingleMaxAmount < req.Amount {
		l.Errorf("提现金额[%v]超出最大[%v]和最小[%v]限制", req.Amount, config.SingleMaxAmount, config.SingleMinAmount)
		return common.NewCodeError(common.AmountOutOfLimit)
	}

	// 4.下单
	order := new(model.MerchantWithdrawOrder)
	order.OrderNo = utils.GetDailyId()
	order.MerchantId = merchant.Id
	order.OrderAmount = req.Amount
	order.Remark = req.Remark
	order.BankName = req.BankName
	order.BranchName = req.BranchName
	order.PayeeName = req.PayeeName
	order.CardNumber = req.CardNumber
	order.OrderStatus = model.WithdrawOrderStatusPending
	order.DeductionMethod = config.DeductionMethod
	order.Currency = merchant.Currency
	order.BankCode = req.BankCode
	order.AreaId = merchant.AreaId

	// 商户手续费
	order.MerchantFee = utils.CalculatePayOrderFeeMerchant(req.Amount, config.SingleFee, config.Rate)

	switch config.DeductionMethod {
	case model.MerchantWithdrawDeductionInner:
		// 商户实际到账金额
		order.RealAmount = order.OrderAmount - order.MerchantFee

		// 商户账户扣减金额
		order.DecreaseAmount = req.Amount

		if merchant.Balance < order.OrderAmount {
			l.Errorf("商户[%v]余额[%v]不足", l.merchantId, merchant.Balance)
			return common.NewCodeError(common.InsufficientBalance)
		}
	case model.MerchantWithdrawDeductionOut:
		// 商户实际到账金额
		order.RealAmount = order.OrderAmount

		// 商户账户扣减金额
		order.DecreaseAmount = order.OrderAmount + order.MerchantFee

		if merchant.Balance < order.DecreaseAmount {
			l.Errorf("商户[%v]余额[%v]不足", l.merchantId, merchant.Balance)
			return common.NewCodeError(common.InsufficientBalance)
		}
	default:
		l.Errorf("提现配置有问题, config: %+v", config)
		return common.NewCodeError(common.ApplyFail)
	}

	l.Infof("订单：%+v", order)

	if err := l.InsertOrder(order); err != nil {
		l.Errorf("提现失败, err=%v", err)
		return common.NewCodeError(common.ApplyFail)
	}

	return nil
}

func (l *WithdrawLogic) InsertOrder(order *model.MerchantWithdrawOrder) error {
	txErr := l.svcCtx.DbEngine.Transaction(func(tx *gorm.DB) error {
		// 1.插入订单
		if err := model.NewMerchantWithdrawOrderModel(tx).Insert(order); err != nil {
			l.Errorf("添加提现订单失败, err=%v", err)
			return err
		}

		// 2.减商户余额，加冻结金额
		log := model.WalletLogExt{
			BusinessNo: order.OrderNo,
			Source:     model.AmountSourceWithdraw,
			Remark:     "",
		}
		if err := model.NewMerchantModel(tx).MinusBalanceFreezeTx(l.merchantId, order.DecreaseAmount, log); err != nil {
			if strings.Contains(err.Error(), model.BalanceErr.Error()) {
				l.Errorf("商户余额不足")
				return err
			} else {
				l.Errorf("减商户余额增加冻结金额失败, err=%v", err)
				return err
			}

		}
		return nil
	})

	return txErr
}
