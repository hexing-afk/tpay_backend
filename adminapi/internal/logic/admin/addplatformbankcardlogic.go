package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddPlatformBankCardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddPlatformBankCardLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddPlatformBankCardLogic {
	return AddPlatformBankCardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddPlatformBankCardLogic) AddPlatformBankCard(req types.AddPlatformBankCardRequest) error {
	isHave, err := model.NewPlatformBankCardModel(l.svcCtx.DbEngine).CheckBankCard(req.BankName, req.CardNumber, req.Currency)
	if err != nil {
		l.Errorf("检查银行卡是否重复失败, err=%v", err)
		return common.NewCodeError(common.SysDBAdd)
	}

	if isHave {
		l.Errorf("检查银行卡[%v-%v-%v]重复添加", req.BankName, req.CardNumber, req.Currency)
		return common.NewCodeError(common.PlatformBankCardRepetition)
	}

	data := &model.PlatformBankCard{
		BankName:    req.BankName,
		AccountName: req.AccountName,
		CardNumber:  req.CardNumber,
		BranchName:  req.BranchName,
		Currency:    req.Currency,
		MaxAmount:   req.MaxAmount,
		QrCode:      req.QrCode,
		Remark:      req.Remark,
		Status:      model.PlatformBankCardEnable,
	}
	if err := model.NewPlatformBankCardModel(l.svcCtx.DbEngine).Insert(data); err != nil {
		l.Errorf("添加平台收款卡失败, err=%v", err)
		return common.NewCodeError(common.SysDBAdd)
	}

	return nil
}
