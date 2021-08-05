package public

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"
	"tpay_backend/model"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetBaseInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBaseInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetBaseInfoLogic {
	return GetBaseInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBaseInfoLogic) GetBaseInfo(userId int64) (*types.GetBaseInfoResponse, error) {
	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneById(userId)
	if err != nil {
		logx.Error("查询商户信息出错,err=[%v]", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}
	data := types.BaseInfoData{
		Username:   merchant.Username,
		MerchantId: merchant.Id,
		Phone:      merchant.Phone,
		Email:      merchant.Email,
		Md5Key:     merchant.Md5Key,
		MerchantNo: merchant.MerchantNo,
	}
	response := &types.GetBaseInfoResponse{
		Data: data,
	}

	return response, nil
}
