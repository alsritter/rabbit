package service

import (
	"alsritter.icu/rabbit-template/internal/service/helloworld_service"

	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(
	helloworld_service.NewHelloworldService,
)
