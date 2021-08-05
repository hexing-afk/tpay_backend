package logic

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetOtherConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOtherConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetOtherConfigLogic {
	return GetOtherConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOtherConfigLogic) GetOtherConfig() (*types.GetOtherConfigResponse, error) {
	payTradeTypeSlice, err := model.NewGlobalConfigModel(l.svcCtx.DbEngine).FindValueByKey(model.ConfigPayTradeTypeSlice)
	if err != nil {
		logx.Errorf("查询代收交易类型出错,key:%v, err:%v", model.ConfigPayTradeTypeSlice, err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	transferTradeTypeSlice, err := model.NewGlobalConfigModel(l.svcCtx.DbEngine).FindValueByKey(model.ConfigTransferTradeTypeSlice)
	if err != nil {
		logx.Errorf("查询代付交易类型出错,key:%v, err:%v", model.ConfigTransferTradeTypeSlice, err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	return &types.GetOtherConfigResponse{
		PayTradeTypeSlice:      payTradeTypeSlice,
		TransferTradeTypeSlice: transferTradeTypeSlice,
	}, nil
}
