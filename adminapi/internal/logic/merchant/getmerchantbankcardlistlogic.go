package merchant

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetMerchantBankCardListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMerchantBankCardListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetMerchantBankCardListLogic {
	return GetMerchantBankCardListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMerchantBankCardListLogic) GetMerchantBankCardList(req types.GetMerchantBankCardListRequest) (*types.GetMerchantBankCardListResponse, error) {
	f := model.FindMerchantBankCardList{
		Search:   req.Search,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	list, total, err := model.NewMerchantBankCardModel(l.svcCtx.DbEngine).FindList(f)
	if err != nil {
		l.Errorf("查询商户银行卡列表失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var cardList []types.MerchantBankCard
	for _, v := range list {
		cardList = append(cardList, types.MerchantBankCard{
			CardId:      v.Id,
			Username:    v.MerchantUsername,
			BankName:    v.BankName,
			AccountName: v.AccountName,
			CardNumber:  v.CardNumber,
			BranchName:  v.BranchName,
			Currency:    v.Currency,
			Remark:      v.Remark,
			CreateTime:  v.CreateTime,
			UpdateTime:  v.UpdateTime,
		})
	}

	return &types.GetMerchantBankCardListResponse{
		Total: total,
		List:  cardList,
	}, nil
}
