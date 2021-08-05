package public

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetWithdrawConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWithdrawConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetWithdrawConfigLogic {
	return GetWithdrawConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWithdrawConfigLogic) GetWithdrawConfig(merchantId int64) (*types.GetWithdrawConfigResponse, error) {
	data, err := model.NewMerchantWithdrawConfigModel(l.svcCtx.DbEngine).FindOneByMerchantId(merchantId)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("商户[%v]没有提现配置", merchantId)
			return nil, common.NewCodeError(common.NotWithdrawConfig)
		} else {
			l.Errorf("查询商户[%v]提现配置失败, err=%v", merchantId, err)
			return nil, common.NewCodeError(common.SysDBGet)
		}
	}

	return &types.GetWithdrawConfigResponse{
		Config: types.WithdrawConfig{
			SingleMinAmount: data.SingleMinAmount,
			SingleMaxAmount: data.SingleMaxAmount,
			DeductionMethod: data.DeductionMethod,
			Rate:            data.Rate,
			SingleFee:       data.SingleFee,
		},
	}, nil
}
