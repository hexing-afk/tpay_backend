package merchant

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type SaveMerchantChannelUpstreamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveMerchantChannelUpstreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) SaveMerchantChannelUpstreamLogic {
	return SaveMerchantChannelUpstreamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveMerchantChannelUpstreamLogic) SaveMerchantChannelUpstream(req types.SaveMerchantChannelUpstreamRequest) error {
	if err := model.NewMerchantChannelUpstreamModel(l.svcCtx.DbEngine).UpdateWeight(req.Id, req.UpstreamChannelWeight); err != nil {
		l.Errorf("更新商户上游通道权重配置失败, err=%v", err)
		return common.NewCodeError(common.SysDBSave)
	}
	return nil
}
