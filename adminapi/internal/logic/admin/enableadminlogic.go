package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type EnableAdminLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEnableAdminLogic(ctx context.Context, svcCtx *svc.ServiceContext) EnableAdminLogic {
	return EnableAdminLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EnableAdminLogic) EnableAdmin(req types.EnableAdminRequest) error {
	exist, err := model.NewAdminModel(l.svcCtx.DbEngine).CheckById(req.AdminId)
	if err != nil {
		l.Logger.Errorf("查询管理员[%v]是否存在失败, err=%v", req.AdminId, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	if !exist {
		l.Logger.Errorf("管理员[%v]不存在")
		return common.NewCodeError(common.UserNotExist)
	}

	switch req.Enable {
	case model.AdminEnableStatus:
		err = model.NewAdminModel(l.svcCtx.DbEngine).Enable(req.AdminId)
		if err != nil {
			l.Logger.Errorf("修改管理员[%v]账号开启状态失败, err=%v", req.AdminId, err)
			return common.NewCodeError(common.SysDBUpdate)
		}

	case model.AdminDisableStatus:
		err = model.NewAdminModel(l.svcCtx.DbEngine).Disable(req.AdminId)
		if err != nil {
			l.Logger.Errorf("修改管理员[%v]账号开启状态失败, err=%v", req.AdminId, err)
			return common.NewCodeError(common.SysDBUpdate)
		}

		//清除管理员的登录token
		if errDel := l.svcCtx.RedisSession.CleanOtherLogined(req.AdminId); errDel != nil {
			l.Errorf("删除管理员[%v]的登录token失败", req.AdminId)
			return common.NewCodeError(common.SystemInternalErr)
		}

	default:
		l.Logger.Errorf("无效参数, req.enable=%v", req.Enable)
		return common.NewCodeError(common.InvalidParam)
	}

	return nil
}
