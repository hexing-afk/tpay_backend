package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type DeletePlatformBankCardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePlatformBankCardLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeletePlatformBankCardLogic {
	return DeletePlatformBankCardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePlatformBankCardLogic) DeletePlatformBankCard(req types.DeletePlatformBankCardRequest) error {
	if err := model.NewPlatformBankCardModel(l.svcCtx.DbEngine).Delete(req.CardId); err != nil {
		l.Errorf("删除银行卡[%v]失败, err=%v", req.CardId, err)
		return common.NewCodeError(common.SysDBDelete)
	}

	return nil
}
