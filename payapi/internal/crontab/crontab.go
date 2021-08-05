package crontab

import (
	"fmt"
	"strings"
	"time"
	"tpay_backend/payapi/internal/svc"
	"tpay_backend/utils"

	"github.com/tal-tech/go-zero/core/logx"

	"github.com/robfig/cron/v3"
)

var c *cron.Cron

type CronBase struct {
	LogCat     string        // 用来记录日志类型
	Running    bool          // 任务是否正在运行中
	LockExpire time.Duration // 分布式锁: 有效期
}

// 获取分布式锁的redis键名(当前服务名称+即结构体名称)
func GetLockKey(prefix string, t interface{}) string {
	s := fmt.Sprintf("%sLock%T", prefix, t)
	s = strings.Replace(s, "*", "", -1) // 去掉符号 *
	s = strings.Replace(s, ".", "", -1) // 去掉符号 .
	return s
}

func Start(serverCtx *svc.ServiceContext) {
	// 添加定时任务
	c = cron.New(cron.WithLocation(utils.DefaultLocation), cron.WithSeconds())
	addCrontab(c, serverCtx)
	c.Start()

	logx.Info("定时任务已启动...")
}

func Stop() {
	if c != nil {
		c.Stop()
	}
}

// 添加定时任务到系统
func addCrontab(c *cron.Cron, serverCtx *svc.ServiceContext) {
	/*
	   Seconds      | Yes        | 0-59            | * / , -
	   Minutes      | Yes        | 0-59            | * / , -
	   Hours        | Yes        | 0-23            | * / , -
	   Day of month | Yes        | 1-31            | * / , - ?
	   Month        | Yes        | 1-12 or JAN-DEC | * / , -
	   Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?
	*/

	// 每天凌晨00:02分执行一次
	c.AddFunc("0 2 0 * * ?", NewPlatformCardReceivedClearTask(serverCtx).Run)

	// 每隔5分钟执行一次
	c.AddFunc("0 */3 * * * ?", NewPayOrderNotifyTask(serverCtx).Run)

	// 每隔5分钟执行一次
	c.AddFunc("0 */3 * * * ?", NewTransferOrderNotifyTask(serverCtx).Run)
}
