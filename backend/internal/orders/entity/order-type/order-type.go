package order_type

type OrderType string

const (
	OrderTypeTakeaway OrderType = "takeaway"
	OrderTypeDineIn   OrderType = "dine_in"
	OrderTypeDelivery OrderType = "delivery"
)
