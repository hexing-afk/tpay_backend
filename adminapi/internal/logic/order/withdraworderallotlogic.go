package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"gorm.io/gorm"
	"strings"
	"time"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"
	"tpay_backend/model"
	"tpay_backend/utils"
)

type WithdrawOrderAllotLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWithdrawOrderAllotLogic(ctx context.Context, svcCtx *svc.ServiceContext) WithdrawOrderAllotLogic {
	return WithdrawOrderAllotLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WithdrawOrderAllotLogic) WithdrawOrderAllot(req types.WithdrawOrderAllotRequest) (*types.WithdrawOrderAllotResponse, error) {
	// 1.查询订单
	withdrawOrder, err := model.NewMerchantWithdrawOrderModel(l.svcCtx.DbEngine).FindByOrderNo(req.OrderNo)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("订单[%v]不存在", req.OrderNo)
			return nil, common.NewCodeError(common.OrderNotExist)
		} else {
			l.Errorf("查询订单[%v]失败, err=%v", req.OrderNo, err)
			return nil, common.NewCodeError(common.SysDBSave)
		}

	}

	if err := l.checkOrderStatus(withdrawOrder.OrderNo, withdrawOrder.OrderStatus); err != nil {
		return nil, err
	}

	// 2.查询提现商户
	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneById(withdrawOrder.MerchantId)
	if err != nil {
		l.Errorf("查询提现商户[%v]失败, err=%v", withdrawOrder.MerchantId, err)
		return nil, common.NewCodeError(common.SystemInternalErr)
	}

	// 3.获取商户通道
	mchChannel, err := l.getMerchantChannel(merchant.Id, merchant.AreaId, merchant.Currency)
	if err != nil {
		l.Errorf("获取商户通道失败：%v", err)
		return nil, common.NewCodeError(common.SysDBSave)
	}

	// 4.获取payapi站点代付接口地址
	transferUrl, err := l.getTransferApiUrl()
	if err != nil {
		l.Errorf("取payapi站点域名：%v", err)
		return nil, common.NewCodeError(common.SysDBSave)
	}

	// 5.发送代付请求
	param := utils.TransferApiRequest{
		MerchantNo:         merchant.MerchantNo,
		Timestamp:          time.Now().Unix(),
		Amount:             withdrawOrder.RealAmount,
		Currency:           withdrawOrder.Currency,
		MchOrderNo:         withdrawOrder.OrderNo,
		TradeType:          mchChannel.PlatformChannelCode,
		OrderSource:        model.TransferOrderSourceWithdrawAllot,
		BankName:           withdrawOrder.BankName,
		BankCardHolderName: withdrawOrder.PayeeName,
		BankCardNo:         withdrawOrder.CardNumber,
		BankBranchName:     withdrawOrder.BranchName,
		BankCode:           withdrawOrder.BankCode,
	}
	resp, err := utils.SendTransfer(transferUrl, param, merchant.Md5Key)
	if err != nil {
		l.Errorf("请求失败, param: %+v, err=%v", param, err)
		if err := l.AllotFail(withdrawOrder); err != nil {
			l.Errorf("修改提现订单[%v]状态失败, err=%v", withdrawOrder.OrderNo, err)
		}
		return nil, common.NewCodeError(common.SysDBSave)
	}

	// 6.提现订单派单中
	if err := model.NewMerchantWithdrawOrderModel(l.svcCtx.DbEngine).UpdateStatusToAllot(withdrawOrder.Id, resp.Data.OrderNo); err != nil {
		l.Errorf("同步提现订单状态失败, orderId:%v, err:%v", withdrawOrder.Id, err)
		return nil, common.NewCodeError(common.SysDBSave)
	}

	return &types.WithdrawOrderAllotResponse{}, nil
}

// 检查提现订单状态
func (l *WithdrawOrderAllotLogic) checkOrderStatus(orderNo string, orderStatus int64) error {
	if orderStatus == model.WithdrawOrderStatusSuccess {
		l.Errorf("订单[%v]已经提现成功", orderNo)
		return common.NewCodeError(common.OrderCompleted)
	}

	if orderStatus == model.WithdrawOrderStatusAllot {
		l.Errorf("订单[%v]已经派单，等待上游返回结果", orderNo)
		return common.NewCodeError(common.OrderDispatch)
	}

	if orderStatus == model.WithdrawOrderStatusAllotSuccess {
		l.Errorf("订单[%v]已经派单成功，不可重复派单", orderNo)
		return common.NewCodeError(common.OrderDispatchSuccess)
	}

	if orderStatus == model.WithdrawOrderStatusAllotFail {
		l.Errorf("订单[%v]已经派单失败, 不可以再次派单", orderNo)
		return common.NewCodeError(common.OrderNotOp)
	}

	if orderStatus != model.WithdrawOrderStatusPass {
		l.Errorf("订单[%v]当前状态[%v]不支持当前操作", orderNo, orderStatus)
		return common.NewCodeError(common.OrderNotOp)
	}

	return nil
}

// 选择平台商户通道
func (l *WithdrawOrderAllotLogic) getMerchantChannel(merchantId, areaId int64, currency string) (model.MerchantChannelList, error) {
	mchChList, err := model.NewMerchantChannelModel(l.svcCtx.DbEngine).FindPlfMchTransferChannel(merchantId, areaId, currency)
	if err != nil {
		return model.MerchantChannelList{}, errors.New(fmt.Sprintf("获取商户通道失败：%v", err))
	}

	if len(mchChList) == 0 {
		return model.MerchantChannelList{}, errors.New("没有可用的通道")
	}

	idx := utils.RandInt64(0, int64(len(mchChList)-1))

	return mchChList[idx], nil
}

// 获取代付接口地址
func (l *WithdrawOrderAllotLogic) getTransferApiUrl() (string, error) {
	host, err := model.NewGlobalConfigModel(l.svcCtx.DbEngine).FindValueByKey(model.ConfigPayapiHostAddr)
	if err != nil {
		return "", errors.New(fmt.Sprintf("查询全局配置[%v]失败：%v", model.ConfigPayapiHostAddr, err))
	}

	if strings.TrimSpace(host) == "" {
		l.Errorf("", model.ConfigPayapiHostAddr, err)
		return "", errors.New(fmt.Sprintf("系统没有配全局配置[%v]", model.ConfigPayapiHostAddr))
	}

	return strings.TrimRight(host, "/") + utils.PayApiTransferPath, nil
}

// 派单失败
func (l *WithdrawOrderAllotLogic) AllotFail(withdrawOrder *model.MerchantWithdrawOrder) error {
	txErr := l.svcCtx.DbEngine.Transaction(func(tx *gorm.DB) error {
		// 1.修改提现订单状态
		if err := model.NewMerchantWithdrawOrderModel(tx).UpdateStatusToAllotFail(withdrawOrder.Id); err != nil {
			return err
		}

		// 2.加提现商户余额，减冻结金额
		log1 := model.WalletLogExt{
			BusinessNo: withdrawOrder.OrderNo,
			Source:     model.AmountSourceWithdraw,
			Remark:     "Withdrawal Failed",
		}
		if err := model.NewMerchantModel(tx).PlusBalanceUnfreezeTx(withdrawOrder.MerchantId, withdrawOrder.DecreaseAmount, log1); err != nil {
			return err
		}

		// 3.xcxcx

		return nil
	})

	return txErr
}
