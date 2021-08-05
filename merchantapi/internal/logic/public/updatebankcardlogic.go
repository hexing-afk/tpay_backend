package public

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type UpdateBankCardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateBankCardLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateBankCardLogic {
	return UpdateBankCardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateBankCardLogic) UpdateBankCard(merchantId int64, req types.UpdateBankCardRequest) error {
	data := model.MerchantBankCard{
		BankName:    req.BankName,
		BranchName:  req.BranchName,
		AccountName: req.AccountName,
		CardNumber:  req.CardNumber,
		Remark:      req.Remark,
	}
	if err := model.NewMerchantBankCardModel(l.svcCtx.DbEngine).Update(req.CardId, merchantId, data); err != nil {
		l.Errorf("修改商户[%v]银行卡[%v]失败, err=%v", merchantId, req.CardId, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	return nil
}
