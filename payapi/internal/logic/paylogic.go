package logic

import (
	"context"
	"fmt"
	"strings"
	"time"
	"tpay_backend/model"
	"tpay_backend/payapi/internal/common"
	"tpay_backend/payapi/internal/order_notify"
	"tpay_backend/payapi/internal/svc"
	"tpay_backend/payapi/internal/types"
	"tpay_backend/upstream"
	"tpay_backend/utils"

	"gorm.io/gorm"

	"github.com/tal-tech/go-zero/core/logx"
)

type PayLogic struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	merchant *model.Merchant
}

func NewPayLogic(ctx context.Context, svcCtx *svc.ServiceContext, merchant *model.Merchant) PayLogic {
	return PayLogic{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		merchant: merchant,
	}
}

// 验证各项参数
func (l *PayLogic) VerifyParam(req types.PayReq) error {
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
		logx.Errorf("代收trade_type不被支持,请求:%v,系统配置:%v", req.TradeType, payTradeTypeSlice)
		return common.NewCodeErrorWithMsg(common.VerifyParamFailed, "trade_type不被支持")
	}

	if strings.TrimSpace(req.NotifyUrl) == "" {
		return common.NewCodeErrorWithMsg(common.VerifyParamFailed, "notify_url不能为空")
	}

	// 订单的币种和商户账号支持的币种不一致
	if req.Currency != l.merchant.Currency {
		return common.NewCodeErrorWithMsg(common.VerifyParamFailed, fmt.Sprintf("currency错误,该商户只支持(%s)货币类型", l.merchant.Currency))
	}

	// 查询商户订单是否已经存在
	exist, err := model.NewPayOrderModel(l.svcCtx.DbEngine).MerchantOrderNoExist(req.MerchantNo, req.MchOrderNo)
	if err != nil {
		logx.Errorf("查询商户订单是否已经存在出错,MerchantNo:%v, MchOrderNo:%v, err:%v", req.MerchantNo, req.MchOrderNo, err)
		return common.NewCodeError(common.SystemInternalErr)
	}

	if exist {
		return common.NewCodeErrorWithMsg(common.VerifyParamFailed, "mch_order_no重复")
	}

	return nil
}

func (l *PayLogic) Pay(req types.PayReq) (*types.PayReply, error) {
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

	planformChannel, err := model.NewPlatformChannelModel(l.svcCtx.DbEngine).FindOneById(upChannelData.PlatformChannelId)
	if err != nil {
		l.Errorf("查询通道配置失败, ChannelId=%v, err=%v", upChannelData.PlatformChannelId, err)
		return nil, err
	}

	nowStr := utils.TodaySec(time.Local)
	if nowStr < planformChannel.StartTime || nowStr > planformChannel.EndTime {
		logx.Errorf("不在时间内%v=>[%v,%v]", nowStr, planformChannel.StartTime, planformChannel.EndTime)
		return nil, common.NewCodeError(common.NotInTime)
	}
	if planformChannel.StartAmount > req.Amount || planformChannel.EndAmount < req.Amount {
		logx.Errorf("不在金额范围内%v=>[%v,%v]", req.Amount, planformChannel.StartAmount, planformChannel.EndAmount)
		return nil, common.NewCodeError(common.NotInAmt)
	}

	// 3.选择上游
	funcLogic := NewFuncLogic(l.svcCtx)
	upstreamObj, err := funcLogic.GetUpstream(upChannelData.UpstreamChannelId)
	if err != nil {
		logx.Errorf("代收-选择上游失败:err:%v,req:%+v,merchant:%+v", err, req, l.merchant)
		return nil, common.NewCodeError(common.ChannelUnusable)
	}

	logx.Infof("代收-选择上游成功:%+v", upstreamObj)

	// 组合异步回调地址
	upstreamConfig := upstreamObj.GetUpstreamConfig()
	logx.Infof("upstreamConfig=%v", upstreamConfig)
	upNotifyUrl, err := funcLogic.GetUpstreamNotifyUrl(upstreamConfig.PayNotifyPath)
	if err != nil {
		l.Errorf("获取上游异步回调地址失败,payNotifyPath:%v,err=%v", upstreamConfig.PayNotifyPath, err)
		return nil, common.NewCodeError(common.SystemInternalErr)
	}

	// 4.下单
	payOrder := new(model.PayOrder)
	payOrder.OrderNo = payOrder.GenerateOrderNo(model.PayModePro)
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
	payOrder.Mode = model.PayModePro
	payOrder.MerchantRate = upChannelData.Rate
	payOrder.MerchantSingleFee = upChannelData.SingleFee

	// 计算手续费
	payOrder.MerchantFee = utils.CalculatePayOrderFeeMerchant(req.Amount, upChannelData.SingleFee, upChannelData.Rate)

	// 账户增加的金额=订单请求金额-商户手续费
	payOrder.IncreaseAmount = payOrder.ReqAmount - payOrder.MerchantFee

	// 计算上游手续费
	payOrder.UpstreamAmount = req.Amount
	payOrder.UpstreamFee = utils.CalculatePayOrderFeeUpstream(payOrder.UpstreamAmount, upChannelData.UpChannelSingleFee, upChannelData.UpChannelRate)

	if err := l.InsertOrder(payOrder); err != nil {
		logx.Errorf("代收-下单失败:err:%v,payOrder:%+v", err, payOrder)
		return nil, err
	}

	logx.Infof("代收-下单成功:payOrder:%+v", payOrder)

	payRequest := &upstream.PayRequest{
		Amount:      payOrder.UpstreamAmount,
		Currency:    payOrder.Currency,
		OrderNo:     payOrder.OrderNo,
		NotifyUrl:   upNotifyUrl,
		ReturnUrl:   payOrder.ReturnUrl,
		ProductType: upChannelData.UpChannelCode,
		Subject:     req.Subject,
	}

	// 5.请求上游
	payRes, err := upstreamObj.Pay(payRequest)
	if err != nil {
		logx.Errorf("代收-请求上游失败:err:%v", err)
		if err := l.UpdateOrderFail(payOrder, err.Error()); err != nil {
			l.Errorf("修改订单状态失败, err=%v", err)
		}
		return nil, common.NewCodeError(common.OrderFailed)
	}

	logx.Infof("代收-请求上游成功:payRes:%+v", payRes)

	// 5.处理结果，走到这里上下游都已下单成功，因此无论是否处理成功都应该返回订单信息
	updateErr := model.NewPayOrderModel(l.svcCtx.DbEngine).UpdateUpstreamInfo(payOrder.Id, payRes.UpstreamOrderNo)
	if updateErr != nil {
		logx.Errorf("代收-处理结果失败:err:%v, payOrderId:%v, UpstreamOrderNo:%v", updateErr, payOrder.Id, payRes.UpstreamOrderNo)
		//return nil, updateErr
	}

	// 6.返回给下游
	return &types.PayReply{
		MchOrderNo: req.MchOrderNo,
		OrderNo:    payOrder.OrderNo,
		PayUrl:     payRes.PayUrl,
	}, nil
}

