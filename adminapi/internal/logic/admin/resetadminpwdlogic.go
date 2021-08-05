package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ResetAdminPwdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetAdminPwdLogic(ctx context.Context, svcCtx *svc.ServiceContext) ResetAdminPwdLogic {
	return ResetAdminPwdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetAdminPwdLogic) ResetAdminPwd(adminId int64, req types.ResetAdminPwdRequest) error {
	//查询被重置的管理员是否存在
	exist, err := model.NewAdminModel(l.svcCtx.DbEngine).CheckById(req.AdminId)
	if err != nil {
		l.Logger.Errorf("查询管理员[%v]是否存在失败, err=%v", req.AdminId, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	if !exist {
		l.Logger.Errorf("管理员[%v]不存在")
		return common.NewCodeError(common.UserNotExist)
	}

	//检查当前登录的管理员密码
	admin, err := model.NewAdminModel(l.svcCtx.DbEngine).FindOneById(adminId)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Logger.Errorf("管理员[%v]不存在")
			return common.NewCodeError(common.UserNotExist)
		} else {
			l.Logger.Errorf("查询管理员[%v]是否存在失败, err=%v", req.AdminId, err)
			return common.NewCodeError(common.SysDBUpdate)
		}
	}

	//解密
	plainPassword, err := common.DecryptPassword(req.LoginUserPwd)
	if err != nil {
		l.Errorf("密码解密发生错误,err:%v, password:%v", err, req.LoginUserPwd)
		return common.NewCodeError(common.SysDBErr)
	}

	if admin.Password != common.CreateAdminPassword(plainPassword) {
		l.Logger.Errorf("当前登录的管理员[%v]账号密码错误", adminId)
		return common.NewCodeError(common.GetLoginPasswordError)
	}

	//设置新密码
	newPwd := common.CreateAdminPassword(common.AdminDefaultPassword)
	err = model.NewAdminModel(l.svcCtx.DbEngine).UpdatePassword(req.AdminId, newPwd)
	if err != nil {
		l.Logger.Errorf("重置管理员[%v]登录密码失败, err=%v", req.AdminId, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	return nil
}
