package logic

import (
	"context"
	"github.com/tal-tech/go-zero/core/logx"
	"gorm.io/gorm"
	"time"
	"tpay_backend/model"
	"tpay_backend/payapi/internal/common"
	"tpay_backend/payapi/internal/order_notify"
	"tpay_backend/payapi/internal/svc"
	"tpay_backend/upstream"
	"tpay_backend/utils"
)

type TransferPlaceOrder struct {
	logx.Logger
	ctx      context.Context
	svcCtx   *svc.ServiceContext
	merchant *model.Merchant
}

type TransferPlaceOrderRequest struct {
	Mode               string // 模式
	MchOrderNo         string // 商户订单号
	Amount             int64  // 订单金额
	Currency           string // 币种
	OrderSource        int64  // 订单来源
	TradeType          string // 交易类型
	BankName           string // 银行名称
	BankCardNo         string // 卡号
	BankCardHolderName string // 持卡人姓名
	BankBranchName     string // 支行名称
	BankCode           string // 银行代码
	NotifyUrl          string // 异步通知地址
	ReturnUrl          string // 同步跳转地址
	Attach             string // 原样返回字段
	Remark             string // 付款备注
	BatchNo            string // 批量付款批次号
	BatchRowNo         string // 批量付款批次行号
}

type TransferPlaceOrderResponse struct {
	MchOrderNo  string
	OrderNo     string
	OrderStatus int64
}

func NewTransferPlaceOrder(ctx context.Context, svcCtx *svc.ServiceContext, merchant *model.Merchant) *TransferPlaceOrder {
	return &TransferPlaceOrder{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		svcCtx:   svcCtx,
		merchant: merchant,
	}
}

