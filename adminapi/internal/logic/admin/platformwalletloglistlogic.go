package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type PlatformWalletLogListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPlatformWalletLogListLogic(ctx context.Context, svcCtx *svc.ServiceContext) PlatformWalletLogListLogic {
	return PlatformWalletLogListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PlatformWalletLogListLogic) PlatformWalletLogList(req types.PlatformWalletLogListRequest) (*types.PlatformWalletLogListReply, error) {

	q := model.PlatformWalletLogListReq{
		BusinessNo: req.BusinessNo,
		Source:     req.Source,
		Currency:   req.Currency,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	datas, total, err := model.NewPlatformWalletLogModel(l.svcCtx.DbEngine).FindList(q)
	if err != nil {
		l.Errorf("查询平台收益列表失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	var list []types.PlatformWalletLogListData
	for _, data := range datas {
		list = append(list, types.PlatformWalletLogListData{
			Id:          data.Id,
			BusinessNo:  data.BusinessNo,
			Source:      data.Source,
			MerchantFee: data.MerchantFee,
			UpstreamFee: data.UpstreamFee,
			Income:      data.Income,
			CreateTime:  data.CreateTime,
			Currency:    data.Currency,
		})
	}

	return &types.PlatformWalletLogListReply{
		Total: total,
		List:  list,
	}, nil
}
