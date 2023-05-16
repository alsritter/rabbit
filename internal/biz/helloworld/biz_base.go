package biz

import (
	"alsritter.icu/rabbit-template/internal/data"

	"github.com/go-kratos/kratos/v2/log"
)

type base struct {
	log  *log.Helper
	data *data.Data
}
