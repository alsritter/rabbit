package qy_middleware

import (
	"context"
	"runtime"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"go.uber.org/zap"
)

// Option is recovery option.
type Option func(*options)

type options struct {
	handler recovery.HandlerFunc
}

// WithHandler with recovery handler.
func WithHandler(h recovery.HandlerFunc) Option {
	return func(o *options) {
		o.handler = h
	}
}

// WithLogger with recovery logger.
// Deprecated: use global logger instead.
func WithLogger(logger log.Logger) Option {
	return func(o *options) {}
}

// Recovery is a server middleware that recovers from any panics.
func Recovery(opts ...Option) middleware.Middleware {
	op := options{
		handler: func(ctx context.Context, req, err interface{}) error {
			return recovery.ErrUnknownRequest
		},
	}
	for _, o := range opts {
		o(&op)
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			defer func() {
				if rerr := recover(); rerr != nil {
					buf := make([]byte, 64<<10) //nolint:gomnd
					n := runtime.Stack(buf, false)
					buf = buf[:n]

					stackStr := strings.ReplaceAll(string(buf), "\t", "")
					stackArr := strings.Split(stackStr, "\n")
					// log.Context(ctx).Errorf("%v: %+v\n %v\n", rerr, req, stackArr)
					log.Context(ctx).Log(log.LevelError, "reason", rerr, "stack", zap.Strings("stack", stackArr))

					err = op.handler(ctx, req, rerr)
				}
			}()
			return handler(ctx, req)
		}
	}
}
