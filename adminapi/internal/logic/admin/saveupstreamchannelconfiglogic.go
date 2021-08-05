package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type SaveUpstreamChannelConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveUpstreamChannelConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) SaveUpstreamChannelConfigLogic {
	return SaveUpstreamChannelConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveUpstreamChannelConfigLogic) SaveUpstreamChannelConfig(req types.SaveUpstreamChannelConfigRequest) error {
	channel, err := model.NewUpstreamChannelModel(l.svcCtx.DbEngine).FindOneById(req.ChannelId)
	if err != nil {
		if err == model.ErrRecordNotFound {
			l.Errorf("通道[%v]不存在", req.ChannelId)
			return common.NewCodeError(common.ChannelNotExist)
		} else {
			l.Errorf("查询通道[%v]失败, err=%v", req.ChannelId, err.Error())
			return common.NewCodeError(common.SysDBUpdate)
		}
	}

	if channel.Status == model.UpstreamChannelStatusDisable {
		l.Errorf("通道[%v]已被禁用", req.ChannelId)
		return common.NewCodeError(common.ChannelDisable)
	}

	// 如果通道类型是代收，扣手续费方式只能是内扣
	// 如果通道类型是代付，扣手续费方式只能是外扣
	switch channel.ChannelType {
	case model.UpstreamChannelTypeCollection:
		req.DeductionMethod = model.UpstreamChannelDeductionInner
		//if req.DeductionMethod != model.UpstreamChannelDeductionInner {
		//	l.Errorf("通道[%v]手续费扣费方式只能是内扣, req.DeductionMethod=%v", channel.Id, req.DeductionMethod)
		//	return common.NewCodeError(common.PayChannelIsInnerDeduction)
		//}
	case model.UpstreamChannelTypeTransfer:
		req.DeductionMethod = model.UpstreamChannelDeductionOut
		//if req.DeductionMethod != model.UpstreamChannelDeductionOut {
		//	l.Errorf("通道[%v]手续费扣费方式只能是外扣, req.DeductionMethod=%v", channel.Id, req.DeductionMethod)
		//	return common.NewCodeError(common.PayOutChannelIsOutDeduction)
		//}
	}

	data := model.UpstreamChannel{
		DeductionMethod: req.DeductionMethod,
		Rate:            req.Rate,
		SingleFee:       req.SingleFee,
		SingleMaxAmount: req.SingleMaxAmount,
	}

	if err := model.NewUpstreamChannelModel(l.svcCtx.DbEngine).Update(req.ChannelId, data); err != nil {
		l.Errorf("修改通道[%v]配置失败, err=%v", req.ChannelId, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	return nil
}
