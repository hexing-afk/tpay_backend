package logic

import (
	"context"
	"strings"
	"tpay_backend/model"
	"tpay_backend/payapi/internal/common"

	"tpay_backend/payapi/internal/svc"
	"tpay_backend/payapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type PayOrderQueryLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	merchant *model.Merchant
}

func NewPayOrderQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext, merchant *model.Merchant) PayOrderQueryLogic {
	return PayOrderQueryLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		merchant: merchant,
	}
}

// 验证各项参数
func (l *PayOrderQueryLogic) VerifyParam(req types.PayOrderQueryReq) error {
	if strings.TrimSpace(req.MchOrderNo) == "" && strings.TrimSpace(req.OrderNo) == "" {
		return common.NewCodeErrorWithMsg(common.VerifyParamFailed, "mch_order_no和order_no不能同时为空")
	}

	return nil
}

func (l *PayOrderQueryLogic) PayOrderQuery(req types.PayOrderQueryReq) (*types.PayOrderQueryReply, error) {
	// 1.验证各项参数
	if err := l.VerifyParam(req); err != nil {
		return nil, err
	}

	var payOrder *model.PayOrder
	var err error

	if req.OrderNo != "" {
		// 通过平台单号查询
		payOrder, err = model.NewPayOrderModel(l.svcCtx.DbEngine).FindOneByOrderNo(req.OrderNo)
	} else {
		// 通过"商户号"和"商户单号"查询
		payOrder, err = model.NewPayOrderModel(l.svcCtx.DbEngine).FindMerchantOrder(l.merchant.MerchantNo, req.MchOrderNo)
	}

	if err != nil {
		if err == model.ErrRecordNotFound {
			return nil, common.NewCodeError(common.OrderNotExist)
		} else {
			return nil, common.NewCodeError(common.SystemInternalErr)
		}
	}

	// 后续可以在订单状态是待支付时向上游主动发起一次查询
	return &types.PayOrderQueryReply{
		MchOrderNo:    payOrder.MerchantOrderNo,
		OrderNo:       payOrder.OrderNo,
		Amount:        payOrder.ReqAmount,
		ReqAmount:     payOrder.ReqAmount,
		PaymentAmount: payOrder.PaymentAmount,
		Status:        payOrder.OrderStatus,
		Currency:      payOrder.Currency,
	}, nil
}
