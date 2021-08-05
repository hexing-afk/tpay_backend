package public

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddBankCardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddBankCardLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddBankCardLogic {
	return AddBankCardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddBankCardLogic) AddBankCard(merchantId int64, req types.AddBankCardRequest) (*types.AddBankCardResponse, error) {
	exist, err := model.NewMerchantBankCardModel(l.svcCtx.DbEngine).CheckMerchantCard(merchantId, req.BankName, req.CardNumber)
	if err != nil {
		l.Errorf("查询商户银行卡是否重复失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBAdd)
	}

	if exist {
		l.Errorf("商户[%v]已经添加了银行卡[%v-%v]是否重复失败", merchantId, req.BankName, req.CardNumber)
		return nil, common.NewCodeError(common.BankCardRepeatAdd)
	}

	merchantInfo, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneById(merchantId)
	if err != nil {
		l.Errorf("获取商户[%v]信息失败", merchantId)
		return nil, common.NewCodeError(common.BankCardRepeatAdd)
	}

	data := &model.MerchantBankCard{
		MerchantId:  merchantId,
		BankName:    req.BankName,
		BranchName:  req.BranchName,
		AccountName: req.AccountName,
		CardNumber:  req.CardNumber,
		//Currency:    req.Currency,
		Currency: merchantInfo.Currency,
		Remark:   req.Remark,
	}
	if err := model.NewMerchantBankCardModel(l.svcCtx.DbEngine).Insert(data); err != nil {
		l.Errorf("添加商户[%v]银行卡[%v-%v]失败, err=%v", merchantId, req.BankName, req.CardNumber, err)
		return nil, common.NewCodeError(common.SysDBAdd)
	}

	return &types.AddBankCardResponse{CardId: data.Id}, nil
}
