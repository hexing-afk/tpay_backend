package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetUpstreamChannelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUpstreamChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetUpstreamChannelLogic {
	return GetUpstreamChannelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUpstreamChannelLogic) GetUpstreamChannel(req types.GetUpstreamChannelRequest) (*types.GetUpstreamChannelResponse, error) {
	data, err := model.NewUpstreamChannelModel(l.svcCtx.DbEngine).FindOneById(req.ChannelId)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("通道[%v]不存在", req.ChannelId)
			return nil, common.NewCodeError(common.ChannelNotExist)
		} else {
			l.Errorf("查询通道[%v]失败, err=%v", req.ChannelId, err.Error())
			return nil, common.NewCodeError(common.SysDBGet)
		}
	}

	channel := types.UpstreamChannel{
		ChannelId:    data.Id,
		ChannelName:  data.ChannelName,
		ChannelCode:  data.ChannelCode,
		ChannelDesc:  data.ChannelDesc,
		Currency:     data.Currency,
		UpstreamId:   data.UpstreamId,
		UpstreamName: data.UpstreamName,
		ChannelType:  data.ChannelType,
	}

	return &types.GetUpstreamChannelResponse{
		Channel: channel,
	}, nil
}
