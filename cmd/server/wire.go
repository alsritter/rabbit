//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"alsritter.icu/rabbit-template/internal/biz"
	"alsritter.icu/rabbit-template/internal/conf"
	"alsritter.icu/rabbit-template/internal/data"
	"alsritter.icu/rabbit-template/internal/server"
	"alsritter.icu/rabbit-template/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server,
	*conf.Data,
	*conf.Tracer,
	log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(
		data.ProviderSet,
		biz.ProviderSet,
		server.ProviderSet,
		service.ProviderSet,
		newApp))
}