func (l *TransferPlaceOrder) TransferPlaceOrder(req TransferPlaceOrderRequest) (*TransferPlaceOrderResponse, error) {
	// 1.按权重选择上游通道相关信息
	upChannelData, err := NewPickUpstreamChannel(l.ctx, l.svcCtx, l.merchant).PickTransferUpstreamChannelByWeigh(req.TradeType)
	if err != nil {
		l.Errorf("获取通道信息失败, err=%v", err)
		return nil, common.NewCodeError(common.ChannelUnusable)
	}

	l.Infof("代付-选中的通道:%+v", upChannelData)

	// 2.选择上游
	funcLogic := NewFuncLogic(l.svcCtx)
	upstreamObj, err := funcLogic.GetUpstream(upChannelData.UpstreamChannelId)
	if err != nil {
		l.Errorf("代付-选择上游失败:err:%v,req:%+v,merchant:%+v", err, req, l.merchant)
		return nil, common.NewCodeError(common.ChannelUnusable)
	}
	// 检测通道配置
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

	l.Infof("代付-选择上游成功:%+v", upstreamObj)

	// 3.下单 手续费默认为外扣
	order := new(model.TransferOrder)
	order.OrderNo = order.GenerateOrderNo(req.Mode)
	order.MerchantNo = l.merchant.MerchantNo
	order.MerchantOrderNo = req.MchOrderNo
	order.ReqAmount = req.Amount
	order.Currency = req.Currency
	order.NotifyUrl = req.NotifyUrl
	order.ReturnUrl = req.ReturnUrl
	order.OrderStatus = model.TransferOrderStatusPending
	order.NotifyStatus = model.TransferNotifyStatusNot
	order.OrderSource = req.OrderSource
	order.PlatformChannelId = upChannelData.PlatformChannelId
	order.UpstreamChannelId = upChannelData.UpstreamChannelId
	order.BankName = req.BankName
	order.AccountName = req.BankCardHolderName
	order.CardNumber = req.BankCardNo
	order.BranchName = req.BankBranchName
	order.Remark = req.Remark
	order.FeeDeductType = upChannelData.UpChannelDeductionMethod
	order.BankCode = req.BankCode
	order.AreaId = l.merchant.AreaId
	order.MerchantRate = upChannelData.Rate
	order.MerchantSingleFee = upChannelData.SingleFee
	order.Mode = req.Mode
	order.BatchNo = req.BatchNo
	order.BatchRowNo = req.BatchRowNo

	// 提现派单过来的订单，再提现那边就已经算好手续费和商户扣减余额，这边就不计算了
	if order.OrderSource == model.TransferOrderSourceWithdrawAllot {
		// 计算手续费
		order.MerchantFee = 0
		// 收款方实际增加金额 = 请求金额
		order.PayeeRealAmount = order.ReqAmount
		// 账户扣减金额
		order.DecreaseAmount = 0
		// 计算上游金额
		order.UpstreamAmount = order.ReqAmount
	} else {
		// 计算手续费
		order.MerchantFee = utils.CalculatePayOrderFeeMerchant(order.ReqAmount, upChannelData.SingleFee, upChannelData.Rate)

		switch order.FeeDeductType {
		case model.UpstreamChannelDeductionInner: // 内扣
			// 账户扣减金额
			order.DecreaseAmount = order.ReqAmount

			// 收款方实际增加金额 = 请求金额 - 手续费
			order.PayeeRealAmount = order.ReqAmount - order.MerchantFee

			// 计算上游手续费
			order.UpstreamAmount = utils.CalculateUpstreamInnerAmount(order.PayeeRealAmount, upChannelData.UpChannelSingleFee, upChannelData.UpChannelRate)
		case model.UpstreamChannelDeductionOut: // 外扣
			// 收款方实际增加金额 = 请求金额
			order.PayeeRealAmount = order.ReqAmount

			// 账户扣减金额
			order.DecreaseAmount = order.ReqAmount + order.MerchantFee

			// 计算上游手续费
			order.UpstreamAmount = order.ReqAmount
		default:
			l.Errorf("手续费扣款方式有误,FeeDeductType:%v", order.FeeDeductType)
			return nil, common.NewCodeError(common.SystemInternalErr)
		}

		// 检查商户余额是否足够支撑代付扣款
		if l.merchant.Balance < order.DecreaseAmount {
			l.Errorf("商户余额不足,merchant:%+v,DecreaseAmount:%v", l.merchant, order.DecreaseAmount)
			return nil, common.NewCodeError(common.InsufficientBalance)
		}
	}

	order.UpstreamFee = utils.CalculatePayOrderFeeUpstream(order.UpstreamAmount, upChannelData.UpChannelSingleFee, upChannelData.UpChannelRate)

	if err := l.InsertOrder(order); err != nil {
		l.Errorf("代付-下单失败：err:%v, order:%+v", err, order)
		return nil, common.NewCodeError(common.OrderFailed)
	}

	l.Infof("代付-下单成功:Order:%+v", order)

	// 订单属于测试模式还是生产模式
	switch order.Mode {
	case model.TransferModePro:
		// 组合异步回调地址
		upstreamConfig := upstreamObj.GetUpstreamConfig()
		upNotifyUrl, err := funcLogic.GetUpstreamNotifyUrl(upstreamConfig.TransferNotifyPath)
		if err != nil {
			l.Errorf("获取上游异步回调地址失败,transferNotifyPath:%v,err=%v", upstreamConfig.TransferNotifyPath, err)
			return nil, common.NewCodeError(common.SystemInternalErr)
		}

		// 4.请求上游
		payRequest := &upstream.TransferRequest{
			Amount:             order.UpstreamAmount,
			Currency:           order.Currency,
			OrderNo:            order.OrderNo,
			NotifyUrl:          upNotifyUrl,
			ProductType:        upChannelData.UpChannelCode,
			BankName:           req.BankName,
			BankBranchName:     req.BankBranchName,
			BankCardNo:         req.BankCardNo,
			BankCode:           req.BankCode,
			BankCardHolderName: req.BankCardHolderName,
		}
		resp, err := upstreamObj.Transfer(payRequest)
		if err != nil {
			l.Errorf("代付-请求上游失败, req:%+v, err:%v", payRequest, err)
			// 订单失败
			if err := l.UpdateOrderFail(order, err.Error()); err != nil {
				l.Errorf("修改代付订单状态失败, order.Id:%v, err:%v", order.Id, err)
			}
			return nil, common.NewCodeError(common.OrderFailed)
		}

		l.Infof("代付-请求上游成功:payRes:%+v", resp)

		// 5.处理结果，走到这里上下游都已下单成功，因此无论是否处理成功都应该返回订单信息
		data := model.TransferOrder{
			UpstreamOrderNo: resp.UpstreamOrderNo,
		}
		updateErr := model.NewTransferOrderModel(l.svcCtx.DbEngine).UpdateUpstreamOrderInfo(order.Id, data)
		if updateErr != nil {
			l.Errorf("代付-处理结果失败:err:%v, orderId:%v, UpstreamOrderNo:%v", updateErr, order.Id, resp.UpstreamOrderNo)
		}
	case model.TransferModeTest:
		l.Infof("测试模式的代付订单[%v]-下单成功", order.Mode)
	}

	return &TransferPlaceOrderResponse{
		MchOrderNo:  req.MchOrderNo,
		OrderNo:     order.OrderNo,
		OrderStatus: order.OrderStatus,
	}, nil
}

