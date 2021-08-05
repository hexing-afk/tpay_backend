package merchant

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetMerchantChannelUpstreamListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMerchantChannelUpstreamListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetMerchantChannelUpstreamListLogic {
	return GetMerchantChannelUpstreamListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMerchantChannelUpstreamListLogic) GetMerchantChannelUpstreamList(req types.GetMerchantChannelUpstreamListRequest) (*types.GetMerchantChannelUpstreamListResponse, error) {
	f := model.FindMerchantChannelUpstreamList{
		MerchantId:        req.MerchantId,
		MerchantChannelId: req.ChannelId,
		Page:              req.Page,
		PageSize:          req.PageSize,
	}
	data, total, err := model.NewMerchantChannelUpstreamModel(l.svcCtx.DbEngine).FindList(f)
	if err != nil {
		l.Errorf("查询商户[%v]通道[%v]上游失败, err=%v", req.MerchantId, req.ChannelId, err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var list []types.MerchantChannelUpstreamList
	for _, v := range data {
		list = append(list, types.MerchantChannelUpstreamList{
			Id:                    v.Id,
			UpstreamChannelName:   v.UpstreamChannelName,
			UpstreamChannelWeight: v.Weight,
		})
	}

	return &types.GetMerchantChannelUpstreamListResponse{
		Total: total,
		List:  list,
	}, nil
}
