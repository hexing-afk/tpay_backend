package order

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetTransferOrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTransferOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetTransferOrderDetailLogic {
	return GetTransferOrderDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTransferOrderDetailLogic) GetTransferOrderDetail(req types.GetTransferOrderDetailRequest) (*types.GetTransferOrderDetailResponse, error) {

	data, err := model.NewTransferOrderModel(l.svcCtx.DbEngine).FindDetail(req.OrderNo)
	if err != nil {
		l.Errorf("查询代付订单失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	dataInfo := types.TransferOrderDetail{
		OrderNo:         data.OrderNo,
		MerchantOrderNo: data.MerchantOrderNo,
		UpstreamOrderNo: data.UpstreamOrderNo,
		MerchantName:    data.MerchantName,
		MerchantNo:      data.MerchantNo,
		ReqAmount:       data.ReqAmount,       //订单请求金额
		DecreaseAmount:  data.DecreaseAmount,  //账户扣除的金额
		MerchantFee:     data.MerchantFee,     //	商户手续费
		UpstreamFee:     data.UpstreamFee,     //上游手续费
		PayeeRealAmount: data.PayeeRealAmount, //收款方实际到账金额
		FeeDeductType:   data.FeeDeductType,   //手续费扣款方式(1内扣,2外扣)
		UpstreamAmount:  data.UpstreamAmount,  //请求上游的金额
		Currency:        data.Currency,
		CreateTime:      data.CreateTime,
		UpdateTime:      data.UpdateTime,
		NotifyUrl:       data.NotifyUrl,    //	异步通知地址
		BankName:        data.BankName,     //	收款银行名称
		CardNumber:      data.CardNumber,   //	收款卡号
		AccountName:     data.AccountName,  //	银行卡开户名
		BranchName:      data.BranchName,   //支行名称
		NotifyStatus:    data.NotifyStatus, //	异步通知状态(0未通知,1成功,2通知进行中,3超时)
		OrderSource:     data.OrderSource,
	}

	return &types.GetTransferOrderDetailResponse{
		Data: dataInfo,
	}, nil
}
