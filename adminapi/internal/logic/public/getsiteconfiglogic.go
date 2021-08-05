package logic

import (
	"context"
	"encoding/json"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type GetSiteConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSiteConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) GetSiteConfigLogic {
	return GetSiteConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSiteConfigLogic) GetSiteConfig() (*types.GetSiteConfigResponse, error) {
	configValue, err := model.NewGlobalConfigModel(l.svcCtx.DbEngine).FindValueByKey(model.ConfigSiteConfig)
	if err != nil {
		l.Errorf("查询网站配置失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	if configValue == "" {
		l.Errorf("查询网站配置为空", err)
		return nil, nil
	}

	var site SiteConfig
	if err := json.Unmarshal([]byte(configValue), &site); err != nil {
		l.Errorf("网站配置json解析失败, err=%v", err)
		return nil, common.NewCodeError(common.SysDBGet)
	}

	return &types.GetSiteConfigResponse{
		SiteName: site.SiteName,
		SiteLogo: site.SiteLogo,
		SiteLang: site.SiteLang,
	}, nil
}
