package orderStatus

type OrderStatus string

const (
	OrderStatusNew            OrderStatus = "NEW"
	OrderStatusWaitingPayment OrderStatus = "WAITING_PAYMENT"
	OrderStatusCooking        OrderStatus = "COOKING"
	OrderStatusReady          OrderStatus = "READY"
	OrderStatusCompleted      OrderStatus = "COMPLETED"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)
