package merchant

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type EnableMerchantChannelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEnableMerchantChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) EnableMerchantChannelLogic {
	return EnableMerchantChannelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EnableMerchantChannelLogic) EnableMerchantChannel(req types.EnableMerchantChannelRequest) error {
	var err error
	switch req.Enable {
	case model.MerchantChannelStatusEnable:
		err = model.NewMerchantChannelModel(l.svcCtx.DbEngine).EnableMerchantChannel(req.ChannelId)
	case model.MerchantChannelStatusDisable:
		err = model.NewMerchantChannelModel(l.svcCtx.DbEngine).DisableMerchantChannel(req.ChannelId)
	default:
		l.Errorf("操作不支持, req.Enable=%v", req.Enable)
		return common.NewCodeError(common.InvalidParam)
	}

	if err != nil {
		l.Errorf("启用|禁用商户通道失败, err=%v", err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	return nil
}
