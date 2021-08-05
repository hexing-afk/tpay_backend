package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"
	"tpay_backend/model"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetAdminListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAdminListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetAdminListLogic {
	return GetAdminListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAdminListLogic) GetAdminList(req types.GetAdminListRequest) (*types.GetAdminListResponse, error) {
	f := model.FindAdminList{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	data, total, err := model.NewAdminModel(l.svcCtx.DbEngine).FindList(f)
	if err != nil {
		l.Logger.Errorf("查询管理员账号列表失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var list []types.Admin
	for _, v := range data {
		list = append(list, types.Admin{
			AdminId:      v.Id,
			Username:     v.Username,
			EnableStatus: v.EnableStatus,
			CreateTime:   v.CreateTime,
			TotpSecret:   v.TotpSecret,
		})
	}

	return &types.GetAdminListResponse{
		Total: total,
		List:  list,
	}, nil
}
