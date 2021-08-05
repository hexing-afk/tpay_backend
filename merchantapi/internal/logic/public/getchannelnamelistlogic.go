package public

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetChannelNameListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetChannelNameListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetChannelNameListLogic {
	return GetChannelNameListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChannelNameListLogic) GetChannelNameList(merchantId int64, req types.GetChannelNameListReq) (*types.GetChannelNameListResponse, error) {

	p := model.FindMerchantChannelListReq{
		MerchantId:            merchantId,
		ChannelType:           req.ChannelType,
		MerchantChannelStatus: req.ChannelStatus,
	}
	data, total, err := model.NewMerchantChannelModel(l.svcCtx.DbEngine).FindMerchantChannelList(p)
	if err != nil {
		l.Errorf("查询商户[%v]通道失败, err=%v", merchantId, err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var list []types.ChannelNameList
	for _, v := range data {
		list = append(list, types.ChannelNameList{
			ChannelId:   v.PlatformChannelId,
			ChannelName: v.PlatformChannelName,
		})
	}

	return &types.GetChannelNameListResponse{
		Total: total,
		List:  list,
	}, nil
}
