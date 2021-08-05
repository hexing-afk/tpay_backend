package upstream_notify

import (
	"context"
	"errors"
	"github.com/tal-tech/go-zero/core/logx"
	"tpay_backend/model"
	"tpay_backend/payapi/internal/logic"
	"tpay_backend/payapi/internal/svc"
	"tpay_backend/utils"
)

type SyncOrder struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSyncOrder(ctx context.Context, svcCtx *svc.ServiceContext) *SyncOrder {
	return &SyncOrder{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SyncOrder) SyncPayOrder(order *model.PayOrder, orderStatus int64, failReason string) error {
	// 1.查询商户
	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneByMerchantNo(order.MerchantNo)
	if err != nil {
		l.Errorf("查询订单商户信息失败, MerchantNo:%v, err:%v", order.MerchantNo, err)
		return errors.New("系统内部错误")
	}

	l.Infof("商户信息：%+v", merchant)

	var upErr error
	payLogic := logic.NewPayLogic(context.TODO(), l.svcCtx, merchant)
	switch orderStatus {
	case model.PayOrderStatusPaid:
		// 1.判断实际支付金额是否与请求上游金额一致
		if order.UpstreamAmount != order.PaymentAmount {
			l.Errorf("代收订单[%v]上游请求金额与实际支付金额不一致, reqAmount:%v, payAmount:%v", order.OrderNo, order.UpstreamAmount, order.PaymentAmount)

			// 计算手续费
			order.MerchantFee = utils.CalculatePayOrderFeeMerchant(order.PaymentAmount, order.MerchantSingleFee, order.MerchantRate)

			// 账户增加的金额=订单请求金额-商户手续费
			order.IncreaseAmount = order.PaymentAmount - order.MerchantFee
		}
		upErr = payLogic.UpdateOrderPaid(order)
	case model.PayOrderStatusFail:
		upErr = payLogic.UpdateOrderFail(order, failReason)
	default:
		l.Errorf("未知的订单状态, order_status:%v", orderStatus)
		return errors.New("订单状态不对")
	}

	if upErr != nil {
		l.Errorf("修改订单状态和商户余额失败, orderNo:%v, MerchantNo:%v, err:%v", order.OrderNo, order.MerchantNo, err)
		return errors.New("系统内部错误")
	}

	return nil
}

func (l *SyncOrder) SyncTransferOrder(order *model.TransferOrder, orderStatus int64, failReason string) error {
	// 1.查询商户
	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneByMerchantNo(order.MerchantNo)
	if err != nil {
		l.Errorf("查询订单商户信息失败, MerchantNo:%v, err:%v", order.MerchantNo, err)
		return errors.New("系统内部错误")
	}

	l.Infof("商户信息：%+v", merchant)

	var upErr error
	transferOrder := logic.NewTransferPlaceOrder(context.TODO(), l.svcCtx, merchant)
	switch orderStatus {
	case model.TransferOrderStatusPaid:
		upErr = transferOrder.UpdateOrderPaid(order)
	case model.TransferOrderStatusFail:
		upErr = transferOrder.UpdateOrderFail(order, failReason)
	default:
		l.Errorf("未知的订单状态, order_status:%v", orderStatus)
		return errors.New("订单状态不对")
	}

	if upErr != nil {
		l.Errorf("修改订单状态和商户余额失败, orderNo:%v, MerchantNo:%v, err:%v", order.OrderNo, order.MerchantNo, err)
		return errors.New("系统内部错误")
	}

	return nil
}
