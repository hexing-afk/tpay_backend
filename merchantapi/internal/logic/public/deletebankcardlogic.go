package public

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type DeleteBankCardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteBankCardLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteBankCardLogic {
	return DeleteBankCardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteBankCardLogic) DeleteBankCard(merchantId int64, req types.DeleteBankCardRequest) error {
	if err := model.NewMerchantBankCardModel(l.svcCtx.DbEngine).Delete(req.CardId, merchantId); err != nil {
		l.Errorf("删除商户[%v]银行卡失败, err=%v", merchantId, req.CardId)
		return common.NewCodeError(common.SysDBDelete)
	}

	return nil
}
