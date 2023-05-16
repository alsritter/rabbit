package proc

import (
	"github.com/go-kratos/grpc-gateway/v2/protoc-gen-openapiv2/generator"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"net/http/pprof"
)

func Register(server *http.Server) {
	RegisterPprof(server)
	RegisterMetric(server)
	RegisterOpenAPI(server)
}

// 注册 pprof
func RegisterPprof(server *http.Server) {
	server.HandleFunc("/debug/pprof/", pprof.Index)
	server.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	server.HandleFunc("/debug/pprof/profile", pprof.Profile)
	server.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	server.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

// 注册监控体系指标
func RegisterMetric(server *http.Server) {
	server.Handle("/metrics", promhttp.Handler())
}

// 注册 openAPI
func RegisterOpenAPI(server *http.Server) {
	server.HandlePrefix("/q/", openapiv2.NewHandler(
		openapiv2.WithGeneratorOptions(generator.UseJSONNamesForFields(false), generator.EnumsAsInts(true)),
	))
}
