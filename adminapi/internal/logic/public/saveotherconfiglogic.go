package logic

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type SaveOtherConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveOtherConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) SaveOtherConfigLogic {
	return SaveOtherConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveOtherConfigLogic) SaveOtherConfig(req types.SaveOtherConfigRequest) error {
	//代收交易方式
	payTradeTypeSlice := model.GlobalConfig{
		ConfigKey:   model.ConfigPayTradeTypeSlice,
		ConfigValue: req.PayTradeTypeSlice,
		IsChange:    model.IsChangeFalse,
	}
	if err := model.NewGlobalConfigModel(l.svcCtx.DbEngine).InsertOrUpdate(payTradeTypeSlice); err != nil {
		l.Errorf("修改代收交易方式失败，err=%v", err)
		return common.NewCodeError(common.SysDBSave)
	}

	//代付交易类型
	transferTradeTypeSlice := model.GlobalConfig{
		ConfigKey:   model.ConfigTransferTradeTypeSlice,
		ConfigValue: req.TransferTradeTypeSlice,
		IsChange:    model.IsChangeFalse,
	}
	if err := model.NewGlobalConfigModel(l.svcCtx.DbEngine).InsertOrUpdate(transferTradeTypeSlice); err != nil {
		l.Errorf("修改代付交易方式失败，err=%v", err)
		return common.NewCodeError(common.SysDBSave)
	}

	return nil
}
