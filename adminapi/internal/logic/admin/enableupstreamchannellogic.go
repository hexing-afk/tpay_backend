package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type EnableUpstreamChannelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEnableUpstreamChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) EnableUpstreamChannelLogic {
	return EnableUpstreamChannelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EnableUpstreamChannelLogic) EnableUpstreamChannel(req types.EnableUpstreamChannelRequest) error {

	var err error
	switch req.Enable {
	case model.UpstreamChannelStatusEnable:
		err = model.NewUpstreamChannelModel(l.svcCtx.DbEngine).EnableChannel(req.ChannelId)
	case model.UpstreamChannelStatusDisable:
		err = model.NewUpstreamChannelModel(l.svcCtx.DbEngine).DisableChannel(req.ChannelId)
	default:
		l.Errorf("不支持该操作[%v]", req.Enable)
		return common.NewCodeError(common.InvalidParam)
	}
	if err != nil {
		l.Errorf("启用|禁用通道[%v]失败, err=%v", req.ChannelId, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	return nil
}
