package public

import (
	"context"
	"tpay_backend/merchantapi/internal/common"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	userId int64
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext, userId int64) LogoutLogic {
	return LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		userId: userId,
	}
}

func (l *LogoutLogic) Logout(token string) (*types.LogoutResponse, error) {

	if err := l.svcCtx.RedisSession.Logout(l.userId, token); err != nil {
		l.Errorf("退出登录失败:err:%v,userId:%v,token:%v", err, l.userId, token)
		return nil, common.NewCodeError(common.UserLogoutFailed)
	}

	return &types.LogoutResponse{}, nil
}
