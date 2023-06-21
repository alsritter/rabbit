// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"alsritter.icu/rabbit-template/internal/biz/helloworld"
	"alsritter.icu/rabbit-template/internal/conf"
	"alsritter.icu/rabbit-template/internal/data"
	"alsritter.icu/rabbit-template/internal/server"
	"alsritter.icu/rabbit-template/internal/service/helloworld_service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(confServer *conf.Server, confData *conf.Data, tracer *conf.Tracer, logger log.Logger) (*kratos.App, func(), error) {
	db := data.NewDB(confData, logger)
	pool := data.NewRedisConn(confData)
	dataData, cleanup, err := data.NewData(confData, logger, db, pool)
	if err != nil {
		return nil, nil, err
	}
	helloworldUsecase := biz.NewHelloworldUsecase(dataData, logger)
	helloworldService := helloworld_service.NewHelloworldService(helloworldUsecase)
	httpServer := server.NewHTTPServer(confServer, helloworldService, logger)
	cronServer := server.NewCronServer()
	app := newApp(logger, httpServer, cronServer)
	return app, func() {
		cleanup()
	}, nil
}
