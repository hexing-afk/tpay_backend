package order

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetPayOrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPayOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetPayOrderDetailLogic {
	return GetPayOrderDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPayOrderDetailLogic) GetPayOrderDetail(req types.GetPayOrderDetailRequest) (*types.GetPayOrderDetailResponse, error) {
	order, err := model.NewPayOrderModel(l.svcCtx.DbEngine).FindOneDetailByOrderNo(req.OrderNo)
	if err != nil {
		l.Errorf("查询订单详情失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	dataInfo := types.PayOrderDetail{
		OrderNo:             order.OrderNo, //平台订单号
		MerchantOrderNo:     order.MerchantOrderNo,
		UpstreamOrderNo:     order.UpstreamOrderNo,
		MerchantName:        order.MerchantName,
		MerchantNo:          order.MerchantNo,
		ReqAmount:           order.ReqAmount,      //请求上游的金额
		IncreaseAmount:      order.IncreaseAmount, //账户增加的金额
		MerchantFee:         order.MerchantFee,    //商户手续费
		UpstreamFee:         order.UpstreamFee,    //上游手续费
		OrderStatus:         order.OrderStatus,    //订单状态:1-待支付;2-已支付;3-超时
		Currency:            order.Currency,
		CreateTime:          order.CreateTime,
		UpdateTime:          order.UpdateTime,
		NotifyUrl:           order.NotifyUrl,           //异步通知url
		ReturnUrl:           order.ReturnUrl,           //同步跳转url
		PlatformChannelId:   order.PlatformChannelId,   //平台通道id
		UpstreamChannelId:   order.UpstreamChannelId,   //上游通道id
		NotifyStatus:        order.NotifyStatus,        //异步通知状态(0未通知,1成功,2通知进行中,3超时)
		Subject:             order.Subject,             //商品的标题/交易标题/订单标题/订单关键字等
		PlatformChannelName: order.PlatformChannelName, //平台通道名称
		UpstreamChannelName: order.UpstreamChannelName, //上游通道名称
		PaymentAmount:       order.PaymentAmount,
	}

	return &types.GetPayOrderDetailResponse{
		Data: dataInfo,
	}, nil
}
