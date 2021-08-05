package merchant

import (
	"context"
	"gorm.io/gorm"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type DeleteMerchantChannelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteMerchantChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) DeleteMerchantChannelLogic {
	return DeleteMerchantChannelLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMerchantChannelLogic) DeleteMerchantChannel(req types.DeleteMerchantChannelRequest) error {
	txErr := l.svcCtx.DbEngine.Transaction(func(tx *gorm.DB) error {
		if err := model.NewMerchantChannelModel(tx).Delete(req.ChannelId); err != nil {
			l.Errorf("删除商户通道[%v]失败, err=%v", req.ChannelId, err)
			return err
		}

		if err := model.NewMerchantChannelUpstreamModel(tx).Delete(req.ChannelId); err != nil {
			l.Errorf("删除商户通道[%v]上游失败, err=%v", req.ChannelId, err)
			return err
		}

		return nil
	})

	if txErr != nil {
		l.Errorf("删除商户通道[%v]失败, err=%v", req.ChannelId, txErr)
		return common.NewCodeError(common.SysDBDelete)
	}

	return nil
}
