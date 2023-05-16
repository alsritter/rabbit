package biz

import (
	"alsritter.icu/rabbit-template/internal/data"

	"github.com/go-kratos/kratos/v2/log"
)

type HelloworldUsecase struct {
	*base
}

func NewHelloworldUsecase(data *data.Data, logger log.Logger) *HelloworldUsecase {
	helloUsecase := &HelloworldUsecase{
		base: &base{
			log:  log.NewHelper(logger),
			data: data,
		},
	}
	return helloUsecase
}
