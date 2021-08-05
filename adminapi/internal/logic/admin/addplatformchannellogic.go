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

type AddPlatformChannelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddPlatformChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddPlatformChannelLogic {
	return AddPlatformChannelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddPlatformChannelLogic) AddPlatformChannel(req types.AddPlatformChannelRequest) error {

	b, err := model.NewAreaModel(l.svcCtx.DbEngine).Check(req.AreaId)
	if err != nil {
		l.Errorf("确认地区是否存在出错, req=%v, err=%v", req, err)
		return common.NewCodeError(common.AreaNotExist)
	}

	if !b {
		l.Errorf("地区不存在 req=%v", req)
		return common.NewCodeError(common.AreaNotExist)
	}

	data := &model.PlatformChannel{
		ChannelName: req.ChannelName,
		ChannelCode: req.ChannelCode,
		ChannelDesc: req.ChannelDesc,
		Status:      model.PlatformChannelStatusEnable,
		ChannelType: req.ChannelType,
		AreaId:      req.AreaId,
	}

	if err := model.NewPlatformChannelModel(l.svcCtx.DbEngine).Insert(data); err != nil {
		l.Errorf("添加通道[%v]失败, err=%v", req.ChannelName, err)
		errStr := err.Error()
		if strings.Contains(errStr, model.PlatformChannelNameUniqueErr.Error()) || strings.Contains(errStr, model.PlatformChannelCodeUniqueErr.Error()) {
			return common.NewCodeError(common.ChannelRepetition)
		} else {
			return common.NewCodeError(common.SysDBAdd)
		}
	}

	return nil
}
