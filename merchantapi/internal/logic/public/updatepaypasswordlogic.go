package public

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"
	"tpay_backend/model"

	"github.com/tal-tech/go-zero/core/logx"
)

type UpdatePayPassWordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	userId int64
}

func NewUpdatePayPassWordLogic(ctx context.Context, svcCtx *svc.ServiceContext, userId int64) UpdatePayPassWordLogic {
	return UpdatePayPassWordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		userId: userId,
	}
}

func (l *UpdatePayPassWordLogic) UpdatePayPassWord(req types.UpdatePayPassWordReq) error {
	//校验旧密码
	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneById(l.userId)
	if err != nil {
		logx.Errorf("修改支付密码失败,商户[%v]不存在,req[%v],err=[%v]", l.userId, req, err)
		return common.NewCodeError(common.UpdatePayPwdFailed)
	}

	//解密
	oldPlainPassword, err := common.DecryptPassword(req.OldPayPwd)
	if err != nil {
		l.Errorf("密码解密发生错误,err:%v, password:%v", err, req.OldPayPwd)
		return common.NewCodeError(common.SysDBErr)
	}

	// 比对密码
	if merchant.PayPassword != common.CreateMerchantPayPassword(oldPlainPassword) {
		l.Error("修改支付密码失败:支付密码校验失败")
		return common.NewCodeError(common.PayPasswordErr)
	}

	newPlainPassword, err := common.DecryptPassword(req.NewPayPwd)
	if err != nil {
		l.Errorf("密码解密发生错误,err:%v, password:%v", err, req.NewPayPwd)
		return common.NewCodeError(common.SysDBErr)
	}

	newPayPwdDb := common.CreateMerchantPayPassword(newPlainPassword)

	if err := model.NewMerchantModel(l.svcCtx.DbEngine).UpdatePayPwd(l.userId, newPayPwdDb); err != nil {
		logx.Error("修改商户[%v]支付密码失败, err[%v]", l.userId, err)
		return common.NewCodeError(common.UpdatePayPwdFailed)
	}

	return nil
}
