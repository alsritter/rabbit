package data

import (
	"testing"

	"github.com/go-kratos/kratos/v2/log"
)

func TestNewDB(t *testing.T) {
	NewDB(nil, log.DefaultLogger)
}
