package logic

import (
	"context"
	"tpay_backend/model"
	"tpay_backend/payapi/internal/common"
	"tpay_backend/payapi/internal/svc"
	"tpay_backend/payapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type TransferOrderQueryLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	merchant *model.Merchant
}

func NewTransferOrderQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext, merchant *model.Merchant) TransferOrderQueryLogic {
	return TransferOrderQueryLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		merchant: merchant,
	}
}

func (l *TransferOrderQueryLogic) VerifyParam(req types.TransferOrderQueryReq) error {
	if req.OrderNo == "" && req.MchOrderNo == "" {
		l.Errorf("order_no和mch_order_no同时为空")
		return common.NewCodeError(common.VerifyParamFailed)
	}

	return nil
}

func (l *TransferOrderQueryLogic) TransferOrderQuery(req types.TransferOrderQueryReq) (*types.TransferOrderQueryReply, error) {
	// 1.验证参数
	if err := l.VerifyParam(req); err != nil {
		return nil, err
	}

	// 2.内部查询订单
	var err error
	var order *model.TransferOrder
	if req.OrderNo != "" {
		// 通过平台订单号查询
		order, err = model.NewTransferOrderModel(l.svcCtx.DbEngine).FindByOrderNo(req.OrderNo)
	} else {
		// 通过商户订单号查询
		order, err = model.NewTransferOrderModel(l.svcCtx.DbEngine).FindByMchOrderNo(l.merchant.MerchantNo, req.MchOrderNo)
	}

	if err != nil {
		if err == model.ErrRecordNotFound {
			return nil, common.NewCodeError(common.OrderNotExist)
		} else {
			l.Errorf("查询代付订单失败, err=%v", err)
			return nil, common.NewCodeError(common.SystemInternalErr)
		}
	}

	// 后续可以在订单状态是待支付时向上游主动发起一次查询

	return &types.TransferOrderQueryReply{
		MchOrderNo: order.MerchantOrderNo,
		OrderNo:    order.OrderNo,
		Amount:     order.ReqAmount,
		ReqAmount:  order.ReqAmount,
		Status:     order.OrderStatus,
		Currency:   order.Currency,
	}, nil
}
