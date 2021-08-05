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

type ModifyUpstreamChannelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModifyUpstreamChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) ModifyUpstreamChannelLogic {
	return ModifyUpstreamChannelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModifyUpstreamChannelLogic) ModifyUpstreamChannel(req types.ModifyUpstreamChannelRequest) error {
	channel, err := model.NewUpstreamChannelModel(l.svcCtx.DbEngine).FindOneById(req.ChannelId)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("上游通道[%v]不存在", req.ChannelId)
			return common.NewCodeError(common.ChannelNotExist)
		} else {
			l.Errorf("查询上游通道[%v]是否存在失败, err=%v", req.ChannelId, err)
			return common.NewCodeError(common.SysDBUpdate)
		}
	}

	if req.ChannelType != channel.ChannelType && req.ChannelType != 0 {
		isHave, err := model.NewPlatformChannelUpstreamModel(l.svcCtx.DbEngine).CheckUpstreamByIdAndType(req.ChannelId, channel.ChannelType)
		if err != nil {
			l.Errorf("查询通道[%v]是否被[%v]通道类型的下游渠道关联失败, err=%v", req.ChannelId, channel.ChannelType, err)
			return common.NewCodeError(common.SysDBUpdate)
		}

		if isHave {
			l.Errorf("通道[%v]已被[%v]通道类型的下游渠道关联, 不能修改通道类型", req.ChannelId, channel.ChannelType)
			return common.NewCodeError(common.ChannelHaveLinkedChannelNotUpdate)
		}
	}

	data := model.UpstreamChannel{
		ChannelName: req.ChannelName,
		ChannelCode: req.ChannelCode,
		ChannelDesc: req.ChannelDesc,
		Currency:    req.Currency,
		UpstreamId:  req.UpstreamId,
		ChannelType: req.ChannelType,
	}
	if err := model.NewUpstreamChannelModel(l.svcCtx.DbEngine).Update(req.ChannelId, data); err != nil {
		l.Errorf("修改上游通道[%v]失败, err=%v", req.ChannelId, err)
		errStr := err.Error()
		if strings.Contains(errStr, model.UpstreamChannelNameUniqueErr.Error()) || strings.Contains(errStr, model.UpstreamChannelCodeUniqueErr.Error()) {
			return common.NewCodeError(common.ChannelRepetition)
		} else {
			return common.NewCodeError(common.SysDBUpdate)
		}
	}

	return nil
}
