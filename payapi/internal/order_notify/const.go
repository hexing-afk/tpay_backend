package order_notify

import "fmt"

const (
	//订单过期key
	PayNotifyExpireKey      = "payNotifyKey"
	TransferNotifyExpireKey = "transferNotifyKey"

	// 分布式锁
	ExpireKeyLock = "ExpireKeyLock"

	// 通知成功结果
	ResultNotifySuccess = "success"
)

// 代收订单过期key
func GetPayNotifyExpireKey(orderNo string) string {
	return fmt.Sprintf("%s_%s", PayNotifyExpireKey, orderNo)
}

// 代付订单过期key
func GetTransferNotifyExpireKey(orderNo string) string {
	return fmt.Sprintf("%s_%s", TransferNotifyExpireKey, orderNo)
}

// 分布式锁
func GetLockKey(orderNo string) string {
	return fmt.Sprintf("%s_%s", ExpireKeyLock, orderNo)
}
