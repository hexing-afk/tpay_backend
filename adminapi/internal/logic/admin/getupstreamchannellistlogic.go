package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetUpstreamChannelListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUpstreamChannelListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetUpstreamChannelListLogic {
	return GetUpstreamChannelListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUpstreamChannelListLogic) GetUpstreamChannelList(req types.GetUpstreamChannelListRequest) (*types.GetUpstreamChannelListResponse, error) {
	f := model.FindUpstreamChannelList{
		Search:          req.Search,
		ChannelType:     req.ChannelType,
		StartCreateTime: req.StartCreateTime,
		EndCreateTime:   req.EndCreateTime,
		Page:            req.Page,
		PageSize:        req.PageSize,
	}
	data, total, err := model.NewUpstreamChannelModel(l.svcCtx.DbEngine).FindList(f)
	if err != nil {
		l.Errorf("查询上游通道列表失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var list []types.UpstreamChannelList
	for _, v := range data {
		list = append(list, types.UpstreamChannelList{
			ChannelId:       v.Id,
			ChannelName:     v.ChannelName,
			ChannelCode:     v.ChannelCode,
			ChannelDesc:     v.ChannelDesc,
			Currency:        v.Currency,
			ChannelType:     v.ChannelType,
			UpstreamName:    v.UpstreamName,
			Rate:            v.Rate,
			DeductionMethod: v.DeductionMethod,
			Status:          v.Status,
			UpdateTime:      v.UpdateTime,
			SingleFee:       v.SingleFee,
		})
	}

	return &types.GetUpstreamChannelListResponse{
		Total: total,
		List:  list,
	}, nil
}
