package qyorder

type Option func(*Opt)

// SendSMS 发送短信
type SendSMS func(orderId int32, orderName string) error

// Opt 定义 Handler 所需参数
type Opt struct {
	OrderId   int32
	OrderName string

	HandlerSendSMS SendSMS
}

// WithOrderId 设置订单ID
func WithOrderId(id int32) Option {
	return func(opt *Opt) {
		opt.OrderId = id
	}
}

// WithOrderName 设置订单名称
func WithOrderName(name string) Option {
	return func(opt *Opt) {
		opt.OrderName = name
	}
}

// WithHandlerSendSMS 设置发送短信
func WithHandlerSendSMS(sendSms SendSMS) Option {
	return func(opt *Opt) {
		opt.HandlerSendSMS = sendSms
	}
}
