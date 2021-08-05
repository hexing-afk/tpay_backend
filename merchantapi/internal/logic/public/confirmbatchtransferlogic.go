package public

import (
	"context"
	"time"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"
	"tpay_backend/utils"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ConfirmBatchTransferLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfirmBatchTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) ConfirmBatchTransferLogic {
	return ConfirmBatchTransferLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfirmBatchTransferLogic) ConfirmBatchTransfer(merchantId int64, req types.ConfirmBatchTransferRequest) error {
	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneById(merchantId)
	if err != nil {
		l.Errorf("查询商户信息失败,err:%v", err)
		return common.NewCodeError(common.SysDBGet)
	}

	// google验证码是否关闭
	totpIsClose, err := model.NewGlobalConfigModel(l.svcCtx.DbEngine).TotpIsClose()
	if err != nil {
		l.Errorf("查询totp配置失败:%v", err)
		return common.NewCodeError(common.SysDBErr)
	}

	if !totpIsClose { // 没有关闭
		// 验证TOTP密码code
		if !utils.VerifyTOTPPasscode(req.TotpCode, merchant.TotpSecret) {
			return common.NewCodeError(common.LoginCaptchaNotMatch)
		}
	}

	payPassword, err := common.DecryptPassword(req.PayPassword)
	if err != nil {
		l.Errorf("密码解密发生错误,err:%v, password:%v", err, req.PayPassword)
		return common.NewCodeError(common.TransferFail)
	}

	// 支付密码比对密码
	if merchant.PayPassword != common.CreateMerchantPayPassword(payPassword) {
		l.Error("支付密码校验失败")
		return common.NewCodeError(common.PayPasswordErr)
	}

	batchOrder, err := model.NewTransferBatchOrderModel(l.svcCtx.DbEngine).FindByBatchNo(merchantId, req.BatchNo)
	if err != nil {
		l.Errorf("查询批量付款订单信息失败,BatchNo:%v; err:%v", req.BatchNo, err)
		return common.NewCodeError(common.FileNotExist)
	}

	if batchOrder.MerchantId != merchantId {
		l.Errorf("批量付款订单不是该商户的,不允许进行操作")
		return common.NewCodeError(common.FileNotExist)
	}

	// 2.获取payapi站点批量代付接口地址
	transferBatchUrl, err := NewFuncLogic(l.ctx, l.svcCtx).SystemTransferBatchApiUrl()
	if err != nil {
		l.Errorf("取payapi站点批量代付接口地址失败：%v", err)
		return common.NewCodeError(common.TransferFail)
	}

	postReq := utils.TransferBatchApiRequest{
		MerchantNo: merchant.MerchantNo, // 商户编号
		Timestamp:  time.Now().Unix(),   // 请求时间
		BatchNo:    batchOrder.BatchNo,  // 批量号
	}

	// 5.发送请求到pay服务
	resp, err := utils.SendTransferBatch(transferBatchUrl, postReq, merchant.Md5Key)
	if err != nil {
		l.Errorf("发送批量付款请求失败,err:%v", err)
		//return common.NewCodeErrorWithMsg(common.TransferFail, resp.Msg)
		return common.NewCodeErrorWithMsg(common.TransferFail, err.Error())
	}

	if resp.Code != 0 {
		l.Errorf("请求失败, param: %+v, err=%v", postReq, err)
		return common.NewCodeErrorWithMsg(common.TransferFail, resp.Msg)
	}

	return nil
}
