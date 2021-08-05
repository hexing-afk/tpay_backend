package public

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"
	"tpay_backend/utils"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type RechargeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	userId int64
}

func NewRechargeLogic(ctx context.Context, svcCtx *svc.ServiceContext, userId int64) RechargeLogic {
	return RechargeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		userId: userId,
	}
}

func (l *RechargeLogic) Recharge(req types.RechargeReq) (*types.RechargeReply, error) {

	if req.Amount <= 0 {
		l.Errorf("充值金额不允许为0")
		return nil, common.NewCodeError(common.AmountFail)
	}

	//查询商户的信息
	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneById(l.userId)
	if err != nil {
		l.Errorf("查询商户信息失败,userId=%v,err=%v", l.userId, err)
		return nil, common.NewCodeError(common.ApplyFail)
	}

	//确认平台收款卡是否存在
	platformBankCard, err := model.NewPlatformBankCardModel(l.svcCtx.DbEngine).FindOneById(req.BankCardId)
	if err != nil {
		l.Errorf("查询平台收款卡信息失败,BankCardId=%v,err=%v", req.BankCardId, err)
		return nil, common.NewCodeError(common.HankCardNotExist)
	}

	if platformBankCard == nil {
		l.Errorf("查询平台收款卡信息为空, BankCardId=%v", req.BankCardId)
		return nil, common.NewCodeError(common.HankCardNotExist)
	}

	if platformBankCard.Status != model.PlatformBankCardEnable {
		l.Errorf("平台收款卡状态为禁用, BankCardId=%v", req.BankCardId)
		return nil, common.NewCodeError(common.HankCardNotExist)
	}

	//金额 + 平台收款卡今日已收金额 > 平台收款卡每日最大收款额度  将不允许充值申请
	if platformBankCard.MaxAmount < platformBankCard.TodayReceived+req.Amount {
		l.Error("平台收款卡已超出当天的收款额度, MaxAmount[%v], TodayReceived[%v], Amount[%v]", platformBankCard.MaxAmount, platformBankCard.TodayReceived, req.Amount)
		return nil, common.NewCodeError(common.BankCardMaxAmountLacking)
	}

	if merchant.Currency != platformBankCard.Currency {
		l.Errorf("平台收款卡币种与商户币种不一致, merchant.Currency=%v, platformBankCard.Currency=%v ", merchant.Currency, platformBankCard.Currency)
		return nil, common.NewCodeError(common.HankCardNotExist)
	}

	d := &model.MerchantRechargeOrder{
		OrderNo:            utils.GetDailyId(),               // 订单号
		OrderAmount:        req.Amount,                       // 订单金额
		MerchantId:         l.userId,                         // 商户id
		OrderStatus:        model.RechargeOrderStatusPending, // 订单状态
		RechargeRemark:     req.Remark,                       // 充值备注
		PlatformBankCardId: req.BankCardId,                   // 平台收款卡id
		BankName:           platformBankCard.BankName,        // 收款银行
		CardNumber:         platformBankCard.CardNumber,      // 收款卡号
		PayeeName:          platformBankCard.AccountName,     // 收款人姓名
		BranchName:         platformBankCard.BranchName,      // 支行名称
		Currency:           merchant.Currency,                // 币种
	}

	if err := model.NewMerchantRechargeOrderModel(l.svcCtx.DbEngine).Insert(d); err != nil {
		l.Errorf("插入充值订单失败,d=[%+v] err=[%v]", d, err)
		return nil, common.NewCodeError(common.RechargeFailed)
	}

	return &types.RechargeReply{}, nil
}
