package crontab

import (
	"log"
	"time"
	"tpay_backend/model"
	"tpay_backend/payapi/internal/svc"
	"tpay_backend/utils"
)

// 平台收款卡今日已收清零

type PlatformCardReceivedClear struct {
	CronBase
	serverCtx *svc.ServiceContext
}

func NewPlatformCardReceivedClearTask(serverCtx *svc.ServiceContext) *PlatformCardReceivedClear {
	t := &PlatformCardReceivedClear{}
	t.LogCat = "平台收款卡-今日已收金额清零定时任务:"
	t.LockExpire = time.Minute * 5
	t.serverCtx = serverCtx
	return t
}

// 运行定时任务
func (t *PlatformCardReceivedClear) Run() {
	if t.Running { // 正在运行中
		return
	}

	lockKey := GetLockKey(t.serverCtx.Config.Name, t)
	lockValue := utils.NewUUID()

	// 获取分布式锁
	if !utils.GetDistributedLock(t.serverCtx.Redis, lockKey, lockValue, t.LockExpire) {
		return
	}
	t.Running = true

	t.Doing()

	// 释放分布式锁
	utils.ReleaseDistributedLock(t.serverCtx.Redis, lockKey, lockValue)
	t.Running = false
}

func (t *PlatformCardReceivedClear) Doing() {

	log.Println(t.LogCat + "------>PlatformCardReceivedClear.Doing...Start")

	if err := model.NewPlatformBankCardModel(t.serverCtx.DbEngine).TodayReceivedClear(); err != nil {
		log.Printf(t.LogCat+"平台收款卡今日已收金额清零失败, err=%v", err)
	}

	log.Println(t.LogCat + "------>PlatformCardReceivedClear.Doing...End")
}
