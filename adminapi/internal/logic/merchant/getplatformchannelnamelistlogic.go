package merchant

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetPlatformChannelNameListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPlatformChannelNameListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPlatformChannelNameListLogic {
	return GetPlatformChannelNameListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPlatformChannelNameListLogic) GetPlatformChannelNameList(req types.GetPlatformChannelNameListRequest) (*types.GetPlatformChannelNameListResponse, error) {
	var notIds []int64
	var areaId int64
	if req.MerchantId != 0 {
		merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneById(req.MerchantId)
		if err != nil {
			if err == model.ErrRecordNotFound {
				l.Errorf("商户[%v]不存在", req.MerchantId)
				return nil, common.NewCodeError(common.UserNotExist)
			} else {
				l.Errorf("查询商户[%v]失败, err=%v", req.MerchantId, err)
				return nil, common.NewCodeError(common.SysDBGet)
			}
		}

		if merchant.Status == model.MerchantStatusDisable {
			l.Errorf("商户[%v]已被禁用", req.MerchantId)
			return nil, common.NewCodeError(common.MerchantAlreadyDisable)
		}

		ids, err := model.NewMerchantChannelModel(l.svcCtx.DbEngine).FindMerchantPlatformChannelId(req.MerchantId, req.ChannelType)
		if err != nil {
			l.Errorf("查询商户[%v]已添加的平台通道失败, err=%v", req.MerchantId, err)
			return nil, common.NewCodeError(common.SysDBGet)
		}
		notIds = ids
		areaId = merchant.AreaId
	}

	data, err := model.NewPlatformChannelModel(l.svcCtx.DbEngine).FindManyNotIds(notIds, req.Status, req.ChannelType, areaId)
	if err != nil {
		l.Errorf("查询平台通道失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var list []types.PlatformChannelNameList
	for _, v := range data {
		list = append(list, types.PlatformChannelNameList{
			ChannelId:   v.Id,
			ChannelName: v.ChannelName,
			ChannelType: v.ChannelType,
		})
	}
	return &types.GetPlatformChannelNameListResponse{
		List: list,
	}, nil
}
