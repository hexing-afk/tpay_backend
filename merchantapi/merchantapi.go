package main

import (
	"flag"
	"tpay_backend/merchantapi/internal/common"
	"tpay_backend/utils"

	"github.com/tal-tech/go-zero/rest/httpx"

	"tpay_backend/merchantapi/internal/config"
	"tpay_backend/merchantapi/internal/handler"
	"tpay_backend/merchantapi/internal/svc"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest"
)

var configFile = flag.String("f", "merchantapi/etc/merchantapi-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 重置redis的数据库
	c.Redis.DB = utils.RedisDbMerchantapi

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
