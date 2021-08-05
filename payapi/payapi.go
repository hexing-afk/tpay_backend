package main

import (
	"flag"
	notify "tpay_backend/payapi/internal/order_notify"

	"tpay_backend/payapi/internal/common"
	"tpay_backend/payapi/internal/config"
	"tpay_backend/payapi/internal/crontab"
	"tpay_backend/payapi/internal/handler"
	"tpay_backend/payapi/internal/svc"
	"tpay_backend/utils"

	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest"
	"github.com/tal-tech/go-zero/rest/httpx"
)

var configFile = flag.String("f", "payapi/etc/payapi-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	// 加载配置
	conf.MustLoad(*configFile, &c)

	// 重置redis的数据库
	c.Redis.DB = utils.RedisDbPayapi

	// 设置时区
	utils.SetTimezone(c.Timezone)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	// 设置错误处理函数
	httpx.SetErrorHandler(common.ErrorHandler)

	// 注册服务
	handler.RegisterHandlers(server, ctx)

	// 启动定时任务
	crontab.Start(ctx)
	defer crontab.Stop()

	// redis过期监听
	notify.NewListenExpKeyHandler(ctx).ListenRedisExpKey()

	logx.Info("test-2020-04-24 16:12...")

	logx.Infof("Starting server at %s:%d...", c.Host, c.Port)
	server.Start()
}
