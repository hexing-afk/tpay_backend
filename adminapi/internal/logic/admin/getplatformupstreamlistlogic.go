package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetPlatformUpstreamListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPlatformUpstreamListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPlatformUpstreamListLogic {
	return GetPlatformUpstreamListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPlatformUpstreamListLogic) GetPlatformUpstreamList(req types.GetPlatformUpstreamListRequest) (*types.GetPlatformUpstreamListResponse, error) {
	channel, err := model.NewPlatformChannelModel(l.svcCtx.DbEngine).FindOneById(req.ChannelId)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("通道[%v]不存在", req.ChannelId)
			return nil, common.NewCodeError(common.UpdateContentNotExist)
		} else {
			l.Errorf("查询通道[%v]失败, err=%v", req.ChannelId, err)
			return nil, common.NewCodeError(common.SysDBGet)
		}
	}

	// 查询通道已关联的上游通道
	linkedData, err := model.NewPlatformChannelUpstreamModel(l.svcCtx.DbEngine).FindUpstreamByPlatform(req.ChannelId)
	if err != nil {
		l.Errorf("查询通道[%v]已关联的上游通道失败, err=%v", req.ChannelId, err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var linkedList []types.PlatformUpstreamChannel
	var upChannelIds []int64
	for _, v := range linkedData {
		upChannelIds = append(upChannelIds, v.UpstreamChannelId)
		linkedList = append(linkedList, types.PlatformUpstreamChannel{
			UpstreamChannelId:   v.UpstreamChannelId,
			UpstreamChannelName: v.UpstreamChannelName,
		})
	}

	notLinkData, err := model.NewUpstreamChannelModel(l.svcCtx.DbEngine).FindManyNotIds(upChannelIds, channel.ChannelType, channel.AreaId)
	if err != nil {
		l.Errorf("查询通道[%v]未关联的上游通道失败, err=%v", req.ChannelId, err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var notLinkList []types.PlatformUpstreamChannel
	for _, v := range notLinkData {
		notLinkList = append(notLinkList, types.PlatformUpstreamChannel{
			UpstreamChannelId:   v.Id,
			UpstreamChannelName: v.ChannelName,
		})
	}

	return &types.GetPlatformUpstreamListResponse{
		ChannelId:        req.ChannelId,
		LinkedChannel:    linkedList,
		NotLinkedChannel: notLinkList,
	}, nil
}
