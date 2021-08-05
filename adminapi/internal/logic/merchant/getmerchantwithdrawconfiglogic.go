package merchant

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetMerchantWithdrawConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMerchantWithdrawConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetMerchantWithdrawConfigLogic {
	return GetMerchantWithdrawConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMerchantWithdrawConfigLogic) GetMerchantWithdrawConfig(req types.GetMerchantWithdrawConfigRequest) (*types.GetMerchantWithdrawConfigResponse, error) {
	data, err := model.NewMerchantWithdrawConfigModel(l.svcCtx.DbEngine).FindOneByMerchantId(req.MerchantId)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("商户[%v]还没有配提现配置", req.MerchantId)
			return nil, nil
		} else {
			l.Errorf("查询商户[%v]提现配置失败, err=%v", req.MerchantId, err)
			return nil, common.NewCodeError(common.SysDBGet)
		}
	}

	config := types.MerchantWithdrawConfig{
		ConfigId:        data.Id,
		SingleMinAmount: data.SingleMinAmount,
		SingleMaxAmount: data.SingleMaxAmount,
		DeductionMethod: data.DeductionMethod,
		WithdrawRate:    data.Rate,
		SingleFee:       data.SingleFee,
	}

	return &types.GetMerchantWithdrawConfigResponse{
		Config: config,
	}, nil
}
