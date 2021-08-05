package order

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type UpTransferChannelListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpTransferChannelListLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpTransferChannelListLogic {
	return UpTransferChannelListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpTransferChannelListLogic) UpTransferChannelList(req types.UpTransferChannelListRequest) (*types.UpTransferChannelListResponse, error) {
	f := model.FindUpstreamChannelList{
		ChannelType: model.UpstreamChannelTypeTransfer,
		Currency:    req.Currency,
		Status:      model.UpstreamChannelStatusEnable,
		Page:        1,
		PageSize:    1000,
	}
	data, _, err := model.NewUpstreamChannelModel(l.svcCtx.DbEngine).FindList(f)
	if err != nil {
		l.Errorf("查询上游通道列表失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var list []types.UpTransferChannel
	for _, v := range data {
		list = append(list, types.UpTransferChannel{
			ChannelId:   v.Id,
			ChannelName: v.ChannelName,
			ChannelCode: v.ChannelCode,
		})
	}

	return &types.UpTransferChannelListResponse{
		List: list,
	}, nil
}
