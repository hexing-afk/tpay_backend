package logic

import (
	"context"
	"fmt"
	"strings"
	"tpay_backend/model"
	"tpay_backend/payapi/internal/common"
	"tpay_backend/utils"

	"tpay_backend/payapi/internal/svc"
	"tpay_backend/payapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type PayTestLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	merchant *model.Merchant
}

func NewPayTestLogic(ctx context.Context, svcCtx *svc.ServiceContext, merchant *model.Merchant) PayTestLogic {
	return PayTestLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		merchant: merchant,
	}
}

// 验证各项参数
func (l *PayTestLogic) VerifyParam(req types.PayTestReq) error {
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
	payTradeTypeSlice, err := model.NewGlobalConfigModel(l.svcCtx.DbEngine).FindValueByKey(model.ConfigPayTradeTypeSlice)
	if err != nil {
		logx.Errorf("查询代收交易类型出错,key:%v, err:%v", model.ConfigPayTradeTypeSlice, err)
		return common.NewCodeError(common.SystemInternalErr)
	}

	if !utils.InSlice(strings.TrimSpace(req.TradeType), strings.Split(payTradeTypeSlice, ",")) {
		return common.NewCodeErrorWithMsg(common.VerifyParamFailed, "trade_type不被支持")
	}

	if strings.TrimSpace(req.NotifyUrl) == "" {
		return common.NewCodeErrorWithMsg(common.VerifyParamFailed, "notify_url不能为空")
	}

	// 订单的币种和商户账号支持的币种不一致
	if req.Currency != l.merchant.Currency {
		return common.NewCodeErrorWithMsg(common.VerifyParamFailed, fmt.Sprintf("currency错误,该商户只支持(%s)货币类型", l.merchant.Currency))
	}

	//// 检查商户的余额是否足够
	//if req.Amount > l.merchant.Balance {
	//	return common.NewCodeErrorWithMsg(common.VerifyParamFailed, "商户余额不足")
	//}

	// 查询商户订单是否已经存在
	exist, err := model.NewPayOrderModel(l.svcCtx.DbEngine).MerchantOrderNoExist(req.MerchantNo, req.MchOrderNo)
	if err != nil {
		logx.Errorf("查询商户订单是否已经存在出错,MerchantNo:%v, MchOrderNo:%v, err:%v", err, req.MerchantNo, req.MchOrderNo, err)
		return common.NewCodeError(common.SystemInternalErr)
	}

	if exist {
		return common.NewCodeErrorWithMsg(common.VerifyParamFailed, "mch_order_no重复")
	}

	return nil
}

func (l *PayTestLogic) PayTest(req types.PayTestReq) (*types.PayTestReply, error) {
	// 1.验证各项参数
	if err := l.VerifyParam(req); err != nil {
		return nil, err
	}

	// 2.按权重选择上游通道相关信息
	upChannelData, err := NewPickUpstreamChannel(l.ctx, l.svcCtx, l.merchant).PickPayUpstreamChannelByWeigh(req.TradeType)
	if err != nil {
		return nil, err
	}

	logx.Infof("代收-选中的通道:%+v", upChannelData)

	// 3.选择上游
	funcLogic := NewFuncLogic(l.svcCtx)
	upstreamObj, err := funcLogic.GetUpstream(upChannelData.UpstreamChannelId)
	if err != nil {
		logx.Errorf("代收-选择上游失败:err:%v,req:%+v,merchant:%+v", err, req, l.merchant)
		return nil, common.NewCodeError(common.ChannelUnusable)
	}

	logx.Infof("代收-选择上游成功:%+v", upstreamObj)

	// 4.下单
	payOrder := new(model.PayOrder)
	payOrder.OrderNo = "TEST" + payOrder.GenerateOrderNo(model.PayModeTest)
	payOrder.MerchantOrderNo = req.MchOrderNo
	payOrder.MerchantNo = req.MerchantNo
	payOrder.Currency = req.Currency
	payOrder.ReqAmount = req.Amount
	payOrder.OrderStatus = model.PayOrderStatusPending
	payOrder.NotifyStatus = model.PayNotifyStatusNot
	payOrder.NotifyUrl = req.NotifyUrl
	payOrder.ReturnUrl = req.ReturnUrl
	payOrder.Subject = req.Subject
	payOrder.PlatformChannelId = upChannelData.PlatformChannelId
	payOrder.UpstreamChannelId = upChannelData.UpstreamChannelId
	payOrder.AreaId = l.merchant.AreaId
	payOrder.Mode = model.PayModeTest
	payOrder.MerchantRate = upChannelData.Rate
	payOrder.MerchantSingleFee = upChannelData.SingleFee

	// 计算手续费
	payOrder.MerchantFee = utils.CalculatePayOrderFeeMerchant(req.Amount, upChannelData.SingleFee, upChannelData.Rate)

	// 账户增加的金额=订单请求金额-商户手续费
	payOrder.IncreaseAmount = payOrder.ReqAmount - payOrder.MerchantFee

	// 计算上游手续费
	payOrder.UpstreamAmount = req.Amount
	payOrder.UpstreamFee = utils.CalculatePayOrderFeeUpstream(payOrder.UpstreamAmount, upChannelData.UpChannelSingleFee, upChannelData.UpChannelRate)

	// 插入订单
	if err := model.NewPayOrderModel(l.svcCtx.DbEngine).Insert(payOrder); err != nil {
		l.Errorf("插入订单数据失败, order:%+v, err=%v", payOrder, err)
		return nil, common.NewCodeError(common.OrderFailed)
	}

	logx.Infof("代收-下单成功:payOrder:%+v", payOrder)

	// 6.返回给下游
	return &types.PayTestReply{
		MchOrderNo: req.MchOrderNo,
		OrderNo:    payOrder.OrderNo,
		PayUrl:     "",
	}, nil
}
