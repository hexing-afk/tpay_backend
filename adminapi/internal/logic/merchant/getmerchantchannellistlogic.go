package merchant

import (
	"context"
	"strings"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetMerchantChannelListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMerchantChannelListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetMerchantChannelListLogic {
	return GetMerchantChannelListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMerchantChannelListLogic) GetMerchantChannelList(req types.GetMerchantChannelListRequest) (*types.GetMerchantChannelListResponse, error) {

	p := model.FindMerchantChannelListReq{
		MerchantId: req.MerchantId,
	}
	channelList, total, err := model.NewMerchantChannelModel(l.svcCtx.DbEngine).FindMerchantChannelList(p)
	if err != nil {
		l.Errorf("查询商户[%v]通道失败, err=%v", req.MerchantId, err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var list []types.MerchantChannelList
	for _, v := range channelList {
		if v.Id == 0 {
			continue
		}
		list = append(list, types.MerchantChannelList{
			ChannelId:             v.Id,
			MerchantId:            v.MerchantId,
			ChannelRate:           v.Rate,
			SingleFee:             v.SingleFee,
			Status:                v.Status,
			PlatformChannelName:   v.PlatformChannelName,
			PlatformChannelType:   v.PlatformChannelType,
			UpstreamChannelWeight: strings.Split(v.UpstreamChannelWeight, ","),
		})
	}

	return &types.GetMerchantChannelListResponse{
		Total: total,
		List:  list,
	}, nil
}
