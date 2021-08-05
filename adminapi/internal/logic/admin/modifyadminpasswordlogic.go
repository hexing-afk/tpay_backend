package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ModifyAdminPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModifyAdminPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) ModifyAdminPasswordLogic {
	return ModifyAdminPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModifyAdminPasswordLogic) ModifyAdminPassword(adminId int64, req types.ModifyAdminPasswordRequest) error {
	admin, err := model.NewAdminModel(l.svcCtx.DbEngine).FindOneById(adminId)
	if err != nil {
		l.Logger.Errorf("查询管理员[%v]失败, err=%v", adminId, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	//解密
	oldPlainPassword, err := common.DecryptPassword(req.OldPassword)
	if err != nil {
		l.Errorf("密码解密发生错误,err:%v, password:%v", err, req.OldPassword)
		return common.NewCodeError(common.SysDBErr)
	}

	if admin.Password != common.CreateAdminPassword(oldPlainPassword) {
		l.Logger.Errorf("管理员[%v]密码错误", adminId, err)
		return common.NewCodeError(common.GetLoginPasswordError)
	}

	newPlainPassword, err := common.DecryptPassword(req.NewPassword)
	if err != nil {
		l.Errorf("密码解密发生错误,err:%v, password:%v", err, req.NewPassword)
		return common.NewCodeError(common.SysDBErr)
	}

	if err := model.NewAdminModel(l.svcCtx.DbEngine).UpdatePassword(adminId, common.CreateAdminPassword(newPlainPassword)); err != nil {
		l.Logger.Errorf("修改管理员[%v]密码失败, err=%v", adminId, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	//清除管理员的登录token
	if errDel := l.svcCtx.RedisSession.CleanOtherLogined(adminId); errDel != nil {
		l.Errorf("删除管理员[%v]的登录token失败", adminId)
		return common.NewCodeError(common.SystemInternalErr)
	}

	return nil
}
