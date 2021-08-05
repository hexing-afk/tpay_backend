package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetUpstreamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUpstreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetUpstreamLogic {
	return GetUpstreamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUpstreamLogic) GetUpstream(req types.GetUpstreamRequest) (*types.GetUpstreamResponse, error) {
	data, err := model.NewUpstreamModel(l.svcCtx.DbEngine).FindOneById(req.UpstreamId)
	if err != nil {
		return nil, common.NewCodeError(common.SysDBGet)
	}

	upstream := types.Upstream{
		UpstreamId:         data.Id,
		UpstreamName:       data.UpstreamName,
		CallConfig:         data.CallConfig,
		CreateTime:         data.CreateTime,
		UpstreamMerchantNo: data.UpstreamMerchantNo,
		UpstreamCode:       data.UpstreamCode,
		AreaId:             data.AreaId,
	}

	return &types.GetUpstreamResponse{
		Upstream: upstream,
	}, nil
}
