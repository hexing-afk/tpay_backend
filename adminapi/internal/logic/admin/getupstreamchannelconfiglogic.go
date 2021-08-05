package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetUpstreamChannelConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUpstreamChannelConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetUpstreamChannelConfigLogic {
	return GetUpstreamChannelConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUpstreamChannelConfigLogic) GetUpstreamChannelConfig(req types.GetUpstreamChannelConfigRequest) (*types.GetUpstreamChannelConfigResponse, error) {
	channel, err := model.NewUpstreamChannelModel(l.svcCtx.DbEngine).FindOneById(req.ChannelId)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("通道[%v]不存在", req.ChannelId)
			return nil, common.NewCodeError(common.ChannelNotExist)
		} else {
			l.Errorf("查询通道[%v]失败, err=%v", req.ChannelId, err.Error())
			return nil, common.NewCodeError(common.SysDBGet)
		}
	}

	channelConfig := types.UpstreamChannelConfig{
		ChannelId:       channel.Id,
		DeductionMethod: channel.DeductionMethod,
		Rate:            channel.Rate,
		SingleFee:       channel.SingleFee,
		SingleMaxAmount: channel.SingleMaxAmount,
	}
	return &types.GetUpstreamChannelConfigResponse{
		ChannelConfig: channelConfig,
	}, nil
}
