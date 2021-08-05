package merchant

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ModifyMerchantChannelRateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModifyMerchantChannelRateLogic(ctx context.Context, svcCtx *svc.ServiceContext) ModifyMerchantChannelRateLogic {
	return ModifyMerchantChannelRateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModifyMerchantChannelRateLogic) ModifyMerchantChannelRate(req types.ModifyMerchantChannelRateRequest) error {
	data := model.MerchantChannel{
		Rate:      req.ChannelRate,
		SingleFee: req.SingleFee,
	}
	if err := model.NewMerchantChannelModel(l.svcCtx.DbEngine).UpdateRate(req.ChannelId, data); err != nil {
		l.Errorf("修改商户通道[%v]费率失败, err=%v", req.ChannelId, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	return nil
}
