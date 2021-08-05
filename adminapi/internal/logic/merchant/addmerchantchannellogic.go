package merchant

import (
	"context"
	"gorm.io/gorm"
	"strings"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"
	"tpay_backend/utils"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddMerchantChannelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddMerchantChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddMerchantChannelLogic {
	return AddMerchantChannelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddMerchantChannelLogic) AddMerchantChannel(req types.AddMerchantChannelRequest) error {
	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindMerchantChannel(req.MerchantId)
	if err != nil {
		l.Errorf("查询商户[%v]失败, err=%v", req.MerchantId, err)
		return common.NewCodeError(common.SysDBGet)
	}

	if merchant.Id == 0 {
		l.Errorf("商户[%v]不存在", req.MerchantId)
		return common.NewCodeError(common.UserNotExist)
	}

	if merchant.Status == model.MerchantStatusDisable {
		l.Errorf("商户[%v]已被禁用", req.MerchantId)
		return common.NewCodeError(common.MerchantAlreadyDisable)
	}

	if merchant.PlatformChannelIds != "" {
		if utils.InSlice(string(req.ChannelId), strings.Split(merchant.PlatformChannelIds, ",")) {
			l.Errorf("商户已经添加通道[%v], 不能重复添加", req.ChannelId, err)
			return common.NewCodeError(common.ChannelRepetition)
		}
	}

	// 检查添加的通道是否可用
	plChannel, err := model.NewPlatformChannelModel(l.svcCtx.DbEngine).FindOneById(req.ChannelId)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("通道[%v]不存在", req.ChannelId)
			return common.NewCodeError(common.ChannelNotExist)
		} else {
			l.Errorf("查询通道[%v]失败, err=%v", req.ChannelId, err)
			return common.NewCodeError(common.SysDBAdd)
		}
	}
	if plChannel.Status == model.PlatformChannelStatusDisable {
		l.Errorf("通道[%v]已被禁用, 不能重复添加", req.ChannelId)
		return common.NewCodeError(common.ChannelDisable)
	}

	if merchant.AreaId != 0 && merchant.AreaId != plChannel.AreaId {
		l.Errorf("通道地区与商户地区不一致, req[%v]", req)
		return common.NewCodeError(common.MerchantAreaNotSame)
	}

	// 查询商户可用的上游通道配置
	upChannels, err := model.NewPlatformChannelUpstreamModel(l.svcCtx.DbEngine).FindMerchantUpstreamChannel(req.ChannelId, merchant.Currency)
	if err != nil {
		l.Errorf("查询平台通道[%v]的上游失败, err=%v", req.ChannelId, err)
		return common.NewCodeError(common.SysDBAdd)
	}

	if len(upChannels) == 0 {
		l.Errorf("平台通道[%v]还没有关联到任何一个上游通道, 不能添加", req.ChannelId)
		return common.NewCodeError(common.ChannelNotUpstreamChannel)
	}

	txErr := l.svcCtx.DbEngine.Transaction(func(tx *gorm.DB) error {
		// 商户通道
		channel := &model.MerchantChannel{
			MerchantId:        req.MerchantId,
			PlatformChannelId: req.ChannelId,
			Status:            model.MerchantChannelStatusDisable,
		}
		if err := model.NewMerchantChannelModel(tx).Insert(channel); err != nil {
			l.Errorf("添加商户[%v]通道[%v]失败, err=%v", req.MerchantId, req.ChannelId, err)
			return err
		}

		// 商户通道上游
		var data []*model.MerchantChannelUpstream
		for _, v := range upChannels {
			data = append(data, &model.MerchantChannelUpstream{
				MerchantId:        req.MerchantId,
				MerchantChannelId: channel.Id,
				UpstreamChannelId: v.UpstreamChannelId,
				Weight:            10,
			})
		}
		if err := model.NewMerchantChannelUpstreamModel(l.svcCtx.DbEngine).Inserts(data); err != nil {
			l.Errorf("添加商户[%v]通道[%v]的上游通道配置失败, err=%v", req.MerchantId, req.ChannelId, err)
			return err
		}

		return nil
	})

	if txErr != nil {
		l.Errorf("添加商户[%v]通道[%v]失败, err=%v", req.MerchantId, req.ChannelId, err)
		return common.NewCodeError(common.SysDBAdd)
	}

	return nil
}
