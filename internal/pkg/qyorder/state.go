package qyorder

// 定义订单状态
const (
	StatusCreate        = State(1) // 无效
	StatusPrepare       = State(2) // 待支付
	StatusSuccess       = State(3) // 支付成功
	StatusCancel        = State(4) // 已取消
	StatusRefunding     = State(5) // 退款中
	StatusIsRefund      = State(6) // 已退款
	StatusClose         = State(7) // 已关闭
	StatusFail          = State(8) // 失败
	StatusRefundRefused = State(9) // 拒绝退款
)

// statusText 定义订单状态文案
var statusText = map[State]string{
	StatusCreate:        "无效",
	StatusPrepare:       "待支付",
	StatusSuccess:       "支付成功",
	StatusCancel:        "已取消",
	StatusRefunding:     "退款中",
	StatusIsRefund:      "已退款",
	StatusClose:         "已关闭",
	StatusFail:          "失败",
	StatusRefundRefused: "拒绝退款",
}

// statusEvent 定义订单状态对应的可操作事件
var statusEvent = map[State][]Event{
	StatusCreate:        {EventInit},
	StatusPrepare:       {EventPay, EventCancel, EventClose, EventIsFail},
	StatusSuccess:       {EventRefunding, EventIsFail, EventClose},
	StatusClose:         {EventRefunding},
	StatusRefundRefused: {EventIsFail, EventClose},
	StatusRefunding:     {EventIsRefund, EventRefundRefused, EventIsFail},
	StatusCancel:        {EventClose, EventIsFail},
}

func StatusText(status State) string {
	return statusText[status]
}
