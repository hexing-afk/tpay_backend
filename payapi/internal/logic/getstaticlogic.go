package logic

import (
	"context"

	"tpay_backend/payapi/internal/svc"
	"tpay_backend/payapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetStaticLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetStaticLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetStaticLogic {
	return GetStaticLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetStaticLogic) GetStatic(req types.GetStaticReq) (*types.GetStaticReply, error) {
	// todo: add your logic here and delete this line

	return &types.GetStaticReply{}, nil
}