// 插入订单
func (l *PayLogic) InsertOrder(payOrder *model.PayOrder) error {

	txErr := l.svcCtx.DbEngine.Transaction(func(tx *gorm.DB) error {
		// 1.插入订单
		if err := model.NewPayOrderModel(tx).Insert(payOrder); err != nil {
			return err
		}

		// 2.xxx

		return nil
	})

	return txErr
}

// 修改-订单支付失败
func (l *PayLogic) UpdateOrderFail(payOrder *model.PayOrder, failReason string) error {
	txErr := l.svcCtx.DbEngine.Transaction(func(tx *gorm.DB) error {
		// 1.修改订单状态
		if err := model.NewPayOrderModel(tx).UpdateOrderFail(payOrder.Id); err != nil {
			return err
		}

		// 2.xxx

		return nil
	})

	return txErr
}

// 修改-订单支付成功
func (l *PayLogic) UpdateOrderPaid(payOrder *model.PayOrder) error {
	txErr := l.svcCtx.DbEngine.Transaction(func(tx *gorm.DB) error {
		// 2.修改订单状态和新的手续费
		if err := model.NewPayOrderModel(tx).UpdateOrderPaid(payOrder.Id, payOrder); err != nil {
			return err
		}

		// 3.增加商户余额
		log := model.WalletLogExt{
			BusinessNo: payOrder.OrderNo,
			Source:     model.AmountSourceCollection,
			Remark:     "",
		}
		if err := model.NewMerchantModel(tx).PlusBalance(l.merchant.Id, payOrder.IncreaseAmount, log); err != nil {
			return err
		}

		// 4.计算并记录平台收益
		// 指上下游的手续费差价
		// 例如对接的上游收费0.5%，对接的下游收费1%，那么一笔100元的代收订单平台将获得（1%-0.5%）*100的收益
		data := &model.PlatformWalletLog{
			BusinessNo:  payOrder.OrderNo,
			Source:      model.PlatformIncomeSourcePay,
			MerchantFee: payOrder.MerchantFee,
			UpstreamFee: payOrder.UpstreamFee,
			Income:      payOrder.MerchantFee - payOrder.UpstreamFee,
			Currency:    payOrder.Currency,
		}
		if err := model.NewPlatformWalletLogModel(tx).Insert(data); err != nil {
			return err
		}

		return nil
	})

	// 4.向下游发送代收订单异步通知
	go order_notify.NewPayOrderNotify(context.TODO(), l.svcCtx).OrderNotify(payOrder.OrderNo)

	return txErr
}
