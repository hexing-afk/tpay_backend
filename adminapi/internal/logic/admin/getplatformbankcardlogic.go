package admin

import (
	"context"
	"fmt"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetPlatformBankCardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPlatformBankCardLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPlatformBankCardLogic {
	return GetPlatformBankCardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPlatformBankCardLogic) GetPlatformBankCard(req types.GetPlatformBankCardRequest) (*types.GetPlatformBankCardResponse, error) {
	data, err := model.NewPlatformBankCardModel(l.svcCtx.DbEngine).FindOneById(req.CardId)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("卡[%v]不存在", req.CardId)
			return nil, common.NewCodeError(common.FindContentNotExist)
		} else {
			l.Errorf("查询卡[%v]失败, err=%v", req.CardId, err)
			return nil, common.NewCodeError(common.SysDBGet)
		}
	}

	imageUrl, err := model.NewGlobalConfigModel(l.svcCtx.DbEngine).FindValueByKey(model.ConfigImageBaseUrl)
	if err != nil {
		l.Errorf("查询图片域名地址失败[%v], err=%v", model.ConfigImageBaseUrl, err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	card := types.PlatformBankCard{
		CardId:      data.Id,
		BankName:    data.BankName,
		AccountName: data.AccountName,
		CardNumber:  data.CardNumber,
		BranchName:  data.BranchName,
		Currency:    data.Currency,
		MaxAmount:   data.MaxAmount,
		Remark:      data.Remark,
		QrCode:      data.QrCode,
	}

	if card.QrCode != "" {
		card.QrCode = fmt.Sprintf("%s/%s", imageUrl, data.QrCode)
	}

	return &types.GetPlatformBankCardResponse{
		Card: card,
	}, nil
}
