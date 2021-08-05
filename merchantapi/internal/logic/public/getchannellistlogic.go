package public

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetChannelListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetChannelListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetChannelListLogic {
	return GetChannelListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChannelListLogic) GetChannelList(merchantId int64, req types.GetChannelListReq) (*types.GetChannelListResponse, error) {
	p := model.FindMerchantChannelListReq{
		MerchantId:            merchantId,
		MerchantChannelStatus: model.MerchantChannelStatusEnable,
		ChannelType:           req.ChannelType,
	}

	data, total, err := model.NewMerchantChannelModel(l.svcCtx.DbEngine).FindMerchantChannelList(p)
	if err != nil {
		l.Errorf("查询商户[%v]通道列表失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var list []types.ChannelList
	for _, v := range data {
		list = append(list, types.ChannelList{
			ChannelId:   v.PlatformChannelId,
			ChannelName: v.PlatformChannelName,
			ChannelCode: v.PlatformChannelCode,
			ChannelType: v.PlatformChannelType,
			ChannelRate: v.Rate,
			Status:      v.Status,
			SingleFee:   v.SingleFee,
		})
	}

	return &types.GetChannelListResponse{
		Total: total,
		List:  list,
	}, nil
}
