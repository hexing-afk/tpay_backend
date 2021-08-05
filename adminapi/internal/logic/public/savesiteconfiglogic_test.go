package logic

import (
	"encoding/json"
	"testing"
)

func TestSaveSiteConfigLogic_SaveSiteConfig(t *testing.T) {
	site := SiteConfig{
		SiteName: "Seastarpay",
		SiteLogo: "logo",
		SiteLang: "",
	}
	bytes, err := json.Marshal(site)
	if err != nil {
		t.Errorf("网站配置数据转JSON失败, err=%v", err)
		return
	}
	t.Logf("json: %v", string(bytes))
}
