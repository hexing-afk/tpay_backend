package merchant

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddMerchantWithdrawConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddMerchantWithdrawConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddMerchantWithdrawConfigLogic {
	return AddMerchantWithdrawConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddMerchantWithdrawConfigLogic) AddMerchantWithdrawConfig(req types.AddMerchantWithdrawConfigRequest) (*types.AddMerchantWithdrawConfigResponse, error) {
	isHave, err := model.NewMerchantWithdrawConfigModel(l.svcCtx.DbEngine).CheckByMerchantId(req.MerchantId)
	if err != nil {
		l.Errorf("检查商户[%v]是否有提现配置失败, err=%v", req.MerchantId, err)
		return nil, common.NewCodeError(common.SysDBAdd)
	}

	if isHave {
		l.Errorf("商户[%v]有提现配置了", req.MerchantId)
		return nil, common.NewCodeError(common.ConfigAlreadyExist)
	}

	data := &model.MerchantWithdrawConfig{
		MerchantId:      req.MerchantId,
		SingleMinAmount: req.SingleMinAmount,
		SingleMaxAmount: req.SingleMaxAmount,
		DeductionMethod: model.MerchantWithdrawDeductionOut,
		Rate:            req.WithdrawRate,
		SingleFee:       req.SingleFee,
	}
	if err := model.NewMerchantWithdrawConfigModel(l.svcCtx.DbEngine).Insert(data); err != nil {
		l.Errorf("添加商户[%v]提现配置失败, err=%v", req.MerchantId, err)
		return nil, common.NewCodeError(common.SysDBAdd)
	}

	return &types.AddMerchantWithdrawConfigResponse{ConfigId: data.Id}, nil
}
