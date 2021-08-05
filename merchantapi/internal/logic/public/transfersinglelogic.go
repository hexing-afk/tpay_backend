package public

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"
	"tpay_backend/utils"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type TransferSingleLogic struct {
	logx.Logger
	ctx        context.Context
	svcCtx     *svc.ServiceContext
	merchantId int64
}

func NewTransferSingleLogic(ctx context.Context, svcCtx *svc.ServiceContext, merchantId int64) TransferSingleLogic {
	return TransferSingleLogic{
		Logger:     logx.WithContext(ctx),
		ctx:        ctx,
		svcCtx:     svcCtx,
		merchantId: merchantId,
	}
}

func (l *TransferSingleLogic) TransferSingle(req types.TransferSingleReq) error {
	// 1.查询商户信息
	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneById(l.merchantId)
	if err != nil {
		l.Errorf("查询商户[%v]信息失败, err=%v", l.merchantId, err)
		return common.NewCodeError(common.SysDBAdd)
	}

	if merchant.Balance < req.ReqAmount {
		l.Errorf("商户[%v]余额[%v]不足, reqAmount=%v", merchant.Balance, req.ReqAmount)
		return common.NewCodeError(common.InsufficientBalance)
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

	// 2.支付密码验证
	//支付密码解密
	payPassword, err := common.DecryptPassword(req.PayPassword)
	if err != nil {
		l.Errorf("密码解密发生错误,err:%v, password:%v", err, req.PayPassword)
		return common.NewCodeError(common.SysDBErr)
	}

	// 支付密码比对密码
	if merchant.PayPassword != common.CreateMerchantPayPassword(payPassword) {
		l.Error("修改支付密码失败:支付密码校验失败")
		return common.NewCodeError(common.PayPasswordErr)
	}

	// 3.根据平台通道id选取商户通道
	mchChannel, err := l.getMerchantChannel(merchant.Currency, merchant.AreaId, req.ChannelId)
	if err != nil {
		l.Errorf("获取商户通道失败, err=%v", err)
		return common.NewCodeError(common.NotChannelAvailable)
	}

	merchantFee := utils.CalculatePayOrderFeeMerchant(req.ReqAmount, mchChannel.SingleFee, mchChannel.Rate)
	if merchant.Balance < req.ReqAmount+merchantFee {
		l.Errorf("商户[%v]余额[%v]不足, decreaseAmount=%v", merchant.Balance, req.ReqAmount+merchantFee)
		return common.NewCodeError(common.InsufficientBalance)
	}

	// 4.获取代付URL
	url, err := l.getTransferApiUrl()
	if err != nil {
		l.Errorf("获取代付URL失败, err=%v", err)
		return common.NewCodeError(common.SysDBAdd)
	}

	mchOrderNo := GenerateMerchantOrderNo(utils.GetDailyId())
	param := utils.TransferApiRequest{
		MerchantNo:         merchant.MerchantNo,
		Timestamp:          time.Now().Unix(),
		Amount:             req.ReqAmount,
		Currency:           merchant.Currency,
		MchOrderNo:         mchOrderNo,
		TradeType:          mchChannel.PlatformChannelCode,
		OrderSource:        model.TransferOrderSourceMerchantPayment,
		BankName:           req.BankName,
		BankCardHolderName: req.AccountName,
		BankCardNo:         req.CardNumbers,
		BankBranchName:     req.BankBranchName,
		BankCode:           req.BankCode,
		Remark:             req.Remark,
	}

	resp, err := utils.SendTransfer(url, param, merchant.Md5Key)
	if err != nil {
		l.Errorf("发送代付请求失败, err=%v", err)
		return common.NewCodeError(common.SysDBAdd)
	}

	l.Infof("代付接口返回：%+v", resp)

	return nil
}

// 选择商户通道
func (l *TransferSingleLogic) getMerchantChannel(currency string, areaId, channelId int64) (model.MerchantChannelList, error) {
	mchChannel, err := model.NewMerchantChannelModel(l.svcCtx.DbEngine).FindPlfMchTransferChannelByPlfChannelId(l.merchantId, areaId, currency, channelId)
	if err != nil {
		return model.MerchantChannelList{}, errors.New(fmt.Sprintf("获取商户通道失败：%v", err))
	}

	if &mchChannel == nil {
		return model.MerchantChannelList{}, errors.New("没有可用的通道")
	}

	return mchChannel, nil
}

// 获取代付接口地址
func (l *TransferSingleLogic) getTransferApiUrl() (string, error) {
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

// 生成平台商户付款单号
func GenerateMerchantOrderNo(orderNo string) string {
	return fmt.Sprintf("%s%s", "MPAY", orderNo)
}
