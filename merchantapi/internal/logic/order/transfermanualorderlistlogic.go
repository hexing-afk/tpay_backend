package order

import (
	"context"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type TransferManualOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	userId int64
}

func NewTransferManualOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext, userId int64) TransferManualOrderListLogic {
	return TransferManualOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		userId: userId,
	}
}

func (l *TransferManualOrderListLogic) TransferManualOrderList(req types.TransferManualOrderListRequest) (*types.TransferManualOrderListReply, error) {
	//根据商户id查询商户编号
	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneById(l.userId)
	if err != nil {
		l.Errorf("查询商户出错, userId[%v], err[%v]", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	f := model.FindTransferOrderList{
		MerchantNo:      merchant.MerchantNo,
		Page:            req.Page,
		PageSize:        req.PageSize,
		StartCreateTime: req.StartCreateTime,
		EndCreateTime:   req.EndCreateTime,
		OrderNo:         req.OrderNo,
		OrderStatus:     req.OrderStatus,
		OrderSourceList: []int64{model.TransferOrderSourceMerchantPayment, model.TransferOrderSourceWithdrawAllot},
	}
	dataList, total, err := model.NewTransferOrderModel(l.svcCtx.DbEngine).FindList(f)
	if err != nil {
		l.Errorf("查询代付订单失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var list []types.TransferManualOrderList
	for _, data := range dataList {
		list = append(list, types.TransferManualOrderList{
			OrderNo:         data.OrderNo,             // 平台订单号
			Currency:        data.Currency,            // 币种
			OrderAmount:     data.ReqAmount,           // 订单请求金额
			MerchantFee:     data.MerchantFee,         // 商户手续费
			PayeeRealAmount: data.PayeeRealAmount,     // 收款人到账金额
			BankName:        data.BankName,            // 收款银行
			AccountName:     data.AccountName,         // 收款人姓名
			CardNumber:      data.CardNumber,          // 收款卡号
			BranchName:      data.BranchName,          // 支行名称
			OrderStatus:     data.OrderStatus,         // 订单状态: 1-待支付; 2-支付成功; 3-支付失败
			CreateTime:      data.CreateTime,          // 创建时间
			UpdateTime:      data.UpdateTime,          // 更新时间
			OrderSource:     data.OrderSource,         // 订单来源：1-接口; 2-平台提现派单；3-商户后台付款
			ChannelName:     data.PlatformChannelName, // 通道名称
			Remark:          data.Remark,              // 付款备注
		})
	}
	return &types.TransferManualOrderListReply{
		Total: total,
		List:  list,
	}, nil
}
