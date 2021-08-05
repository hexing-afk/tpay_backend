package public

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"
	"tpay_backend/utils"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type UpdateMd5KeyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	userId int64
}

func NewUpdateMd5KeyLogic(ctx context.Context, svcCtx *svc.ServiceContext, userId int64) UpdateMd5KeyLogic {
	return UpdateMd5KeyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		userId: userId,
	}
}

func (l *UpdateMd5KeyLogic) UpdateMd5Key() (*types.Md5KeyResponse, error) {

	newMd5Key := utils.RandString(32)
	if err := model.NewMerchantModel(l.svcCtx.DbEngine).UpdateMd5Key(l.userId, newMd5Key); err != nil {
		logx.Error("刷新商户[%v]的Md5Key失败, err[%v]", l.userId, err)
		return nil, common.NewCodeError(common.ReSetMd5KeyFailed)
	}

	return &types.Md5KeyResponse{
		Md5Key: newMd5Key,
	}, nil
}
