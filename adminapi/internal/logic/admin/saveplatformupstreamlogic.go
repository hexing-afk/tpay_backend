package admin

import (
	"context"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"
	"tpay_backend/utils"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type SavePlatformUpstreamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSavePlatformUpstreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) SavePlatformUpstreamLogic {
	return SavePlatformUpstreamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SavePlatformUpstreamLogic) SavePlatformUpstream(req types.SavePlatformUpstreamRequest) error {
	// 1.查询平台已绑定的上游通道
	channel, err := model.NewPlatformChannelModel(l.svcCtx.DbEngine).FindUpChannelById(req.ChannelId)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("通道[%v]不存在", req.ChannelId)
			return common.NewCodeError(common.UpdateContentNotExist)
		} else {
			l.Errorf("查询通道[%v]失败, err=%v", req.ChannelId, err)
			return common.NewCodeError(common.SysDBSave)
		}
	}

	if channel.Id == 0 {
		l.Errorf("通道[%v]不存在", req.ChannelId)
		return common.NewCodeError(common.UpdateContentNotExist)
	}

	// 2.去除重复的上游通道
	upChannelIds := utils.DistinctSliceOnIt64(req.UpstreamChannelIds)

	// 3.平台通道原绑定的上游通道
	var oldUpChannel []int64
	if channel.UpstreamChannelIds != "" {
		for _, v := range strings.Split(channel.UpstreamChannelIds, ",") {
			id, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				l.Errorf("数据类型转换失败, err=%v", err)
				return common.NewCodeError(common.SystemInternalErr)
			}
			oldUpChannel = append(oldUpChannel, id)
		}
	}

	// 4.需要新增和删除的绑定关系
	addUp, delUp := utils.TemplateSliceInt64(oldUpChannel, upChannelIds)

	// 4.1需要新增的平台通道上游绑定关系
	var data []*model.PlatformChannelUpstream
	for _, v := range addUp {
		data = append(data, &model.PlatformChannelUpstream{
			PlatformChannelId: req.ChannelId,
			UpstreamChannelId: v,
		})
	}

	if len(data) > 0 {
		// 5.查询需要新增的上游通道是否存在
		upChannel, err := model.NewUpstreamChannelModel(l.svcCtx.DbEngine).FindManyByIds(addUp, channel.ChannelType)
		if err != nil {
			l.Errorf("查询上游通道失败, ids=%v, channelType=%v, err=%v", upChannelIds, channel.ChannelType, err)
			return common.NewCodeError(common.SysDBSave)
		}

		if len(addUp) != len(upChannel) {
			l.Errorf("查询出的上游通道数量[%v]和要添加的数量[%v]不相等", len(upChannel), len(upChannelIds))
			l.Errorf("查询条件为: ids=%v, channelType=%v", addUp, channel.ChannelType)
			return common.NewCodeError(common.InvalidUpstreamChannel)
		}
	}

	// 6.查询已绑定平台通道的商户
	mChannels, err := model.NewMerchantChannelModel(l.svcCtx.DbEngine).FindOneByPlatformId(channel.Id)
	if err != nil {
		l.Errorf("平台通道绑定的商户失败, channelId=%v, err=%v", channel.Id, err)
		return common.NewCodeError(common.InvalidUpstreamChannel)
	}

	// 6.1需要新增的商户上游通道绑定关系
	var upData []*model.MerchantChannelUpstream
	for _, m := range mChannels {
		for _, upId := range addUp {
			upData = append(upData, &model.MerchantChannelUpstream{
				MerchantId:        m.MerchantId,
				MerchantChannelId: m.Id,
				UpstreamChannelId: upId,
				Weight:            10,
			})
		}

	}

	txErr := l.svcCtx.DbEngine.Transaction(func(tx *gorm.DB) error {
		// 7.删除绑定关系
		if len(delUp) > 0 {
			for _, upChannelId := range delUp {
				// 7.1删除平台通道和上游通道绑定关系
				if err := model.NewPlatformChannelUpstreamModel(tx).DeleteByWhere(req.ChannelId, upChannelId); err != nil {
					l.Errorf("删除平台通道[%v]和上游通道[%v]绑定关系失败, err=%v", req.ChannelId, upChannelId, err)
					return err
				}

				// 7.2删除上游通道和商户通道的绑定关系
				for _, mc := range mChannels {
					if err := model.NewMerchantChannelUpstreamModel(tx).DeleteByUpChannel(mc.Id, upChannelId); err != nil {
						l.Errorf("删除商户通道[%v]和上游通道[%v]绑定关系失败, err=%v", mc.Id, upChannelId, err)
						return err
					}
				}
			}

		}

		// 8.添加绑定关系
		if len(data) > 0 {
			// 8.1添加平台通道和上游通道的绑定关系
			if err := model.NewPlatformChannelUpstreamModel(tx).Inserts(data); err != nil {
				l.Errorf("添加平台通道[%v]和上游通道[%+v]的绑定关系失败, err=%v", req.ChannelId, data, err)
				return err
			}

			if upData != nil {
				//8.2添加商户通道和上游通道的绑定关系
				if err := model.NewMerchantChannelUpstreamModel(tx).Inserts(upData); err != nil {
					l.Errorf("绑定关系商户通道和上游通道绑定关系失败, data:%+v err=%v", upData, err)
					return err
				}
			}

		}

		return nil
	})

	if txErr != nil {
		l.Errorf("保存通道[%v]的上游通道关联失败, err=%v", req.ChannelId, err)
		return common.NewCodeError(common.SysDBSave)
	}

	return nil
}
