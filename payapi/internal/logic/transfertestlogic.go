package logic

import (
	"context"
	"fmt"
	"strings"
	"tpay_backend/model"
	"tpay_backend/payapi/internal/common"
	"tpay_backend/payapi/internal/svc"
	"tpay_backend/payapi/internal/types"
	"tpay_backend/utils"

	"github.com/tal-tech/go-zero/core/logx"
)

type TransferTestLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	merchant *model.Merchant
}

func NewTransferTestLogic(ctx context.Context, svcCtx *svc.ServiceContext, merchant *model.Merchant) TransferTestLogic {
	return TransferTestLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		merchant: merchant,
	}
}

func (l *TransferTestLogic) VerifyParam(req types.TransferTestReq) error {
	if req.Amount < 1 {
		return common.NewCodeErrorWithMsg(common.VerifyParamFailed, "amount不能小于1")
	}

	if strings.TrimSpace(req.Currency) == "" {
		return common.NewCodeErrorWithMsg(common.VerifyParamFailed, "currency不能为空")
	}

	if strings.TrimSpace(req.MchOrderNo) == "" {
		return common.NewCodeErrorWithMsg(common.VerifyParamFailed, "mch_order_no不能为空")
	}

	if strings.TrimSpace(req.TradeType) == "" {
		return common.NewCodeErrorWithMsg(common.VerifyParamFailed, "trade_type不能为空")
	}

	// 检查交易类型是否存在
	transferTradeTypeSlice, err := model.NewGlobalConfigModel(l.svcCtx.DbEngine).FindValueByKey(model.ConfigTransferTradeTypeSlice)
	if err != nil {
		logx.Errorf("查询代收交易类型出错,key:%v, err:%v", model.ConfigPayTradeTypeSlice, err)
		return common.NewCodeError(common.SystemInternalErr)
	}

	if !utils.InSlice(strings.TrimSpace(req.TradeType), strings.Split(transferTradeTypeSlice, ",")) {
		return common.NewCodeErrorWithMsg(common.VerifyParamFailed, "trade_type不被支持")
	}

	//if !utils.InSlice(strings.TrimSpace(req.TradeType), common.TransferTradeTypeSlice) {
	//	return common.NewCodeErrorWithMsg(common.VerifyParamFailed, "trade_type不被支持")
	//}

	if strings.TrimSpace(req.NotifyUrl) == "" {
		return common.NewCodeErrorWithMsg(common.VerifyParamFailed, "notify_url不能为空")
	}

	if strings.TrimSpace(req.BankName) == "" {
		return common.NewCodeErrorWithMsg(common.VerifyParamFailed, "bank_name不能为空")
	}

	if strings.TrimSpace(req.BankCardHolderName) == "" {
		return common.NewCodeErrorWithMsg(common.VerifyParamFailed, "bank_card_holder_name不能为空")
	}

	if strings.TrimSpace(req.BankCardNo) == "" {
		return common.NewCodeErrorWithMsg(common.VerifyParamFailed, "bank_card_no不能为空")
	}

	// 订单的币种和商户账号支持的币种不一致
	if req.Currency != l.merchant.Currency {
		return common.NewCodeErrorWithMsg(common.VerifyParamFailed, fmt.Sprintf("currency错误,该商户只支持(%s)货币类型", l.merchant.Currency))
	}

	exist, err := model.NewTransferOrderModel(l.svcCtx.DbEngine).MerchantOrderNoExist(req.MerchantNo, req.MchOrderNo)
	if err != nil {
		logx.Errorf("查询商户订单是否已经存在出错,MerchantNo:%v, MchOrderNo:%v, err:%v", err, req.MerchantNo, req.MchOrderNo, err)
		return common.NewCodeError(common.SystemInternalErr)
	}

	if exist {
		return common.NewCodeError(common.DuplicateOrderNO)
	}

	if l.merchant.Balance < req.Amount {
		return common.NewCodeErrorWithMsg(common.VerifyParamFailed, "商户余额不足")
	}

	return nil
}

func (l *TransferTestLogic) TransferTest(req types.TransferTestReq) (*types.TransferTestReply, error) {
	// 2.下单
	param := TransferPlaceOrderRequest{
		MchOrderNo:         req.MchOrderNo,
		Amount:             req.Amount,
		Currency:           req.Currency,
		OrderSource:        model.TransferOrderSourceInterface,
		TradeType:          req.TradeType,
		BankName:           req.BankName,
		BankCardNo:         req.BankCardNo,
		BankCardHolderName: req.BankCardHolderName,
		BankBranchName:     req.BankBranchName,
		BankCode:           req.BankCode,
		NotifyUrl:          req.NotifyUrl,
		ReturnUrl:          req.ReturnUrl,
		Attach:             req.Attach,
		Remark:             req.Remark,
		Mode:               model.TransferModeTest, // 测试模式订单
	}
	placeOrder := NewTransferPlaceOrder(l.ctx, l.svcCtx, l.merchant)
	resp, err := placeOrder.TransferPlaceOrder(param)
	if err != nil {
		l.Errorf("下单失败")
		return nil, err
	}

	// 3.返回给下游
	return &types.TransferTestReply{
		MchOrderNo: req.MchOrderNo,
		OrderNo:    resp.OrderNo,
		Status:     resp.OrderStatus,
	}, nil
}
