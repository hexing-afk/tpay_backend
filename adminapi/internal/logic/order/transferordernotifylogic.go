package order

import (
	"context"
	"time"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"
	"tpay_backend/utils"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type TransferOrderNotifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTransferOrderNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) TransferOrderNotifyLogic {
	return TransferOrderNotifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TransferOrderNotifyLogic) TransferOrderNotify(req types.TransferOrderNotifyRequest) (*types.TransferOrderNotifyResponse, error) {
	order, err := model.NewTransferOrderModel(l.svcCtx.DbEngine).FindByOrderNo(req.OrderNo)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("代付订单[%v]不存在", req.OrderNo)
			return nil, common.NewCodeError(common.OrderNotExist)
		} else {
			l.Errorf("查询代付订单[%]失败, err=%v", req.OrderNo, err)
			return nil, common.NewCodeError(common.SystemInternalErr)
		}
	}

	// 检查订单状态
	// 只有支付成功和支付失败的订单才可以通知
	if order.OrderStatus != model.TransferOrderStatusPaid && order.OrderStatus != model.TransferOrderStatusFail {
		l.Errorf("代付订单[%]当前支付状态[%v]不能进行通知", order.OrderNo, order.OrderStatus)
		return nil, common.NewCodeError(common.OrderNotOp)
	}

	// 检查订单是否有异步通知URL
	if order.NotifyUrl == "" {
		l.Errorf("代付订单[%v]缺少异步通知地址", order.OrderNo)
		return nil, common.NewCodeError(common.OrderMissingNotifyUrl)
	}

	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneByMerchantNo(order.MerchantNo)
	if err != nil {
		l.Errorf("查询商户失败, MerchantNo=%v, err=%v", order.MerchantNo, err)
		return nil, common.NewCodeError(common.SystemInternalErr)
	}

	postData := &utils.PackTransferNotifyParamsRequest{
		MerchantNo:      order.MerchantNo,
		Timestamp:       time.Now().Unix(),
		NotifyType:      utils.TransferNotifyType,
		OrderNo:         order.OrderNo,
		MerchantOrderNo: order.MerchantOrderNo,
		ReqAmount:       order.ReqAmount,
		Currency:        order.Currency,
		OrderStatus:     order.OrderStatus,
		PayTime:         order.UpdateTime,
	}

	param, err := utils.PackTransferNotifyParams(postData, merchant.Md5Key)
	if err != nil {
		l.Errorf("打包参数失败, data=%v, err=%v", order.MerchantNo, err)
		return nil, common.NewCodeError(common.SystemInternalErr)
	}

	body, err := utils.PostForm(order.NotifyUrl, param)
	if err != nil {
		l.Errorf("发送数据失败, url=%v, param=%v, err=%v", order.NotifyUrl, param, err)
		return nil, common.NewCodeError(common.OrderNotifyFail)
	}

	bodyStr := string(body)
	if bodyStr != "success" {
		l.Errorf("通知失败, body=%v", bodyStr)
	}

	l.Infof("代付订单[%v]通知成功, body:%v", order.OrderNo, bodyStr)

	return &types.TransferOrderNotifyResponse{
		NotifyResponse: bodyStr,
	}, nil
}
