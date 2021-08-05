package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ModifyPlatformBankCardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModifyPlatformBankCardLogic(ctx context.Context, svcCtx *svc.ServiceContext) ModifyPlatformBankCardLogic {
	return ModifyPlatformBankCardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModifyPlatformBankCardLogic) ModifyPlatformBankCard(req types.ModifyPlatformBankCardRequest) error {
	exist, err := model.NewPlatformBankCardModel(l.svcCtx.DbEngine).CheckById(req.CardId)
	if err != nil {
		l.Errorf("查询平台收款卡[%v]是否存在失败, err=%v", req.CardId, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	if !exist {
		l.Errorf("平台收款卡[%v]不存在", req.CardId)
		return common.NewCodeError(common.UpdateContentNotExist)
	}

	data := model.PlatformBankCard{
		BankName:    req.BankName,
		AccountName: req.AccountName,
		CardNumber:  req.CardNumber,
		BranchName:  req.BranchName,
		Currency:    req.Currency,
		MaxAmount:   req.MaxAmount,
		QrCode:      req.QrCode,
		Remark:      req.Remark,
	}
	if err := model.NewPlatformBankCardModel(l.svcCtx.DbEngine).Update(req.CardId, data); err != nil {
		l.Errorf("修改平台收款卡[%v]失败, err=%v", req.CardId, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	return nil
}
