package events

const (
	TopicOrderPaid     = "order.paid"
	TopicShipAssembled = "ship.assembled"
)

type OrderPaid struct {
	EventUUID       string `json:"event_uuid"`
	OrderUUID       string `json:"order_uuid"`
	UserUUID        string `json:"user_uuid"`
	PaymentMethod   string `json:"payment_method"`
	TransactionUUID string `json:"transaction_uuid"`
}

type ShipAssembled struct {
	EventUUID    string `json:"event_uuid"`
	OrderUUID    string `json:"order_uuid"`
	UserUUID     string `json:"user_uuid"`
	BuildTimeSec int64  `json:"build_time_sec"`
}
