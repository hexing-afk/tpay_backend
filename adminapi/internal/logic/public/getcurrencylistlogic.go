package logic

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetCurrencyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCurrencyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetCurrencyListLogic {
	return GetCurrencyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCurrencyListLogic) GetCurrencyList() (*types.GetCurrencyListResponse, error) {
	list, err := model.NewCurrencyModel(l.svcCtx.DbEngine).FindMany()
	if err != nil {
		l.Errorf("查询币种列表失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var currency []types.Currency
	for _, v := range list {
		currency = append(currency, types.Currency{
			Currency:        v.Currency,
			Symbol:          v.Symbol,
			Country:         v.Country,
			IsDivideHundred: v.IsDivideHundred,
		})
	}

	return &types.GetCurrencyListResponse{List: currency}, nil
}
