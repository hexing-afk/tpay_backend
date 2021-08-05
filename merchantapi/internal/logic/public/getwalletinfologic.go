package public

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetWalletInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	userId int64
}

func NewGetWalletInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext, userId int64) GetWalletInfoLogic {
	return GetWalletInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		userId: userId,
	}
}

func (l *GetWalletInfoLogic) GetWalletInfo() (*types.GetWalletInfoReply, error) {
	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneById(l.userId)
	if err != nil {
		logx.Errorf("查询用户[%v]信息失败,err[%v]", l.userId, err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	return &types.GetWalletInfoReply{
		Currency:     merchant.Currency,
		Balance:      merchant.Balance,
		FrozenAmount: merchant.FrozenAmount,
	}, nil
}
