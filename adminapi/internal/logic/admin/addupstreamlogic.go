package admin

import (
	"context"
	"encoding/json"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type AddUpstreamLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddUpstreamLogic(ctx context.Context, svcCtx *svc.ServiceContext) AddUpstreamLogic {
	return AddUpstreamLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddUpstreamLogic) AddUpstream(req types.AddUpstreamRequest) error {
	if req.CallConfig != "" {
		_, err := json.Marshal(req.CallConfig)
		if err != nil {
			l.Errorf("上游通信配置JSON解析失败, req.CallConfig=%v, err=%v", req.CallConfig, err)
			return common.NewCodeError(common.InvalidParam)
		}
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

	upInfo, err := model.NewUpstreamModel(l.svcCtx.DbEngine).FindOneByUpstreamMerchantNo(req.UpstreamMerchantNo)
	if err != nil && err != model.ErrRecordNotFound {
		l.Errorf("确认上游账号uid是否唯一失败, req=%v, err=%v", req, err)
		return common.NewCodeError(common.AccountRepeat)
	}

	if upInfo != nil && upInfo.Id != 0 {
		l.Errorf("与另一个上游配置[%v] 上游账号uid 重复 req=%v,  err=%v", upInfo.Id, req, err)
		return common.NewCodeError(common.AccountRepeat)
	}

	data := &model.Upstream{
		UpstreamName:       req.UpstreamName,
		CallConfig:         req.CallConfig,
		UpstreamMerchantNo: req.UpstreamMerchantNo,
		UpstreamCode:       req.UpstreamCode,
		AreaId:             req.AreaId,
	}

	if err := model.NewUpstreamModel(l.svcCtx.DbEngine).Insert(data); err != nil {
		l.Errorf("添加上游[%v]失败, err=%v", req.UpstreamName, err)
		return common.NewCodeError(common.SysDBAdd)
	}

	return nil
}
