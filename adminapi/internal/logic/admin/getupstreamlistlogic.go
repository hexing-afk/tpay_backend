package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetUpstreamListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUpstreamListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetUpstreamListLogic {
	return GetUpstreamListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUpstreamListLogic) GetUpstreamList(req types.GetUpstreamListRequest) (*types.GetUpstreamListResponse, error) {
	f := model.FindUpstreamList{
		UpstreamName: req.UpstreamName,
		Page:         req.Page,
		PageSize:     req.PageSize,
	}
	data, total, err := model.NewUpstreamModel(l.svcCtx.DbEngine).FindList(f)
	if err != nil {
		l.Errorf("查询上游列表失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var list []types.UpstreamList
	for _, v := range data {
		list = append(list, types.UpstreamList{
			UpstreamId:         v.Id,
			UpstreamName:       v.UpstreamName,
			CreateTime:         v.CreateTime,
			UpstreamMerchantNo: v.UpstreamMerchantNo,
			UpstreamCode:       v.UpstreamCode,
			AreaName:           v.AreaName,
		})
	}

	return &types.GetUpstreamListResponse{
		Total: total,
		List:  list,
	}, nil
}
