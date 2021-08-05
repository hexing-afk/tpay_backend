package order_notify

import (
	"context"
	"fmt"
	"strings"
	"time"
	"tpay_backend/utils"

	"tpay_backend/payapi/internal/svc"

	"github.com/tal-tech/go-zero/core/logx"
)

type ListenExpKeyHandler struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewListenExpKeyHandler(svcCtx *svc.ServiceContext) *ListenExpKeyHandler {
	return &ListenExpKeyHandler{
		svcCtx: svcCtx,
	}
}

// 开始监听过期的key
func (l *ListenExpKeyHandler) ListenRedisExpKey() {
	logx.Info("开始监听redis...")
	go l.subscribeExpireKey()
}

// 监听redis过期key事件
func (l *ListenExpKeyHandler) subscribeExpireKey() {
	pubSub := l.svcCtx.Redis.PSubscribe(context.TODO(), fmt.Sprintf("__keyevent@%v__:expired", utils.RedisDbPayapi))
	defer pubSub.Close()

	for {
		msg, err := pubSub.ReceiveMessage(context.TODO())
		if err != nil {
			logx.Errorf("获取订阅过期key出错,err:%v", err)
			time.Sleep(time.Millisecond * 500) // 出错时稍微等待
			continue
		}

		l.handleExpireKeyEvent(msg.Channel, msg.Pattern, msg.Payload)
	}
}

// 处理过期key事件
func (l *ListenExpKeyHandler) handleExpireKeyEvent(channel, pattern, payload string) {
	logx.Infof("ListenRedisExpKey()获的一条过期数据-----payload:%s", payload)

	expireKeys := strings.Split(payload, "_")
	if len(expireKeys) != 2 {
		return
	}

	key := expireKeys[0]
	if key == "" {
		return
	}

	orderNo := expireKeys[1]
	if orderNo == "" {
		return
	}

	var notifyObj Notify
	switch key {
	case PayNotifyExpireKey:
		notifyObj = NewPayOrderNotify(context.TODO(), l.svcCtx)
	case TransferNotifyExpireKey:
		notifyObj = NewTransferOrderNotify(context.TODO(), l.svcCtx)
	default:
		return
	}

	lockKey := GetLockKey(orderNo)
	lockValue := utils.NewUUID()

	// 获取分布式锁
	if !utils.GetDistributedLock(l.svcCtx.Redis, lockKey, lockValue, 30*time.Second) {
		return
	}

	notifyObj.OrderNotify(orderNo)

	// 释放分布式锁
	utils.ReleaseDistributedLock(l.svcCtx.Redis, lockKey, lockValue)
}