// 插入订单数据
func (l *TransferPlaceOrder) InsertOrder(order *model.TransferOrder) error {
	txErr := l.svcCtx.DbEngine.Transaction(func(tx *gorm.DB) error {
		// 1.插入订单
		if err := model.NewTransferOrderModel(tx).Insert(order); err != nil {
			return err
		}

		if order.Mode == model.TransferModeTest {
			return nil
		}

		// 2.根据订单来源进行操作
		if order.OrderSource == model.TransferOrderSourceWithdrawAllot {
			// 查询提现订单
			withdrawOrder, err := model.NewMerchantWithdrawOrderModel(tx).FindByOrderNo(order.MerchantOrderNo)
			if err != nil {
				return err
			}

			// 绑定提现订单
			if err := model.NewMerchantWithdrawOrderModel(tx).UpdateTransferNo(withdrawOrder.Id, order.OrderNo); err != nil {
				return err
			}
		} else {
			// 扣减商户余额，增加冻结金额
			log := model.WalletLogExt{
				BusinessNo: order.OrderNo,
				Source:     model.AmountSourceTransfer,
				Remark:     "",
			}

			if err := model.NewMerchantModel(tx).MinusBalanceFreezeTx(l.merchant.Id, order.DecreaseAmount, log); err != nil {
				return err
			}
		}

		return nil
	})

	return txErr
}

// 修改-订单失败
func (l *TransferPlaceOrder) UpdateOrderFail(order *model.TransferOrder, failReason string) error {
	txErr := l.svcCtx.DbEngine.Transaction(func(tx *gorm.DB) error {
		// 1.修改订单状态
		if err := model.NewTransferOrderModel(l.svcCtx.DbEngine).UpdateOrderFailById(order.Id, failReason); err != nil {
			return err
		}

		if order.Mode == model.TransferModeTest {
			return nil
		}

		// 2.根据订单来源进行操作
		if order.OrderSource == model.TransferOrderSourceWithdrawAllot {
			// 查询提现订单
			withdrawOrder, err := model.NewMerchantWithdrawOrderModel(tx).FindByOrderNo(order.MerchantOrderNo)
			if err != nil {
				return err
			}

			// 修改提现订单状态
			if err := model.NewMerchantWithdrawOrderModel(tx).UpdateStatusToAllotFail(withdrawOrder.Id); err != nil {
				return err
			}

			// 加提现商户余额，减冻结金额
			log2 := model.WalletLogExt{
				BusinessNo: order.OrderNo,
				Source:     model.AmountSourceWithdraw,
				Remark:     "Withdrawal Failed",
			}
			if err := model.NewMerchantModel(tx).PlusBalanceUnfreezeTx(withdrawOrder.MerchantId, withdrawOrder.DecreaseAmount, log2); err != nil {
				return err
			}
		} else {
			// 2.加商户余额，减冻结金额
			log1 := model.WalletLogExt{
				BusinessNo: order.OrderNo,
				Source:     model.AmountSourceTransfer,
				Remark:     "Payment failed",
			}
			if err := model.NewMerchantModel(tx).PlusBalanceUnfreezeTx(l.merchant.Id, order.DecreaseAmount, log1); err != nil {
				return err
			}
		}

		return nil
	})

	// 3.向下游发送代付订单异步通知
	go order_notify.NewTransferOrderNotify(context.TODO(), l.svcCtx).OrderNotify(order.OrderNo)

	return txErr
}

