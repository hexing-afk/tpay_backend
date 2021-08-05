package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type DeleteUpstreamChannelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteUpstreamChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteUpstreamChannelLogic {
	return DeleteUpstreamChannelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUpstreamChannelLogic) DeleteUpstreamChannel(req types.DeleteUpstreamChannelRequest) error {
	// 查询上游通道是否有下游通道关联
	isHave, err := model.NewPlatformChannelUpstreamModel(l.svcCtx.DbEngine).CheckUpstreamById(req.ChannelId)
	if err != nil {
		l.Errorf("查询通道[%v]是否被下游通道关联失败, err=%v", req.ChannelId, err)
		return common.NewCodeError(common.SysDBDelete)
	}

	if isHave {
		l.Errorf("查询通道[%v]有下游通道关联, 不能删除", req.ChannelId)
		return common.NewCodeError(common.ChannelHaveLinkedChannelNotDelete)
	}

	if err := model.NewUpstreamChannelModel(l.svcCtx.DbEngine).Delete(req.ChannelId); err != nil {
		l.Errorf("删除通道[%v]失败, err=%v", req.ChannelId, err)
		return common.NewCodeError(common.SysDBDelete)
	}
	return nil
}
