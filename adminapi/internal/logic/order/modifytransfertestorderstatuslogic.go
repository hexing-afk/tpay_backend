package order

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ModifyTransferTestOrderStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModifyTransferTestOrderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) ModifyTransferTestOrderStatusLogic {
	return ModifyTransferTestOrderStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModifyTransferTestOrderStatusLogic) ModifyTransferTestOrderStatus(req types.ModifyTransferTestOrderStatusRequest) error {
	order, err := model.NewTransferOrderModel(l.svcCtx.DbEngine).FindByOrderNo(req.OrderNo)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("代付订单[%v]不存在", req.OrderNo)
			return common.NewCodeError(common.OrderNotExist)
		} else {
			l.Errorf("查询代付订单[%v]失败, err=%v", req.OrderNo, err)
			return common.NewCodeError(common.SystemInternalErr)
		}
	}

	if order.Mode != model.TransferModeTest {
		l.Errorf("代付订单[%v]不是测试订单，不可手动修改订单状态, mode=%v", req.OrderNo, order.Mode)
		return common.NewCodeError(common.OrderNotOp)
	}

	if order.OrderStatus != model.TransferOrderStatusPending {
		l.Errorf("代付订单[%v]不是待支付状态，不可修改订单状态, status=%v", req.OrderNo, order.OrderStatus)
		return common.NewCodeError(common.OrderNotOp)
	}

	switch req.OrderStatus {
	case model.TransferOrderStatusPaid:
		err = model.NewTransferOrderModel(l.svcCtx.DbEngine).UpdateOrderStatus(order.Id, model.TransferOrderStatusPaid)
	case model.TransferOrderStatusFail:
		err = model.NewTransferOrderModel(l.svcCtx.DbEngine).UpdateOrderStatus(order.Id, model.TransferOrderStatusFail)
	default:
		l.Errorf("未知的订单状态，req.OrderStatus=%v", req.OrderStatus)
		return common.NewCodeError(common.InvalidParam)
	}
	if err != nil {
		l.Errorf("修改代付订单[%v]支付状态失败, err=%v", order.Id, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	return nil
}
