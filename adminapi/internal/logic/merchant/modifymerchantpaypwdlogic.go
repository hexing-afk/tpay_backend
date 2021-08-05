package merchant

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ModifyMerchantPayPwdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModifyMerchantPayPwdLogic(ctx context.Context, svcCtx *svc.ServiceContext) ModifyMerchantPayPwdLogic {
	return ModifyMerchantPayPwdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModifyMerchantPayPwdLogic) ModifyMerchantPayPwd(adminId int64, req types.ModifyMerchantPayPwdRequest) error {
	admin, err := model.NewAdminModel(l.svcCtx.DbEngine).FindOneById(adminId)
	if err != nil {
		l.Errorf("查询管理员[%v]密码失败, err=%v", adminId, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	//解密
	plainPassword, err := common.DecryptPassword(req.LoginUserPwd)
	if err != nil {
		l.Errorf("密码解密发生错误,err:%v, password:%v", err, req.LoginUserPwd)
		return common.NewCodeError(common.SysDBUpdate)
	}

	if admin.Password != common.CreateAdminPassword(plainPassword) {
		l.Errorf("管理员[%v]密码错误", adminId)
		return common.NewCodeError(common.GetLoginPasswordError)
	}

	//解密
	payPassword, err := common.DecryptPassword(req.PayPassword)
	if err != nil {
		l.Errorf("密码解密发生错误,err:%v, password:%v", err, req.LoginUserPwd)
		return common.NewCodeError(common.SysDBUpdate)
	}

	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneById(req.MerchantId)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("商户[%v]不存在", req.MerchantId)
			return common.NewCodeError(common.MerchantNotExist)
		} else {
			l.Errorf("查询商户[%v]信息失败, err:%v", req.MerchantId, err)
			return common.NewCodeError(common.SysDBUpdate)
		}
	}

	// 修改商户支付密码
	mchPayPwd := common.CreateMerchantPayPassword(payPassword)
	if err := model.NewMerchantModel(l.svcCtx.DbEngine).UpdatePayPwd(merchant.Id, mchPayPwd); err != nil {
		l.Errorf("修改商户[%v]支付密码失败, password:%v, err:%v", merchant.Id, mchPayPwd, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	return nil
}
