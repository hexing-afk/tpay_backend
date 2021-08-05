package order_notify

type Notify interface {
	// 订单异步通知
	OrderNotify(orderNo string)
}
