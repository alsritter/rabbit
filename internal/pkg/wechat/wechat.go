package wechat

import (
	"errors"

	"alsritter.icu/rabbit-template/api/common"

	"github.com/ArtisanCloud/PowerWeChat/v2/src/miniProgram"
	"github.com/ArtisanCloud/PowerWeChat/v2/src/payment"
)

var (
	ErrInstanceNotSupport = errors.New("当前小程序未配置 SDK")
)

type MiniProgram struct {
	MarketApp   *miniProgram.MiniProgram
	MerchantApp *miniProgram.MiniProgram
	AssetsApp   *miniProgram.MiniProgram
	ManagerApp  *miniProgram.MiniProgram
}

func (mp *MiniProgram) GetMiniProgram(sourceCode common.SourceType) (*miniProgram.MiniProgram, error) {
	switch sourceCode {
	case common.SourceType_SOURCE_TYPE_MARKET:
		return mp.MarketApp, nil
	case common.SourceType_SOURCE_TYPE_MERCHANT:
		return mp.MerchantApp, nil
	case common.SourceType_SOURCE_TYPE_MANAGER:
		return mp.ManagerApp, nil
	case common.SourceType_SOURCE_TYPE_ASSETS:
		return mp.AssetsApp, nil
	}
	return nil, ErrInstanceNotSupport
}

// NewMiniProgram 初始化小程序实例
func NewMiniProgram(appId, secret string) (*miniProgram.MiniProgram, error) {
	application, err := miniProgram.NewMiniProgram(&miniProgram.UserConfig{
		AppID:  appId,
		Secret: secret,
	})
	if err != nil {
		return nil, err
	}
	return application, nil
}

type MiniPayment struct {
	MarketApp   *payment.Payment
	MerchantApp *payment.Payment
	AssetsApp   *payment.Payment
	ManagerApp  *payment.Payment
}

func (mp *MiniPayment) GetMiniPayment(sourceCode common.SourceType) (*payment.Payment, error) {
	switch sourceCode {
	case common.SourceType_SOURCE_TYPE_MARKET:
		return mp.MarketApp, nil
	case common.SourceType_SOURCE_TYPE_MERCHANT:
		return mp.MerchantApp, nil
	case common.SourceType_SOURCE_TYPE_MANAGER:
		return mp.ManagerApp, nil
	case common.SourceType_SOURCE_TYPE_ASSETS:
		return mp.AssetsApp, nil
	}
	return nil, ErrInstanceNotSupport
}

// NewPaymentApp 初始化支付实例
func NewPaymentApp(conf *payment.UserConfig) (*payment.Payment, error) {
	Payment, err := payment.NewPayment(conf)
	if err != nil {
		return nil, err
	}
	return Payment, nil
}
