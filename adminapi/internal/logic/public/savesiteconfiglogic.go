package logic

import (
	"context"
	"encoding/json"
	"time"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/model"

	"tpay_backend/adminapi/internal/svc"
	"tpay_backend/adminapi/internal/types"

	"github.com/tal-tech/go-zero/core/logx"
)

type SaveSiteConfigLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type SiteConfig struct {
	SiteName string
	SiteLogo string
	SiteLang string
}

func NewSaveSiteConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) SaveSiteConfigLogic {
	return SaveSiteConfigLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveSiteConfigLogic) SaveSiteConfig(req types.SaveSiteConfigRequest) error {
	site := SiteConfig{
		SiteName: req.SiteName,
		SiteLogo: req.SiteLogo,
		SiteLang: req.SiteLang,
	}
	bytes, err := json.Marshal(site)
	if err != nil {
		l.Errorf("网站配置数据转JSON失败, data=%+v, err=%v", site, err)
		return common.NewCodeError(common.SysDBSave)
	}

	data := model.GlobalConfig{
		ConfigKey:   model.ConfigSiteConfig,
		ConfigValue: string(bytes),
		Remark:      "网站配置",
		CreateTime:  time.Now().Unix(),
		IsChange:    model.IsChangeTrue,
	}
	//_, err = model.NewGlobalConfigModel(l.svcCtx.MysqlConn).InsertOrUpdate(data)
	err = model.NewGlobalConfigModel(l.svcCtx.DbEngine).InsertOrUpdate(data)
	if err != nil {
		l.Errorf("保存网站配置失败, err=%v", err)
		return common.NewCodeError(common.SysDBSave)
	}

	return nil
}
