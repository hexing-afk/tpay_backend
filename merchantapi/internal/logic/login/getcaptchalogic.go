package login

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/utils"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetCaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetCaptchaLogic {
	return GetCaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCaptchaLogic) GetCaptcha() (*types.CaptchaReply, error) {
	captcha := utils.NewCaptcha(l.svcCtx.Redis, utils.CaptchaConfig{
		KeyPrefix: common.CaptchaPrefix,
		Expire:    common.CaptchaExpire,
	})

	id, b64s, err := captcha.GenerateCaptcha()
	if err != nil {
		l.Errorf("获取登录验证码失败,err:%v", err)
		return nil, common.NewCodeError(common.GetLoginCaptchaFailed)
	}
	return &types.CaptchaReply{CaptchaId: id, Base64png: b64s}, nil
}
