package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetPlatformBankCardListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPlatformBankCardListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPlatformBankCardListLogic {
	return GetPlatformBankCardListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPlatformBankCardListLogic) GetPlatformBankCardList(req types.GetPlatformBankCardListRequest) (*types.GetPlatformBankCardListResponse, error) {
	f := model.FindPlatformBankCardList{
		Search:   req.Search,
		Currency: req.Currency,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	list, total, err := model.NewPlatformBankCardModel(l.svcCtx.DbEngine).FindList(f)
	if err != nil {
		l.Errorf("查询平台收款卡列表失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var cardList []types.PlatformBankCardList
	for _, v := range list {
		cardList = append(cardList, types.PlatformBankCardList{
			CardId:        v.Id,
			BankName:      v.BankName,
			AccountName:   v.AccountName,
			CreateTime:    v.CreateTime,
			CardNumber:    v.CardNumber,
			BranchName:    v.BranchName,
			Currency:      v.Currency,
			MaxAmount:     v.MaxAmount,
			Remark:        v.Remark,
			ReceivedToday: v.TodayReceived,
			Status:        v.Status,
		})
	}

	return &types.GetPlatformBankCardListResponse{
		Total: total,
		List:  cardList,
	}, nil
}
