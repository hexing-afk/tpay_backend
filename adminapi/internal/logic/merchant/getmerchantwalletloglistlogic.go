package merchant

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetMerchantWalletLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMerchantWalletLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetMerchantWalletLogListLogic {
	return GetMerchantWalletLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMerchantWalletLogListLogic) GetMerchantWalletLogList(req types.GetMerchantWalletLogListRequest) (*types.GetMerchantWalletLogListResponse, error) {
	f := model.FindMerchantWalletLogList{
		Page:            req.Page,
		PageSize:        req.PageSize,
		Id:              req.SerialNo,
		BusinessNo:      req.BusinessNo,
		Username:        req.MerchantAccount,
		OpType:          req.OpType,
		Source:          req.OrderType,
		Currency:        req.Currency,
		StartCreateTime: req.StartCreateTime,
		EndCreateTime:   req.EndCreateTime,
		IdOrBusinessNo:  req.IdOrBusinessNo,
		OpTypeList:      []int64{model.OpTypeAddBalance, model.OpTypeMinusBalance},
	}
	data, total, err := model.NewMerchantWalletLogModel(l.svcCtx.DbEngine).FindList(f)
	if err != nil {
		l.Errorf("查询商户钱包日志列表失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	if total == 0 {
		return nil, nil
	}

	var list []types.MerchantWalletLogList
	for _, v := range data {
		list = append(list, types.MerchantWalletLogList{
			SerialNo:           v.Id,
			BusinessNo:         v.BusinessNo,
			MerchantAccount:    v.Username,
			Currency:           v.Currency,
			ChangeAmount:       v.ChangeAmount,
			ChangeAfterBalance: v.AfterBalance,
			OpType:             v.OpType,
			OrderType:          v.Source,
			Remark:             v.Remark,
			CreateTime:         v.CreateTime,
		})
	}

	return &types.GetMerchantWalletLogListResponse{
		Total: total,
		List:  list,
	}, nil
}
