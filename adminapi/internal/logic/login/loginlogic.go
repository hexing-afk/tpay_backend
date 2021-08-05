package login

import (
	"context"
	"errors"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"
	"tpay_backend/model"
	"tpay_backend/utils"

	"github.com/tal-tech/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) LoginLogic {
	return LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req types.LoginReq) (*types.LoginReply, error) {
	/*
		// 校验登录验证码
		captcha := utils.NewCaptcha(l.svcCtx.Redis, utils.CaptchaConfig{
			KeyPrefix: common.CaptchaPrefix,
			Expire:    common.CaptchaExpire,
		})

		if ok, err := captcha.VerifyCaptcha(req.CaptchaId, req.CaptchaCode); err != nil {
			l.Errorf("校验登录验证码出错:err:%v", err)
			return nil, common.NewCodeError(common.LoginCaptchaNotMatch)
		} else {
			if !ok {
				return nil, common.NewCodeError(common.LoginCaptchaNotMatch)
			}
		}
	*/

	// 通过用户名查找用户
	admin, err := model.NewAdminModel(l.svcCtx.DbEngine).FindOneByUsername(req.Username)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			l.Errorf("登录失败:用户不存在:err--%v", err)
			return nil, common.NewCodeError(common.UserNotExist)
		} else {
			l.Errorf("登录失败:查询数据库失败:%v", err)
			return nil, common.NewCodeError(common.SysDBErr)
		}
	}

	// google验证码是否关闭
	totpIsClose, err := model.NewGlobalConfigModel(l.svcCtx.DbEngine).TotpIsClose()
	if err != nil {
		l.Errorf("查询totp配置失败:%v", err)
		return nil, common.NewCodeError(common.SysDBErr)
	}

	if !totpIsClose { // 没有关闭
		// 验证TOTP密码code
		if !utils.VerifyTOTPPasscode(req.TotpCode, admin.TotpSecret) {
			return nil, common.NewCodeError(common.LoginCaptchaNotMatch)
		}
	}

	//解密
	plainPassword, err := common.DecryptPassword(req.Password)
	if err != nil {
		l.Errorf("密码解密发生错误,err:%v, password:%v", err, req.Password)
		return nil, common.NewCodeError(common.SysDBErr)
	}
	//l.Infof("password:%v, plainPassword=[%v]", req.Password, plainPassword)

	// 比对密码
	if admin.Password != common.CreateAdminPassword(plainPassword) {
		l.Error("登录失败:密码错误")
		return nil, common.NewCodeError(common.GetLoginPasswordError)
	}

	if admin.EnableStatus != model.AdminEnableStatus {
		l.Errorf("登录失败:用户已被禁用")
		return nil, common.NewCodeError(common.AccountDisable)
	}

	//删除可能存在的旧token
	if err := l.svcCtx.RedisSession.CleanOtherLogined(admin.Id); err != nil {
		l.Error("登录失败:删除旧的redis失败:%v", err)
		return nil, common.NewCodeError(common.UserLoginFailed)
	}

	token, err := l.svcCtx.RedisSession.Login(admin.Id)
	if err != nil {
		l.Error("登录失败:存入redis失败:%v", err)
		return nil, common.NewCodeError(common.UserLoginFailed)
	}

	loginToken := common.LoginTokenGenerate(admin.Id, token)

	return &types.LoginReply{
		Username:   admin.Username,
		LoginToken: loginToken,
	}, nil
}
