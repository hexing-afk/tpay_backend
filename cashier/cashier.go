package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"tpay_backend/cashier/internal/config"
	"tpay_backend/cashier/internal/handler"
	"tpay_backend/cashier/internal/lang"
	"tpay_backend/cashier/internal/svc"
	"tpay_backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/tal-tech/go-zero/core/conf"
)

//go:embed templates/*
var templateFiles embed.FS

//go:embed static/*
var staticFiles embed.FS

var configFile = flag.String("f", "cashier/etc/cashier.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	// 加载配置
	conf.MustLoad(*configFile, &c)

	// 设置时区
	utils.SetTimezone(c.Timezone)

	// 重置redis的数据库
	c.Redis.DB = utils.RedisDbCashier
	sctx := svc.NewServiceContext(c)

	engine := gin.Default()

	// 3.添加静态文件目录
	// 访问示例: /assets/static/img/apple.png
	engine.StaticFS("/assets", http.FS(staticFiles))

	// 4.解析模板并添加
	templ := template.Must(template.New("").Funcs(template.FuncMap{
		"Lang": lang.Lang, // 设置模板函数(在模板中使用的函数)
	}).ParseFS(templateFiles, "templates/*.html"))

	engine.SetHTMLTemplate(templ)

	// 注册服务
	handler.RegisterHandlers(engine, sctx)

	engine.Run(fmt.Sprintf("%v:%v", c.Host, c.Port)) // listen and serve on 0.0.0.0:8080
}
