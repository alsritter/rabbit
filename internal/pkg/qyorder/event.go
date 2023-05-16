package qyorder

// 定义订单事件
const (
	EventInit          = Event("初始化")
	EventPay           = Event("支付")
	EventCancel        = Event("取消")
	EventClose         = Event("关闭")
	EventRefunding     = Event("发起退款")
	EventIsRefund      = Event("已退款")
	EventIsFail        = Event("失败")
	EventRefundRefused = Event("拒绝退款")
)

// 定义订单事件对应的处理方法
var eventHandler = map[Event]Handler{
	EventInit:          handlerInit,
	EventCancel:        handlerCancel,
	EventClose:         handlerClose,
	EventPay:           handlerPay,
	EventRefunding:     handlerRefund,
	EventIsRefund:      handlerIsRefund,
	EventIsFail:        handlerFail,
	EventRefundRefused: handlerRefundRefused,
}
