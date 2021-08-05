package public

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"
	"tpay_backend/model"

	"github.com/tal-tech/go-zero/core/logx"
)

type UpdateLoginPassWordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	userId int64
}

func NewUpdateLoginPassWordLogic(ctx context.Context, svcCtx *svc.ServiceContext, userId int64) UpdateLoginPassWordLogic {
	return UpdateLoginPassWordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		userId: userId,
	}
}

func (l *UpdateLoginPassWordLogic) UpdateLoginPassWord(req types.UpdateLoginPassWordReq) error {

	//校验旧密码
	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneById(l.userId)
	if err != nil {
		logx.Errorf("修改密码失败,商户[%v]不存在,req[%v],err=[%v]", l.userId, req, err)
		return common.NewCodeError(common.ReSetLoginPwdFailed)
	}

	oldPlainPassword, err := common.DecryptPassword(req.OldPwd)
	if err != nil {
		l.Errorf("密码解密发生错误,err:%v, password:%v", err, req.OldPwd)
		return common.NewCodeError(common.SysDBErr)
	}

	// 比对密码
	if merchant.Password != common.CreateMerchantPassword(oldPlainPassword) {
		l.Error("修改登录密码失败:登录密码校验失败")
		return common.NewCodeError(common.LoginPwdFailed)
	}

	//解密
	newPlainPassword, err := common.DecryptPassword(req.NewPwd)
	if err != nil {
		l.Errorf("密码解密发生错误,err:%v, password:%v", err, req.NewPwd)
		return common.NewCodeError(common.SysDBErr)
	}

	//产生新密码
	newPwd := common.CreateMerchantPassword(newPlainPassword)

	//更新为新密码
	if err := model.NewMerchantModel(l.svcCtx.DbEngine).UpdatePassword(l.userId, newPwd); err != nil {
		logx.Errorf("商户[%v]修改密码失败,req[%v],err=[%v]", l.userId, req, err)
		return common.NewCodeError(common.ReSetLoginPwdFailed)
	}

	//清除登录的reids token
	//if err := l.svcCtx.RedisSession.CleanOtherLogined(l.userId); err != nil {
	//	logx.Errorf("清除商户[%v]登录的redis失败,err=[%v]", l.userId, err)
	//}

	return nil
}
