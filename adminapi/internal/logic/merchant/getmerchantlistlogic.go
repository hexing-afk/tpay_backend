package merchant

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetMerchantListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMerchantListLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetMerchantListLogic {
	return GetMerchantListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMerchantListLogic) GetMerchantList(req types.GetMerchantListRequest) (*types.GetMerchantListResponse, error) {

	f := model.FindMerchantList{
		Username:        req.Username,
		ContactDetails:  req.ContactDetails,
		Currency:        req.Currency,
		StartCreateTime: req.StartCreateTime,
		EndCreateTime:   req.EndCreateTime,
		Page:            req.Page,
		PageSize:        req.PageSize,
	}
	list, total, err := model.NewMerchantModel(l.svcCtx.DbEngine).FindList(f)
	if err != nil {
		l.Errorf("查询商户列表数据失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var merchantList []types.Merchant
	for _, v := range list {
		merchantList = append(merchantList, types.Merchant{
			MerchantId:   v.Id,
			MerchantNo:   v.MerchantNo,
			Username:     v.Username,
			Phone:        v.Phone,
			Email:        v.Email,
			Currency:     v.Currency,
			CreateTime:   v.CreateTime,
			Status:       v.Status,
			Balance:      v.Balance,
			IpWhiteList:  v.IpWhiteList,
			TotpSecret:   v.TotpSecret,
			AreaName:     v.AreaName,
			FrozenAmount: v.FrozenAmount,
		})
	}

	return &types.GetMerchantListResponse{
		Total: total,
		List:  merchantList,
	}, nil
}
