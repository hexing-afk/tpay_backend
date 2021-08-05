package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type DeletePlatformChannelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeletePlatformChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeletePlatformChannelLogic {
	return DeletePlatformChannelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeletePlatformChannelLogic) DeletePlatformChannel(req types.DeletePlatformChannelRequest) error {
	// 检查平台通道是否有关联上游
	isHave, err := model.NewPlatformChannelUpstreamModel(l.svcCtx.DbEngine).CheckPlatformById(req.ChannelId)
	if err != nil {
		l.Errorf("查询通道[%v]是否有关联上游通道失败, err=%v", req.ChannelId, err)
		return common.NewCodeError(common.SysDBDelete)
	}

	if isHave {
		l.Errorf("通道[%v]有关联上游通道, 不能删除", req.ChannelId)
		return common.NewCodeError(common.ChannelHaveLinkedChannelNotDelete)
	}

	// 检查是否有商户绑定平台通道
	mChannels, err := model.NewMerchantChannelModel(l.svcCtx.DbEngine).FindOneByPlatformId(req.ChannelId)
	if err != nil {
		l.Errorf("平台通道绑定的商户失败, channelId=%v, err=%v", req.ChannelId, err)
		return common.NewCodeError(common.InvalidUpstreamChannel)
	}

	if len(mChannels) > 0 {
		l.Errorf("通道[%v]有关联商户, 不能删除", req.ChannelId)
		return common.NewCodeError(common.ChannelHaveLinkedChannelNotDelete)
	}

	if err := model.NewPlatformChannelModel(l.svcCtx.DbEngine).Delete(req.ChannelId); err != nil {
		l.Errorf("删除通道[%v]失败, err=%v", req.ChannelId, err)
		return common.NewCodeError(common.SysDBDelete)
	}

	return nil
}
