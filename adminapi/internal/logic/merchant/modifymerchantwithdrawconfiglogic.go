package merchant

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ModifyMerchantWithdrawConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModifyMerchantWithdrawConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) ModifyMerchantWithdrawConfigLogic {
	return ModifyMerchantWithdrawConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModifyMerchantWithdrawConfigLogic) ModifyMerchantWithdrawConfig(req types.ModifyMerchantWithdrawConfigRequest) error {
	data := model.MerchantWithdrawConfig{
		SingleMinAmount: req.SingleMinAmount,
		SingleMaxAmount: req.SingleMaxAmount,
		Rate:            req.WithdrawRate,
		SingleFee:       req.SingleFee,
	}
	if err := model.NewMerchantWithdrawConfigModel(l.svcCtx.DbEngine).Update(req.ConfigId, data); err != nil {
		l.Errorf("修改商户配置[%v]失败, data=%+v, err=%v", req.ConfigId, data, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	return nil
}
