package public

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type UpdateBaseInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	userId int64
}

func NewUpdateBaseInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext, userId int64) UpdateBaseInfoLogic {
	return UpdateBaseInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		userId: userId,
	}
}

func (l *UpdateBaseInfoLogic) UpdateBaseInfo(req types.UpdateBaseInfoReq) error {
	err := model.NewMerchantModel(l.svcCtx.DbEngine).Update(l.userId, model.Merchant{
		Phone: req.Phone,
		Email: req.Email,
	})
	if err != nil {
		l.Errorf("修改商户[%v]信息失败, err=%v", l.userId, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	return nil
}
