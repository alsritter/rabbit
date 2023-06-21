package server

import (
	helloworld_v1 "alsritter.icu/rabbit-template/api/helloworld/v1"
	"alsritter.icu/rabbit-template/internal/conf"
	"alsritter.icu/rabbit-template/internal/pkg/proc"
	"alsritter.icu/rabbit-template/internal/pkg/qy_middleware"
	"alsritter.icu/rabbit-template/internal/service/helloworld_service"

	prom "github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/gorilla/handlers"
)

// 需要记录响应的接口列表
var _responseLogList = map[string]struct{}{}

// NewHTTPServer new a HTTP server.
func NewHTTPServer(
	c *conf.Server,
	helloworldService *helloworld_service.HelloworldService,
	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		http.Middleware(
			metrics.Server(
				metrics.WithSeconds(prom.NewHistogram(_metricSeconds)),
				metrics.WithRequests(prom.NewCounter(_metricRequests)),
			),
			tracing.Server(),
			qy_middleware.Recovery(),
			validate.Validator(),
			qy_middleware.Trace(),
			qy_middleware.Server(logger, _responseLogList),
		),
		http.Filter(
			handlers.CORS(
				handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
				handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"}),
				handlers.AllowedOrigins([]string{"*"}),
			),
		),
	}

	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}

	srv := http.NewServer(opts...)

	// 注册程序内部接口
	proc.Register(srv)
	helloworld_v1.RegisterHelloServiceHTTPServer(srv, helloworldService)
	return srv
}
