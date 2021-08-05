package admin

import (
	"context"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetPlanformChannelSettingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPlanformChannelSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPlanformChannelSettingLogic {
	return GetPlanformChannelSettingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPlanformChannelSettingLogic) GetPlanformChannelSetting(req types.GetPlanformChannelSettingRequest) (*types.GetPlanformChannelSettingReply, error) {
	planformChannel, err := model.NewPlatformChannelModel(l.svcCtx.DbEngine).FindOneById(req.ChannelId)
	if err != nil {
		l.Errorf("查询通道配置失败, ChannelId=%v, err=%v", req.ChannelId, err)
		planformChannel.StartTime = 0
		planformChannel.EndTime = 0
		planformChannel.StartAmount = 0
		planformChannel.EndAmount = 0
	}

	return &types.GetPlanformChannelSettingReply{
		StartTime:   planformChannel.StartTime,
		EndTime:     planformChannel.EndTime,
		StartAmount: planformChannel.StartAmount,
		EndAmount:   planformChannel.EndAmount,
	}, nil
}
