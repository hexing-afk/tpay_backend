package order

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type TransferOrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	userId int64
}

func NewTransferOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext, userId int64) TransferOrderDetailLogic {
	return TransferOrderDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		userId: userId,
	}
}

func (l *TransferOrderDetailLogic) TransferOrderDetail(req types.TransferOrderDetailRequest) (*types.TransferOrderDetailReply, error) {
	//根据商户id查询商户编号
	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneById(l.userId)
	if err != nil {
		l.Errorf("查询商户出错, userId[%v], err[%v]", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	order, err := model.NewTransferOrderModel(l.svcCtx.DbEngine).FindByOrderNoAndMerchantNo(req.OrderNo, merchant.MerchantNo)
	if err != nil {
		l.Errorf("查询订单出错, MerchantNo:%v,OrderNo:%v err:%v", merchant.MerchantNo, req.OrderNo, err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	data := types.TransferOrderDetailDetail{
		OrderNo:         order.OrderNo,         //平台订单号
		MerchantOrderNo: order.MerchantOrderNo, //商户订单号
		Currency:        order.Currency,        //币种
		OrderAmount:     order.ReqAmount,       //订单请求金额
		MerchantFee:     order.MerchantFee,     //商户手续费
		DecreaseAmount:  order.DecreaseAmount,  //账户扣除的金额
		FeeDeductType:   order.FeeDeductType,   //手续费扣款方式(1内扣,2外扣)
		PayeeRealAmount: order.PayeeRealAmount, //收款人到账金额
		CreateTime:      order.CreateTime,      //创建时间
		UpdateTime:      order.UpdateTime,      //更新时间
		BankName:        order.BankName,        //收款银行名称
		AccountName:     order.AccountName,     //银行卡开户名
		CardNumber:      order.CardNumber,      //收款卡号
		BranchName:      order.BranchName,      //支行名称
		OrderSource:     order.OrderSource,
		ChannelName:     order.PlatformChannelName,
	}

	return &types.TransferOrderDetailReply{
		Data: data,
	}, nil
}
