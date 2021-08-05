package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"
	"tpay_backend/upstream"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetUpstreamBalanceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUpstreamBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetUpstreamBalanceLogic {
	return GetUpstreamBalanceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUpstreamBalanceLogic) GetUpstreamBalance(req types.GetUpstreamBalanceRequest) (*types.GetUpstreamBalanceResponse, error) {
	up, err := model.NewUpstreamModel(l.svcCtx.DbEngine).FindOneById(req.UpstreamId)
	if err != nil {
		l.Errorf("查询上游[%v]失败, err=%v", req.UpstreamId, err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	if up.UpstreamMerchantNo == "" || up.CallConfig == "" {
		l.Errorf("上游[%v]没有配通信配置", req.UpstreamId)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var upstreamObj upstream.Upstream
	switch up.UpstreamCode {
	case upstream.UpstreamCodeTotopay:
		upstreamObj, err = upstream.NewTotopay(up.UpstreamMerchantNo, up.CallConfig)
	case upstream.UpstreamCodeGoldPays:
		upstreamObj, err = upstream.NewGoldPays(up.UpstreamMerchantNo, up.CallConfig)
	default:
		l.Errorf("未知的上游代码[%v]", up.UpstreamCode)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	if err != nil || upstreamObj == nil {
		l.Errorf("获取上游失败,UpstreamCode=%v, err=%v", up.UpstreamCode, err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	resp, err := upstreamObj.QueryBalance()
	if err != nil {
		l.Errorf("查询上游失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	// todo 现在不知道上游会返回什么数据，所以先以ToToPay为示例。后续要重新调整
	var isDivideHundred bool
	switch resp.Currency {
	case model.CurrencyUSD:
		fallthrough
	case model.CurrencyCNY:
		isDivideHundred = true
	case model.CurrencyKHR:
		isDivideHundred = false
	case model.CurrencyINR:
		isDivideHundred = true
	default:
		isDivideHundred = true
	}

	return &types.GetUpstreamBalanceResponse{
		List: []types.UpstreamBalance{
			{
				Balance:           resp.Balance,
				PayoutBalance:     resp.PayOutBalance,
				PayAmountLimit:    resp.PayAmountLimit,
				PayoutAmountLimit: resp.PayoutAmountLimit,
				IsDivideHundred:   isDivideHundred,
				Currency:          resp.Currency,
			},
		},
	}, nil
}
