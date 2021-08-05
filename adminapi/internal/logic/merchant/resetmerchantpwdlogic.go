package merchant

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ResetMerchantPwdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetMerchantPwdLogic(ctx context.Context, svcCtx *svc.ServiceContext) ResetMerchantPwdLogic {
	return ResetMerchantPwdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetMerchantPwdLogic) ResetMerchantPwd(adminId int64, req types.ResetMerchantPwdRequest) error {
	admin, err := model.NewAdminModel(l.svcCtx.DbEngine).FindOneById(adminId)
	if err != nil {
		l.Errorf("查询管理员[%v]密码失败, err=%v", adminId, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	//解密
	plainPassword, err := common.DecryptPassword(req.LoginUserPwd)
	if err != nil {
		l.Errorf("密码解密发生错误,err:%v, password:%v", err, req.LoginUserPwd)
		return common.NewCodeError(common.SysDBErr)
	}

	if admin.Password != common.CreateAdminPassword(plainPassword) {
		l.Errorf("管理员[%v]密码错误", adminId)
		return common.NewCodeError(common.GetLoginPasswordError)
	}

	newPassword := common.CreateAdminPassword(common.MerchantDefaultPassword)
	err = model.NewMerchantModel(l.svcCtx.DbEngine).UpdatePassword(req.MerchantId, newPassword)
	if err != nil {
		l.Errorf("充值商户[%v]登录密码失败, err=%v", req.MerchantId, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	return nil
}
