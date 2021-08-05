package order

import (
	"context"
	"tpay_backend/model"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetWalletLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	userId int64
}

func NewGetWalletLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext, userId int64) GetWalletLogListLogic {
	return GetWalletLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		userId: userId,
	}
}

func (l *GetWalletLogListLogic) GetWalletLogList(req types.GetWalletLogListReq) (*types.GetWalletLogListReply, error) {

	f := model.FindMerchantWalletLogList{
		Page:            req.Page,
		PageSize:        req.PageSize,
		OpType:          req.OpType,
		Source:          req.OrderType,
		StartCreateTime: req.StartCreateTime,
		EndCreateTime:   req.EndCreateTime,
		IdOrBusinessNo:  req.IdOrBusinessNo,
		MerchantId:      l.userId,
		OpTypeList:      []int64{model.OpTypeAddBalance, model.OpTypeMinusBalance},
	}
	walletLogs, total, err := model.NewMerchantWalletLogModel(l.svcCtx.DbEngine).FindList(f)
	if err != nil {
		logx.Errorf("查询账户流水日志失败, err=[%v]", err)
	}

	var list []types.WalletLogList
	for _, v := range walletLogs {
		list = append(list, types.WalletLogList{
			Id:           v.Id,
			CreateTime:   v.CreateTime,
			OpType:       v.OpType,
			ChangeAmount: v.ChangeAmount,
			AfterBalance: v.AfterBalance,
			BusinessNo:   v.BusinessNo,
			OrderType:    v.Source,
			Remark:       v.Remark,
		})

	}

	return &types.GetWalletLogListReply{
		List:  list,
		Total: total,
	}, nil
}
