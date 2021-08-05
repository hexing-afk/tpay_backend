package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"tpay_backend/model"
	"tpay_backend/payapi/internal/svc"
	"tpay_backend/utils"
)

type PickUpstreamChannel struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	merchant *model.Merchant
}

func NewPickUpstreamChannel(ctx context.Context, svcCtx *svc.ServiceContext, merchant *model.Merchant) *PickUpstreamChannel {
	return &PickUpstreamChannel{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		merchant: merchant,
	}
}

// 按权重选择代收上游通道相关信息
func (l *PickUpstreamChannel) PickPayUpstreamChannelByWeigh(tradeType string) (model.MerchantUpstreamChannelData, error) {
	// 查询可用的通道
	upChannelList, err := model.NewMerchantChannelModel(l.svcCtx.DbEngine).QueryPayUpstreamChannel(l.merchant.Id, l.merchant.AreaId, tradeType, l.merchant.Currency)
	if err != nil {
		l.Errorf("代收-查询可用的通道:err:%v,merchantId:%v,TradeType:%v,Currency:%v", err, l.merchant.Id, tradeType, l.merchant.Currency)
		return model.MerchantUpstreamChannelData{}, errors.New("查询可用的通道失败")
	}

	l.Infof("代收-查询可用的通道列表:upChannelList:%+v", upChannelList)

	if len(upChannelList) == 0 {
		return model.MerchantUpstreamChannelData{}, errors.New("没有可用的通道")
	}

	var weightList []int64
	for _, v := range upChannelList {
		weightList = append(weightList, v.Weight)
	}
	idx := utils.PickByWeight(weightList)

	return upChannelList[idx], nil
}

// 按权重选择代付上游通道相关信息
func (l *PickUpstreamChannel) PickTransferUpstreamChannelByWeigh(tradeType string) (model.MerchantUpstreamChannelData, error) {
	// 查询可用的通道
	upChannelList, err := model.NewMerchantChannelModel(l.svcCtx.DbEngine).QueryTransferUpstreamChannel(l.merchant.Id, l.merchant.AreaId, tradeType, l.merchant.Currency)
	if err != nil {
		return model.MerchantUpstreamChannelData{}, errors.New(fmt.Sprintf("代付-查询可用的通道:err:%v,merchantId:%v,TradeType:%v,Currency:%v", err, l.merchant.Id, tradeType, l.merchant.Currency))
	}

	if len(upChannelList) == 0 {
		return model.MerchantUpstreamChannelData{}, errors.New("没有可用的通道")
	}

	var weightList []int64
	for _, v := range upChannelList {
		weightList = append(weightList, v.Weight)
	}
	idx := utils.PickByWeight(weightList)

	return upChannelList[idx], nil
}
