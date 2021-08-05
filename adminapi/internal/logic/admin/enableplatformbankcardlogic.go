package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"
	"tpay_backend/model"

	"github.com/tal-tech/go-zero/core/logx"
)

type EnablePlatformBankCardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEnablePlatformBankCardLogic(ctx context.Context, svcCtx *svc.ServiceContext) EnablePlatformBankCardLogic {
	return EnablePlatformBankCardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EnablePlatformBankCardLogic) EnablePlatformBankCard(req types.EnablePlatformBankCardRequest) error {
	var err error
	switch req.Enable {
	case model.PlatformBankCardEnable:
		err = model.NewPlatformBankCardModel(l.svcCtx.DbEngine).EnableCard(req.CardId)
	case model.PlatformBankCardDisable:
		err = model.NewPlatformBankCardModel(l.svcCtx.DbEngine).DisableCard(req.CardId)
	default:
		l.Errorf("req.Enable参数的值有误, %v", req.Enable)
		return common.NewCodeError(common.InvalidParam)
	}

	if err != nil {
		l.Errorf("启用|禁用银行卡[%v]失败, err=%v", req.CardId, req.Enable)
		return common.NewCodeError(common.SysDBUpdate)
	}

	return nil
}
