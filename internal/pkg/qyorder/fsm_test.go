package qyorder

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	// 通过订单ID 或 其他信息查询到订单状态
	orderStatus := StatusPrepare
	orderMachine, err := NewFSM(orderStatus)
	if err != nil {
		t.Log(err.Error())
		return
	}

	// 关闭订单
	if _, err = orderMachine.Call(EventClose); err != nil {
		t.Log(err.Error())
	}

	// 创建订单，订单创建成功后再给用户发送短信
	if _, err = orderMachine.Call(EventInit,
		WithOrderId(1),
		WithOrderName("测试订单"),
	); err != nil {
		t.Log(err.Error())
	}

	// 支付订单
	if _, err = orderMachine.Call(EventPay,
		WithHandlerSendSMS(sendSMS)); err != nil {
		t.Log(err.Error())
	}

	// 退款订单
	if _, err = orderMachine.Call(EventRefunding); err != nil {
		t.Log(err.Error())
	}

	// 退款完成订单
	if _, err = orderMachine.Call(EventIsRefund); err != nil {
		t.Log(err.Error())
	}

	// 测试状态流转失败
	if _, err = orderMachine.Call(EventCancel); err != nil {
		t.Error(err.Error())
	}
}

// sendSMS 发送短信，为了演示写在了这里
func sendSMS(orderId int32, orderName string) error {
	fmt.Printf("订单 %s, 发送短信 \n", orderName)
	return nil
}
