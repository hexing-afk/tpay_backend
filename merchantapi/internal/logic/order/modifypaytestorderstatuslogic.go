package order

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ModifyPayTestOrderStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModifyPayTestOrderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) ModifyPayTestOrderStatusLogic {
	return ModifyPayTestOrderStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModifyPayTestOrderStatusLogic) ModifyPayTestOrderStatus(merchantId int64, req types.ModifyPayTestOrderStatusRequest) error {
	order, err := model.NewPayOrderModel(l.svcCtx.DbEngine).FindOrderByMerchantId(merchantId, req.OrderNo)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("代收订单[%v]不存在", req.OrderNo)
			return common.NewCodeError(common.OrderNotExist)
		} else {
			l.Errorf("查询代收订单[%v]失败, err=%v", req.OrderNo, err)
			return common.NewCodeError(common.SystemInternalErr)
		}
	}

	if order.Mode != model.PayModeTest {
		l.Errorf("代收订单[%v]不是测试订单，不可手动修改订单状态, mode=%v", req.OrderNo, order.Mode)
		return common.NewCodeError(common.OrderNotOp)
	}

	if order.OrderStatus != model.PayOrderStatusPending {
		l.Errorf("代收订单[%v]不是待支付状态，不可修改订单状态, status=%v", req.OrderNo, order.OrderStatus)
		return common.NewCodeError(common.OrderNotOp)
	}

	switch req.OrderStatus {
	case model.PayOrderStatusFail:
		err = model.NewPayOrderModel(l.svcCtx.DbEngine).UpdateOrderStatus(order.Id, model.PayOrderStatusFail, 0)
	case model.PayOrderStatusPaid:
		err = model.NewPayOrderModel(l.svcCtx.DbEngine).UpdateOrderStatus(order.Id, model.PayOrderStatusPaid, order.ReqAmount)
	default:
		l.Errorf("未知的订单状态，req.OrderStatus=%v", req.OrderStatus)
		return common.NewCodeError(common.InvalidParam)
	}
	if err != nil {
		l.Errorf("修改代收订单[%v]支付状态失败, err=%v", order.Id, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	return nil
}
