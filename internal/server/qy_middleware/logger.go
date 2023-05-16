package qy_middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

// Server is an server logging middleware.
func Server(logger log.Logger, responseLogList map[string]struct{}) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			var (
				code      int32
				reason    string
				kind      string
				operation string
			)
			startTime := time.Now()
			if info, ok := transport.FromServerContext(ctx); ok {
				kind = info.Kind().String()
				operation = info.Operation()
			}

			_ = log.WithContext(ctx, logger).Log(log.LevelInfo,
				"direction", "request",
				"component", kind,
				"operation", operation,
				"args", extractArgs(req),
			)

			reply, err = handler(ctx, req)
			if se := errors.FromError(err); se != nil {
				code = se.Code
				reason = se.Message
			}

			level, _ := extractError(err)
			// 如果是需要记录响应日志的接口，就记录响应日志
			if _, ok := responseLogList[operation]; ok {
				_ = log.WithContext(ctx, logger).Log(level,
					"direction", "response",
					"component", kind,
					"operation", operation,
					"args", extractArgs(req),
					"reply", extractArgs(reply),
					"code", code,
					"reason", reason,
					"latency", time.Since(startTime).Seconds(),
					// common.TraceIDKey, traceId,
				)
				return
			}

			_ = log.WithContext(ctx, logger).Log(level,
				"direction", "response",
				"component", kind,
				"operation", operation,
				"args", extractArgs(req),
				"code", code,
				"reason", reason,
				"latency", time.Since(startTime).Seconds(),
				// common.TraceIDKey, traceId,
			)
			return
		}
	}
}

// extractError returns the string of the error
func extractError(err error) (log.Level, string) {
	if err != nil {
		return log.LevelError, fmt.Sprintf("%+v", err)
	}
	return log.LevelInfo, ""
}

// extractArgs returns the string of the req
func extractArgs(req interface{}) string {
	if stringer, ok := req.(fmt.Stringer); ok {
		return stringer.String()
	}
	return fmt.Sprintf("%+v", req)
}
