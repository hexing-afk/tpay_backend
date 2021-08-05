package merchant

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ModifyMerchantLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModifyMerchantLogic(ctx context.Context, svcCtx *svc.ServiceContext) ModifyMerchantLogic {
	return ModifyMerchantLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModifyMerchantLogic) ModifyMerchant(req types.ModifyMerchantRequest) error {
	err := model.NewMerchantModel(l.svcCtx.DbEngine).Update(req.MerchantId, model.Merchant{
		Phone: req.Phone,
		Email: req.Email,
	})
	if err != nil {
		l.Errorf("修改商户[%v]信息失败, err=%v", req.MerchantId, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	return nil
}
