package logic

import (
	"context"

	"tpay_backend/payapi/internal/svc"
	"tpay_backend/payapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetUnionpayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUnionpayLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetUnionpayLogic {
	return GetUnionpayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUnionpayLogic) GetUnionpay(req types.GetUnionpayReq) (*types.GetUnionpayReply, error) {
	// todo: add your logic here and delete this line

	return &types.GetUnionpayReply{}, nil
}
