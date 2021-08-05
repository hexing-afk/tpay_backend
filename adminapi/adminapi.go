package main

import (
	"flag"
	"tpay_backend/adminapi/internal/common"
	"tpay_backend/utils"

	"github.com/tal-tech/go-zero/core/logx"

	"github.com/tal-tech/go-zero/rest/httpx"

	"tpay_backend/adminapi/internal/config"
	"tpay_backend/adminapi/internal/handler"
	"tpay_backend/adminapi/internal/svc"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/rest"
)

var configFile = flag.String("f", "adminapi/etc/adminapi-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 重置redis的数据库
	c.Redis.DB = utils.RedisDbAdminapi

	// 设置时区
	utils.SetTimezone(c.Timezone)

	ctx := svc.NewServiceContext(c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 设置错误处理函数
	httpx.SetErrorHandler(common.ErrorHandler)

	handler.RegisterHandlers(server, ctx)

	logx.Infof("Starting server at %s:%d...", c.Host, c.Port)
	server.Start()
}
