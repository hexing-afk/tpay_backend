package logic

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"
	"tpay_backend/utils"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type WalletLogExportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWalletLogExportLogic(ctx context.Context, svcCtx *svc.ServiceContext) WalletLogExportLogic {
	return WalletLogExportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WalletLogExportLogic) WalletLogExport(req types.WalletLogExportRequest) (*types.WalletLogExportResponse, error) {
	day := utils.TimeSubToDay(req.StartCreateTime, req.EndCreateTime)
	if day < 0 || day > 31 {
		l.Errorf("导出时间超出限制(31天), start:%v, end:%v", req.StartCreateTime, req.EndCreateTime)
		return nil, common.NewCodeError(common.ExportLimit31Day)
	}

	currency, err := model.NewCurrencyModel(l.svcCtx.DbEngine).FindByCurrency(req.Currency)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("币种不存在, err=%v", err)
			return nil, common.NewCodeError(common.InvalidParam)
		} else {
			l.Errorf("查询币种失败, err=%v", err)
			return nil, common.NewCodeError(common.SystemInternalErr)
		}
	}

	f := model.FindWalletExportData{
		Username:        req.MerchantAccount,
		OpType:          req.OpType,
		Source:          req.OrderType,
		Currency:        req.Currency,
		StartCreateTime: req.StartCreateTime,
		EndCreateTime:   req.EndCreateTime,
		IdOrBusinessNo:  req.IdOrBusinessNo,
		OpTypeList:      []int64{model.OpTypeAddBalance, model.OpTypeMinusBalance},
	}
	data, err := model.NewMerchantWalletLogModel(l.svcCtx.DbEngine).FindExportData(f)
	if err != nil {
		logx.Errorf("查询账户流水日志失败, err=[%v]", err)
		return nil, common.NewCodeError(common.ExportFail)
	}

	if data == nil || data.Total == 0 {
		logx.Errorf("账户流水没数据")
		return nil, common.NewCodeError(common.NotData)
	}

	var list []types.WalletLogExportData
	for _, v := range data.List {
		var opType, orderType string
		switch v.OpType {
		case model.OpTypeAddBalance:
			opType = "入账"
		case model.OpTypeMinusBalance:
			opType = "出账"
		case model.OpTypeAddFrozen:
			opType = "增加冻结金额"
		case model.OpTypeMinusFrozen:
			opType = "减少解冻金额"
		}

		switch v.Source {
		case model.AmountSourcePlatform:
			orderType = "手动调账"
		case model.AmountSourceRecharge:
			orderType = "充值"
		case model.AmountSourceWithdraw:
			orderType = "提现"
		case model.AmountSourceCollection:
			orderType = "收款"
		case model.AmountSourceTransfer:
			orderType = "代付"
		}

		list = append(list, types.WalletLogExportData{
			Id:           v.Id,
			CreateTime:   v.CreateTime,
			OpType:       opType,
			ChangeAmount: v.ChangeAmount,
			AfterBalance: v.AfterBalance,
			BusinessNo:   v.BusinessNo,
			OrderType:    orderType,
			Remark:       v.Remark,
		})
	}

	return &types.WalletLogExportResponse{
		Total:           data.Total,
		IsDivideHundred: currency.IsDivideHundred == model.DivideHundred,
		List:            list,
	}, nil
}
