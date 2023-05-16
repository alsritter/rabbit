package qy_middleware

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// Reference:
// 对于响应加上链路追踪用的 trace id
// * https://opentelemetry.io/docs/instrumentation/go/getting-started/
// * https://juejin.cn/post/6968284700808839176
func Trace() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if ts, ok := transport.FromServerContext(ctx); ok {
				if ht, ok := ts.(http.Transporter); ok {
					ht.ReplyHeader().Set("trace.id", fmt.Sprintf("%s", tracing.TraceID()(ctx)))
				}
			}
			return handler(ctx, req)
		}
	}
}
