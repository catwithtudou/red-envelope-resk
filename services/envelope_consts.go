package services

const (
	DefaultBlessing = "恭喜发财"
)

//订单类型：发布单、退款单
type OrderType int

const (
	OrderTypeSending OrderType = iota + 1
	OrderTypeRefund
)

//支付状态：未支付，支付中，已支付，支付失败
//退款：未退款，退款中，已退款，退款失败
type PayStatus int

const (
	PayNothing PayStatus = iota + 1
	Paying
	Payed
	PayFailure
)

//红包订单状态：创建、发布、过期、失效
type OrderStatus int

const (
	OrderCreate   OrderStatus = iota + 1
	OrderSending
	OrderExpired
	OrderDisabled
)

//红包类型：普通红包，碰运气红包
type EnvelopeType int

const (
	GeneralEnvelopeType = iota + 1
	LuckyEnvelopeType
)
