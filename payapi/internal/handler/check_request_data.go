package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"tpay_backend/model"
	"tpay_backend/payapi/internal/common"
	"tpay_backend/payapi/internal/svc"
	"tpay_backend/payapi/internal/types"
	"tpay_backend/utils"

	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest/httpx"
)

type RequestData struct {
	Data string `form:"data"`
	Sign string `form:"sign"`
}

const (
	RequestExpire int64 = 120    // 请求参数有效期(单位: 秒)
	IP                  = "1111" // 不收白名单限制的id
)

// 请求公共数据及签名验证-外部请求
func CheckRequestData(svcCtx *svc.ServiceContext, r *http.Request, req interface{}) (*model.Merchant, error) {
	var reqData RequestData
	if err := httpx.Parse(r, &reqData); err != nil {
		return nil, common.NewCodeErrorWithMsg(common.VerifyParamFailed, err.Error())
	}

	var commonField types.RequestCommonField
	if err := json.Unmarshal([]byte(reqData.Data), &commonField); err != nil {
		return nil, common.NewCodeError(common.DataFieldFormatErr)
	}

	if commonField.MerchantNo == "" {
		return nil, common.NewCodeErrorWithMsg(common.VerifyParamFailed, "merchant_no字段不能为空")
	}

	if commonField.Timestamp < 1 {
		return nil, common.NewCodeErrorWithMsg(common.VerifyParamFailed, "timestamp字段不能为空")
	}

	if commonField.Timestamp < (time.Now().Unix() - RequestExpire) {
		return nil, common.NewCodeErrorWithMsg(common.VerifyParamFailed, "请求参数已过期")
	}

	// 获取商户信息
	merchant, err := model.NewMerchantModel(svcCtx.DbEngine).FindOneByMerchantNo(commonField.MerchantNo)
	if err != nil {
		if err == model.ErrRecordNotFound {
			return nil, common.NewCodeError(common.MerchantNotExist)
		} else {
			logx.Errorf("获取商户信息失败:%v, MerchantNo:%v", err, commonField.MerchantNo)
			return nil, common.NewCodeError(common.SystemInternalErr)
		}
	}

	if strings.TrimSpace(merchant.IpWhiteList) != IP {
		// 检查ip白名单
		clientIp := utils.ClientIP(r)
		if !utils.InSlice(clientIp, merchant.IpWhiteListSlice) {
			logx.Errorf("检查ip白名单失败clientIp:%v, IpWhiteList:%v", clientIp, merchant.IpWhiteListSlice)
			return nil, common.NewCodeError(common.IpWhitListForbidden)
		}
	}

	// 检查签名
	if utils.GenerateSign(reqData.Data, merchant.Md5Key) != reqData.Sign {
		return nil, common.NewCodeError(common.SignFailed)
	}

	// 商户是否被禁用
	if merchant.Status != model.MerchantStatusEnable {
		return nil, common.NewCodeError(common.MerchantForbidden)
	}

	// 解析数据
	if err := json.Unmarshal([]byte(reqData.Data), &req); err != nil {
		return nil, common.NewCodeError(common.DataFieldFormatErr)
	}

	return merchant, nil
}

// 请求公共数据及签名验证-内部请求
func CheckInternalRequestData(svcCtx *svc.ServiceContext, r *http.Request, req interface{}) (*model.Merchant, error) {
	var reqData RequestData
	if err := httpx.Parse(r, &reqData); err != nil {
		return nil, common.NewCodeErrorWithMsg(common.VerifyParamFailed, err.Error())
	}

	var commonField types.RequestCommonField
	if err := json.Unmarshal([]byte(reqData.Data), &commonField); err != nil {
		return nil, common.NewCodeError(common.DataFieldFormatErr)
	}

	if commonField.MerchantNo == "" {
		return nil, common.NewCodeErrorWithMsg(common.VerifyParamFailed, "merchant_no字段不能为空")
	}

	if commonField.Timestamp < 1 {
		return nil, common.NewCodeErrorWithMsg(common.VerifyParamFailed, "timestamp字段不能为空")
	}

	if commonField.Timestamp < (time.Now().Unix() - RequestExpire) {
		return nil, common.NewCodeErrorWithMsg(common.VerifyParamFailed, "请求参数已过期")
	}

	// 获取商户信息
	merchant, err := model.NewMerchantModel(svcCtx.DbEngine).FindOneByMerchantNo(commonField.MerchantNo)
	if err != nil {
		if err == model.ErrRecordNotFound {
			return nil, common.NewCodeError(common.MerchantNotExist)
		} else {
			logx.Errorf("获取商户信息失败:%v, MerchantNo:%v", err, commonField.MerchantNo)
			return nil, common.NewCodeError(common.SystemInternalErr)
		}
	}

	// 检查签名
	if utils.GenerateSign(reqData.Data, merchant.Md5Key) != reqData.Sign {
		return nil, common.NewCodeError(common.SignFailed)
	}

	// 商户是否被禁用
	if merchant.Status != model.MerchantStatusEnable {
		return nil, common.NewCodeError(common.MerchantForbidden)
	}

	// 解析数据
	if err := json.Unmarshal([]byte(reqData.Data), &req); err != nil {
		return nil, common.NewCodeError(common.DataFieldFormatErr)
	}

	return merchant, nil
}