// 修改-订单成功
func (l *TransferPlaceOrder) UpdateOrderPaid(order *model.TransferOrder) error {
	txErr := l.svcCtx.DbEngine.Transaction(func(tx *gorm.DB) error {
		// 1.修改订单状态
		if err := model.NewTransferOrderModel(l.svcCtx.DbEngine).UpdateOrderPaidById(order.Id); err != nil {
			return err
		}

		if order.Mode == model.TransferModeTest {
			return nil
		}

		// 2.根据订单来源进行操作
		if order.OrderSource == model.TransferOrderSourceWithdrawAllot {
			// 查询提现订单信息
			withdrawOrder, err := model.NewMerchantWithdrawOrderModel(tx).FindByOrderNo(order.MerchantOrderNo)
			if err != nil {
				return err
			}

			// 修改提现订单状态
			if err := model.NewMerchantWithdrawOrderModel(tx).UpdateStatusToAllotSuccess(withdrawOrder.Id); err != nil {
				return err
			}

			// 减少提现商户冻结金额
			log := model.WalletLogExt{
				BusinessNo: withdrawOrder.OrderNo,
				Source:     model.AmountSourceWithdraw,
				Remark:     "Withdraw successfully", // 提现成功
			}
			if err := model.NewMerchantModel(tx).MinusFrozenAmount(withdrawOrder.MerchantId, withdrawOrder.DecreaseAmount, log); err != nil {
				l.Errorf("减商户[%v]冻结金额[%v]失败, err=%v", withdrawOrder.MerchantId, order.DecreaseAmount, err)
				return err
			}

			// 3.4.计算并记录平台收益
			data := &model.PlatformWalletLog{
				BusinessNo:  withdrawOrder.OrderNo,
				Source:      model.PlatformIncomeSourceWithdraw,
				MerchantFee: withdrawOrder.MerchantFee,
				UpstreamFee: order.UpstreamFee,
				Income:      withdrawOrder.MerchantFee - order.UpstreamFee,
				Currency:    withdrawOrder.Currency,
			}
			if err := model.NewPlatformWalletLogModel(tx).Insert(data); err != nil {
				return err
			}

		} else {
			// 减商户冻结金额
			log := model.WalletLogExt{
				BusinessNo: order.OrderNo,
				Source:     model.AmountSourceTransfer,
				Remark:     "",
			}
			if err := model.NewMerchantModel(tx).MinusFrozenAmount(l.merchant.Id, order.DecreaseAmount, log); err != nil {
				return err
			}

			// 记录平台收益
			// 指上下游的手续费差价
			// 例如对接的上游收费0.5%，对接的下游收费1%，那么一笔100元的代收订单平台将获得（1%-0.5%）*100的收益
			data := &model.PlatformWalletLog{
				BusinessNo:  order.OrderNo,
				Source:      model.PlatformIncomeSourceTransfer,
				MerchantFee: order.MerchantFee,
				UpstreamFee: order.UpstreamFee,
				Income:      order.MerchantFee - order.UpstreamFee,
				Currency:    order.Currency,
			}
			if err := model.NewPlatformWalletLogModel(tx).Insert(data); err != nil {
				return err
			}
		}

		return nil
	})

	// 4.向下游发送代付订单异步通知
	go order_notify.NewTransferOrderNotify(context.TODO(), l.svcCtx).OrderNotify(order.OrderNo)

	return txErr
}
