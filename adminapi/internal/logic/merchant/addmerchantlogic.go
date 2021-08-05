package merchant

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"
	"tpay_backend/model"
	"tpay_backend/utils"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddMerchantLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddMerchantLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddMerchantLogic {
	return AddMerchantLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddMerchantLogic) AddMerchant(req types.AddMerchantRequest) (*types.AddMerchantResponse, error) {
	exist, err := model.NewMerchantModel(l.svcCtx.DbEngine).CheckByName(req.Username)
	if err != nil {
		l.Errorf("检查商户[%v]是否已存在失败, err=%v", req.Username, err)
		return nil, common.NewCodeError(common.SysDBAdd)
	}

	if exist {
		l.Errorf("商户[%v]已经存在", req.Username)
		return nil, common.NewCodeError(common.AccountRepeat)
	}

	//解密
	plainPassword, err := common.DecryptPassword(req.Password)
	if err != nil {
		l.Errorf("密码解密发生错误,err:%v, password:%v", err, req.Password)
		return nil, common.NewCodeError(common.SysDBErr)
	}

	// 生成TOTP秘钥
	totpSecret, err := utils.GenerateTOTPSecret(req.Username)
	if err != nil {
		l.Errorf("生成TOTP秘钥错误,err:%v, Username:%v", err, req.Username)
		return nil, common.NewCodeError(common.SystemInternalErr)
	}

	b, err := model.NewAreaModel(l.svcCtx.DbEngine).Check(req.AreaId)
	if err != nil {
		l.Errorf("确认地区是否存在出错, req=%v, err=%v", req, err)
		return nil, common.NewCodeError(common.AreaNotExist)
	}

	if !b {
		l.Errorf("地区不存在 req=%v", req)
		return nil, common.NewCodeError(common.AreaNotExist)
	}

	// 添加商户
	merchant := &model.Merchant{
		MerchantNo:   utils.GetUniqueId(),
		Username:     req.Username,
		Password:     common.CreateAdminPassword(plainPassword),
		Phone:        req.Phone,
		Email:        req.Email,
		Currency:     req.Currency,
		Status:       model.MerchantStatusEnable,
		Md5Key:       utils.RandString(32),
		Balance:      0,
		FrozenAmount: 0,
		PayPassword:  common.CreateMerchantPayPassword(common.MerchantDefaultPayPassword),
		TotpSecret:   totpSecret,
		AreaId:       req.AreaId,
	}

	if err := model.NewMerchantModel(l.svcCtx.DbEngine).Insert(merchant); err != nil {
		l.Errorf("添加商户[%v]失败, err=%v", req.Username, err)
		return nil, common.NewCodeError(common.SysDBAdd)
	}

	return &types.AddMerchantResponse{MerchantId: merchant.Id}, nil
}
