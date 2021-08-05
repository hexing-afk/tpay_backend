package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AreaListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAreaListLogic(ctx context.Context, svcCtx *svc.ServiceContext) AreaListLogic {
	return AreaListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AreaListLogic) AreaList() (*types.AreaListReply, error) {
	dataList, err := model.NewAreaModel(l.svcCtx.DbEngine).FindMany()
	if err != nil {
		return nil, common.NewCodeError(common.Success)
	}

	var list []types.AreaData
	for _, data := range dataList {
		list = append(list, types.AreaData{
			Id:       data.Id,
			AreaName: data.AreaName,
		})
	}

	return &types.AreaListReply{
		List: list,
	}, nil
}
