package admin

import (
	"context"
	"strings"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ModifyPlatformChannelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModifyPlatformChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) ModifyPlatformChannelLogic {
	return ModifyPlatformChannelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModifyPlatformChannelLogic) ModifyPlatformChannel(req types.ModifyPlatformChannelRequest) error {
	channel, err := model.NewPlatformChannelModel(l.svcCtx.DbEngine).FindOneById(req.ChannelId)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("通道[%v]不存在", req.ChannelId)
			return common.NewCodeError(common.UpdateContentNotExist)
		} else {
			l.Errorf("查询通道[%v]失败, err=%v", req.ChannelId, err)
			return common.NewCodeError(common.SysDBUpdate)
		}
	}

	if req.ChannelType != channel.ChannelType && req.ChannelType != 0 {
		isHave, err := model.NewPlatformChannelUpstreamModel(l.svcCtx.DbEngine).CheckPlatformByIdAndType(req.ChannelId, channel.ChannelType)
		if err != nil {
			l.Errorf("查询通道[%v]是否已有[%v]通道类型的上游渠道失败, err=%v", req.ChannelId, channel.ChannelType, err)
			return common.NewCodeError(common.SysDBUpdate)
		}

		if isHave {
			l.Errorf("通道[%v]已有[%v]通道类型的上游渠道, 不能修改通道类型", req.ChannelId, channel.ChannelType)
			return common.NewCodeError(common.ChannelHaveLinkedChannelNotUpdate)
		}
	}

	b, err := model.NewAreaModel(l.svcCtx.DbEngine).Check(req.AreaId)
	if err != nil {
		l.Errorf("确认地区是否存在出错, req=%v, err=%v", req, err)
		return common.NewCodeError(common.AreaNotExist)
	}

	if !b {
		l.Errorf("地区不存在 req=%v", req)
		return common.NewCodeError(common.AreaNotExist)
	}

	data := model.PlatformChannel{
		ChannelName: req.ChannelName,
		ChannelCode: req.ChannelCode,
		ChannelDesc: req.ChannelDesc,
		ChannelType: req.ChannelType,
		AreaId:      req.AreaId,
	}
	if err := model.NewPlatformChannelModel(l.svcCtx.DbEngine).Update(req.ChannelId, data); err != nil {
		l.Errorf("添加通道[%v]失败, err=%v", req.ChannelName, err)
		errStr := err.Error()
		if strings.Contains(errStr, model.PlatformChannelNameUniqueErr.Error()) || strings.Contains(errStr, model.PlatformChannelCodeUniqueErr.Error()) {
			return common.NewCodeError(common.ChannelRepetition)
		} else {
			return common.NewCodeError(common.SysDBUpdate)
		}
	}

	return nil
}
