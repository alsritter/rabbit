package qyorder

var (
	// handlerInit 初始化订单
	handlerInit = Handler(func(opt *Opt) (State, error) {
		return StatusPrepare, nil
	})

	// handlerCancel 关闭订单
	handlerCancel = Handler(func(opt *Opt) (State, error) {
		return StatusCancel, nil
	})

	// handlerClose 关闭订单
	handlerClose = Handler(func(opt *Opt) (State, error) {
		return StatusClose, nil
	})

	// handlerPay 支付订单
	handlerPay = Handler(func(opt *Opt) (State, error) {
		if opt.HandlerSendSMS != nil {
			_ = opt.HandlerSendSMS(opt.OrderId, opt.OrderName)
		}
		return StatusSuccess, nil
	})

	// handlerRefund 退款订单
	handlerRefund = Handler(func(opt *Opt) (State, error) {
		return StatusRefunding, nil
	})

	// handlerRefund 退款完成订单
	handlerIsRefund = Handler(func(opt *Opt) (State, error) {
		return StatusIsRefund, nil
	})

	// handlerFail 订单失败的事件
	handlerFail = Handler(func(opt *Opt) (State, error) {
		return StatusFail, nil
	})

	// handlerRefundRefused
	handlerRefundRefused = Handler(func(opt *Opt) (State, error) {
		return StatusRefundRefused, nil
	})
)
