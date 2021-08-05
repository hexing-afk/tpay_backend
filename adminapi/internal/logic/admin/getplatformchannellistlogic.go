package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetPlatformChannelListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPlatformChannelListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPlatformChannelListLogic {
	return GetPlatformChannelListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPlatformChannelListLogic) GetPlatformChannelList(req types.GetPlatformChannelListRequest) (*types.GetPlatformChannelListResponse, error) {
	f := model.FindPlatformChannelList{
		Search:   req.Search,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	data, total, err := model.NewPlatformChannelModel(l.svcCtx.DbEngine).FindList(f)
	if err != nil {
		l.Errorf("查询下游通道列表失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var list []types.PlatformChannelList
	for _, v := range data {
		list = append(list, types.PlatformChannelList{
			ChannelId:           v.Id,
			ChannelName:         v.ChannelName,
			ChannelCode:         v.ChannelCode,
			ChannelDesc:         v.ChannelDesc,
			ChannelType:         v.ChannelType,
			Status:              v.Status,
			CreateTime:          v.CreateTime,
			UpdateTime:          v.UpdateTime,
			UpstreamChannelName: v.UpstreamChannelName,
			AreaId:              v.AreaId,
			AreaName:            v.AreaName,
		})
	}

	return &types.GetPlatformChannelListResponse{
		Total: total,
		List:  list,
	}, nil
}
