package merchant

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

type ResetMerchantTotpSecretLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetMerchantTotpSecretLogic(ctx context.Context, svcCtx *svc.ServiceContext) ResetMerchantTotpSecretLogic {
	return ResetMerchantTotpSecretLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetMerchantTotpSecretLogic) ResetMerchantTotpSecret(req types.ResetMerchantTotpSecretRequest) (*types.ResetMerchantTotpSecretResponse, error) {
	merchantModel := model.NewMerchantModel(l.svcCtx.DbEngine)

	merchant, err := merchantModel.FindOneById(req.MerchantId)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			return nil, common.NewCodeError(common.MerchantNotExist)
		} else {
			l.Errorf("登录失败:查询数据库失败:%v", err)
			return nil, common.NewCodeError(common.SysDBErr)
		}
	}

	// 生成TOTP秘钥
	totpSecret, err := utils.GenerateTOTPSecret(merchant.Username)
	if err != nil {
		l.Errorf("生成TOTP秘钥错误,err:%v, Username:%v", err, merchant.Username)
		return nil, common.NewCodeError(common.SystemInternalErr)
	}

	merchant.TotpSecret = totpSecret

	// 更新数据
	if err := merchantModel.Save(merchant); err != nil {
		l.Errorf("修改商户信息失败err:%v,merchant:%+v", err, merchant)
		return nil, common.NewCodeError(common.SysDBUpdate)
	}

	return &types.ResetMerchantTotpSecretResponse{TotpSecret: totpSecret}, nil
}
