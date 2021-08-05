package merchant

import (
	"context"
	"errors"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type UpdateMerchantIpWhitelistLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateMerchantIpWhitelistLogic(ctx context.Context, svcCtx *svc.ServiceContext) UpdateMerchantIpWhitelistLogic {
	return UpdateMerchantIpWhitelistLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMerchantIpWhitelistLogic) UpdateMerchantIpWhitelist(req types.UpdateMerchantIpWhitelistRequest) error {
	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneById(req.MerchantId)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			return common.NewCodeError(common.MerchantNotExist)
		} else {
			l.Errorf("登录失败:查询数据库失败:%v", err)
			return common.NewCodeError(common.SysDBErr)
		}
	}

	merchant.IpWhiteList = req.IpWhitelist

	// 更新数据
	updateErr := model.NewMerchantModel(l.svcCtx.DbEngine).Save(merchant)
	if updateErr != nil {
		l.Errorf("修改商户信息失败err:%v,merchant:%+v", updateErr, merchant)
		return common.NewCodeError(common.SysDBUpdate)
	}

	return nil
}
