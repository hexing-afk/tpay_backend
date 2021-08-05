package crontab

import (
	"context"
	"sync"
	"time"
	"tpay_backend/model"
	"tpay_backend/payapi/internal/order_notify"
	"tpay_backend/payapi/internal/svc"
	"tpay_backend/utils"

	"github.com/tal-tech/go-zero/core/logx"
)

type PayOrderNotify struct {
	CronBase
	serverCtx *svc.ServiceContext
}

func NewPayOrderNotifyTask(serverCtx *svc.ServiceContext) *PayOrderNotify {
	t := &PayOrderNotify{}
	t.LogCat = "代收订单-异步通知定时任务:"
	t.LockExpire = time.Minute * 5
	t.serverCtx = serverCtx
	return t
}

// 运行定时任务
func (l *PayOrderNotify) Run() {
	if l.Running { // 正在运行中
		return
	}

	lockKey := GetLockKey(l.serverCtx.Config.Name, l)
	lockValue := utils.NewUUID()

	// 获取分布式锁
	if !utils.GetDistributedLock(l.serverCtx.Redis, lockKey, lockValue, l.LockExpire) {
		return
	}
	l.Running = true

	l.queryNotifyOmissionOrder()
	l.queryNotifyBreakOrder()

	// 释放分布式锁
	utils.ReleaseDistributedLock(l.serverCtx.Redis, lockKey, lockValue)
	l.Running = false
}

// 查询通知遗漏的订单
func (l *PayOrderNotify) queryNotifyOmissionOrder() {
	//logx.Infof(l.LogCat + "检查通知遗漏的订单---------------------------------start")

	// 查询通知遗漏订单,checkTime=15s
	// 第一次通知后等15秒才会进行一下次通知, 如果pay_time < currentTime - 15s; 则第一次通知失败,要再通知一次
	checkTime := order_notify.GetNotifyWaitTimeById(order_notify.WaitTimesOne)
	payTime := time.Now().Add(time.Duration(-checkTime) * time.Second).Unix()

	orderStatus := []int64{model.PayOrderStatusPaid} //, model.PayOrderStatusTimeout
	orderNos, err := model.NewPayOrderModel(l.serverCtx.DbEngine).FindNotifyOmissionOrderNo(orderStatus, payTime)
	if err != nil {
		logx.Errorf(l.LogCat+"查询通知遗漏订单失败, err=[%v]", err)
		return
	}

	omissionOrderNum := len(orderNos)
	if omissionOrderNum <= 0 {
		return
	}
	logx.Infof(l.LogCat+"通知遗漏订单数量:%v, omissionOrders=%v", omissionOrderNum, orderNos)

	//并发处理
	l.parallelDoing(orderNos, omissionOrderNum)

	//logx.Infof(l.LogCat + "检查通知遗漏的订单------------------------- ------------end")
}

// 查询通知中断的订单
func (l *PayOrderNotify) queryNotifyBreakOrder() {
	//logx.Info(l.LogCat + "检查通知中断的订单---------------------------------start")

	//查询通知中断订单
	nextTime := time.Now().Unix()
	orderStatus := []int64{model.PayOrderStatusPaid} //, model.PayOrderStatusTimeout
	orderNos, err := model.NewPayOrderModel(l.serverCtx.DbEngine).FindNotifyBreakOrderNo(orderStatus, nextTime)
	if err != nil {
		logx.Errorf(l.LogCat+"查询通知中断的订单, err=[%v]", err)
		return
	}

	breakOrderNum := len(orderNos)
	if breakOrderNum <= 0 {
		return
	}
	logx.Infof(l.LogCat+"中断通知的订单数量:%v, breakOrders=%v", breakOrderNum, orderNos)

	//并发处理
	l.parallelDoing(orderNos, breakOrderNum)

	//logx.Info(l.LogCat + "检查通知中断的订单------------------------- ------------end")
}

func (l *PayOrderNotify) parallelDoing(orderNos []string, orderNum int) {
	parallelNum := 50 // 并发执行数量
	var wg sync.WaitGroup
	for i := 0; i < orderNum; i++ {
		if i > 0 && i%parallelNum == 0 {
			logx.Infof(l.LogCat+"并发等待----------------------------%v", i)
			wg.Wait()
		}
		wg.Add(1)
		go func(orderNo string) {
			defer wg.Done()
			l.insertRedis(orderNo)
		}(orderNos[i])
	}
	wg.Wait()
}

func (l *PayOrderNotify) insertRedis(orderNo string) {
	redisKey := order_notify.GetPayNotifyExpireKey(orderNo)
	err := l.serverCtx.Redis.SetNX(context.TODO(), redisKey, order_notify.WaitTimesOne, time.Second).Err()
	if err != nil {
		logx.Errorf(l.LogCat+"redis插入数据失败, key=%v, err=%v", redisKey, err)
	}
}
