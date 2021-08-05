package public

import (
	"context"
	"fmt"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/merchantapi/internal/svc"
	"tpay_backend/merchantapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type PlatformBankCardListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPlatformBankCardListLogic(ctx context.Context, svcCtx *svc.ServiceContext) PlatformBankCardListLogic {
	return PlatformBankCardListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PlatformBankCardListLogic) PlatformBankCardList(merchantId int64, req types.PlatformBankCardListReq) (*types.PlatformBankCardListReply, error) {
	merchant, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindOneById(merchantId)
	if err != nil {
		l.Errorf("查询商户[%v]信息失败, err=%v", merchantId, err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	q := model.FindPlatformBankCardList{
		Page:     req.Page,
		PageSize: req.PageSize,
		Currency: merchant.Currency,
		Status:   model.PlatformBankCardEnable,
	}
	dataList, total, err := model.NewPlatformBankCardModel(l.svcCtx.DbEngine).FindList(q)
	if err != nil {
		l.Errorf("查询平台收款卡列表失败,err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	imageUrl, err := model.NewGlobalConfigModel(l.svcCtx.DbEngine).FindValueByKey(model.ConfigImageBaseUrl)
	if err != nil {
		l.Errorf("查询图片域名地址失败[%v], err=%v", model.ConfigImageBaseUrl, err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var list []types.PlatformBankCardList
	for _, data := range dataList {
		d := types.PlatformBankCardList{
			Id:          data.Id,
			BankName:    data.BankName,    //	银行名称
			AccountName: data.AccountName, //	开户名
			CardNumber:  data.CardNumber,  //	银行卡号
			BranchName:  data.BranchName,  //	支行名称
			Currency:    data.Currency,    //	币种
			MaxAmount:   data.MaxAmount,   //	最大收款额度
			QrCode:      data.QrCode,      //	收款二维码
		}

		if d.QrCode != "" {
			d.QrCode = fmt.Sprintf("%s/%s", imageUrl, data.QrCode)
		}

		list = append(list, d)
	}

	return &types.PlatformBankCardListReply{
		Total: total,
		List:  list,
	}, nil
}
