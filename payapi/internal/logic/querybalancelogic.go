package logic

import (
	"context"
	"tpay_backend/model"

	"tpay_backend/payapi/internal/svc"
	"tpay_backend/payapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type QueryBalanceLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	merchant *model.Merchant
}

func NewQueryBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext, merchant *model.Merchant) QueryBalanceLogic {
	return QueryBalanceLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		merchant: merchant,
	}
}

func (l *QueryBalanceLogic) QueryBalance(req types.QueryBalanceReq) (*types.QueryBalanceReply, error) {
	return &types.QueryBalanceReply{
		Balance:  l.merchant.Balance,
		Currency: l.merchant.Currency,
	}, nil
}
