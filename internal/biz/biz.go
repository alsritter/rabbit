package biz

import (
	helloworld_biz "alsritter.icu/rabbit-template/internal/biz/helloworld"

	"github.com/google/wire"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	helloworld_biz.NewHelloworldUsecase,
)
