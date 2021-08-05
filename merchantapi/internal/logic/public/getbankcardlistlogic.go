package public

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetBankCardListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBankCardListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetBankCardListLogic {
	return GetBankCardListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBankCardListLogic) GetBankCardList(merchantId int64, req types.GetBankCardListRequest) (*types.GetBankCardListResponse, error) {
	f := model.FindMerchantBankCardList{
		Search:     req.Search,
		MerchantId: merchantId,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}
	data, total, err := model.NewMerchantBankCardModel(l.svcCtx.DbEngine).FindList(f)
	if err != nil {
		l.Errorf("查询商户[%v]银行卡列表失败, err=%v", merchantId, err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var list []types.BankCardList
	for _, v := range data {
		list = append(list, types.BankCardList{
			CardId:      v.Id,
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

	return &types.GetBankCardListResponse{
		Total: total,
		List:  list,
	}, nil
}
