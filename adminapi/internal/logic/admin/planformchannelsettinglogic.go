package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type PlanformChannelSettingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPlanformChannelSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) PlanformChannelSettingLogic {
	return PlanformChannelSettingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PlanformChannelSettingLogic) PlanformChannelSetting(req types.PlanformChannelSettingRequest) (*types.PlanformChannelSettingReply, error) {
	if req.StartTime > req.EndTime {
		return nil, common.NewCodeError(common.EndTimeOverStartTime)
	}
	if req.StartAmount > req.EndAmount {
		return nil, common.NewCodeError(common.AmountOver)
	}
	err := model.NewPlatformChannelModel(l.svcCtx.DbEngine).Update(req.ChannelId, model.PlatformChannel{
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		StartAmount: req.StartAmount,
		EndAmount:   req.EndAmount,
	})
	if err != nil {
		l.Errorf("修改通道配置失败, ChannelId=%v, err=%v", req.ChannelId, err)
		return nil, common.NewCodeError(common.SysDBUpdate)
	}

	return &types.PlanformChannelSettingReply{}, nil
}
