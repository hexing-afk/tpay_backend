package admin

import (
	"context"
	"errors"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"
	"tpay_backend/utils"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ResetAdminTotpSecretLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetAdminTotpSecretLogic(ctx context.Context, svcCtx *svc.ServiceContext) ResetAdminTotpSecretLogic {
	return ResetAdminTotpSecretLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetAdminTotpSecretLogic) ResetAdminTotpSecret(req types.ResetAdminTotpSecretRequest) (*types.ResetAdminTotpSecretResponse, error) {
	adminModel := model.NewAdminModel(l.svcCtx.DbEngine)

	admin, err := adminModel.FindOneById(req.AdminId)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			return nil, common.NewCodeError(common.UserNotExist)
		} else {
			l.Errorf("登录失败:查询数据库失败:%v", err)
			return nil, common.NewCodeError(common.SysDBErr)
		}
	}

	// 生成TOTP秘钥
	totpSecret, err := utils.GenerateTOTPSecret(admin.Username)
	if err != nil {
		l.Errorf("生成TOTP秘钥错误,err:%v, Username:%v", err, admin.Username)
		return nil, common.NewCodeError(common.SystemInternalErr)
	}

	admin.TotpSecret = totpSecret

	// 更新数据
	if err := adminModel.Save(admin); err != nil {
		l.Errorf("修改用户信息失败err:%v,admin:%+v", err, admin)
		return nil, common.NewCodeError(common.SysDBUpdate)
	}

	return &types.ResetAdminTotpSecretResponse{TotpSecret: totpSecret}, nil
}
