package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type EnablePlatformChannelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEnablePlatformChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) EnablePlatformChannelLogic {
	return EnablePlatformChannelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EnablePlatformChannelLogic) EnablePlatformChannel(req types.EnablePlatformChannelRequest) error {
	var err error
	switch req.Enable {
	case model.PlatformChannelStatusEnable:
		err = model.NewPlatformChannelModel(l.svcCtx.DbEngine).EnableChannel(req.ChannelId)
	case model.PlatformChannelStatusDisable:
		err = model.NewPlatformChannelModel(l.svcCtx.DbEngine).DisableChannel(req.ChannelId)
	default:
		l.Errorf("不支持该操作[%v]", req.Enable)
		return common.NewCodeError(common.InvalidParam)
	}
	if err != nil {
		l.Errorf("启用|禁用通道[%v]失败, err=%v", req.ChannelId, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	return nil
}
