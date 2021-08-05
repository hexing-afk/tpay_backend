package order

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type TransferInterfaceOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	userId int64
}

func NewTransferInterfaceOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext, userId int64) TransferInterfaceOrderListLogic {
	return TransferInterfaceOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		userId: userId,
	}
}

func (l *TransferInterfaceOrderListLogic) TransferInterfaceOrderList(req types.TransferInterfaceOrderListRequest) (*types.TransferInterfaceOrderListReply, error) {
	//根据商户id查询商户编号
	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneById(l.userId)
	if err != nil {
		l.Errorf("查询商户出错, userId[%v], err[%v]", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	// 普通商户没有提现派单订单，只有平台商户有
	f := model.FindTransferOrderList{
		MerchantNo:      merchant.MerchantNo,
		Page:            req.Page,
		PageSize:        req.PageSize,
		StartCreateTime: req.StartCreateTime,
		EndCreateTime:   req.EndCreateTime,
		OrderNo:         req.OrderNo,
		MerchantOrderNo: req.MerchantOrderNo,
		OrderStatus:     req.OrderStatus,
		OrderSourceList: []int64{model.TransferOrderSourceInterface},
		OrderType:       req.OrderType,
	}
	dataList, total, err := model.NewTransferOrderModel(l.svcCtx.DbEngine).FindList(f)
	if err != nil {
		l.Errorf("查询代付订单失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var list []types.TransferInterfaceOrderList
	for _, data := range dataList {
		list = append(list, types.TransferInterfaceOrderList{
			OrderNo:             data.OrderNo,
			MerchantOrderNo:     data.MerchantOrderNo,
			Currency:            data.Currency,
			OrderAmount:         data.ReqAmount,
			MerchantFee:         data.MerchantFee,
			PayeeRealAmount:     data.PayeeRealAmount,
			PlatformChannelName: data.PlatformChannelName,
			OrderStatus:         data.OrderStatus,
			CreateTime:          data.CreateTime,
			UpdateTime:          data.UpdateTime,
			OrderSource:         data.OrderSource,
			Remark:              data.Remark,
		})
	}

	return &types.TransferInterfaceOrderListReply{
		Total: total,
		List:  list,
	}, nil
}
