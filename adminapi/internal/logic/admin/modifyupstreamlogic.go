package admin

import (
	"context"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type ModifyUpstreamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewModifyUpstreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) ModifyUpstreamLogic {
	return ModifyUpstreamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ModifyUpstreamLogic) ModifyUpstream(req types.ModifyUpstreamRequest) error {
	exist, err := model.NewUpstreamModel(l.svcCtx.DbEngine).CheckById(req.UpstreamId)
	if err != nil {
		l.Errorf("查询上游[%v]是否存在失败, err=%v", req.UpstreamId, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	if !exist {
		l.Errorf("上游[%v]不存在", req.UpstreamId)
		return common.NewCodeError(common.UpdateContentNotExist)
	}

	upInfo, err := model.NewUpstreamModel(l.svcCtx.DbEngine).FindOneByUpstreamMerchantNo(req.UpstreamMerchantNo)
	if err != nil && err != model.ErrRecordNotFound {
		l.Errorf("确认上游账号uid、上游代码是否唯一失败, req=%v, err=%v", req, err)
		return common.NewCodeError(common.SysDBAdd)
	}

	b, err := model.NewAreaModel(l.svcCtx.DbEngine).Check(req.AreaId)
	if err != nil {
		l.Errorf("确认地区是否存在出错, req=%v, err=%v", req, err)
		return common.NewCodeError(common.AreaNotExist)
	}

	if !b {
		l.Errorf("地区不存在 req=%v", req)
		return common.NewCodeError(common.AreaNotExist)
	}

	if upInfo != nil && upInfo.Id != req.UpstreamId {
		l.Errorf("与另一个上游配置[%v] 上游账号uid、上游代码重复, req=%v,  err=%v", upInfo.Id, req, err)
		return common.NewCodeError(common.SysDBAdd)
	}

	data := model.Upstream{
		UpstreamName:       req.UpstreamName,
		CallConfig:         req.CallConfig,
		UpstreamMerchantNo: req.UpstreamMerchantNo,
		UpstreamCode:       req.UpstreamCode,
		AreaId:             req.AreaId,
	}

	if err := model.NewUpstreamModel(l.svcCtx.DbEngine).Update(req.UpstreamId, data); err != nil {
		l.Errorf("修改上游[%v]失败, err=%v", req.UpstreamId, err)
		return common.NewCodeError(common.SysDBUpdate)
	}

	return nil
}
