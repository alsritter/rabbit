package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"alsritter.icu/rabbit-template/internal/conf"
	"alsritter.icu/rabbit-template/internal/pkg/cron"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/encoding/json"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	configPath string

	id, _ = os.Hostname()
)

func init() {
	json.MarshalOptions.UseProtoNames = true
	json.MarshalOptions.EmitUnpopulated = true
	json.UnmarshalOptions.DiscardUnknown = false
	flag.StringVar(&configPath, "conf", "../../configs", "config path, eg: -conf config.yaml")
	flag.Parse()
	log.Infof("trying get CONFIG_FILE_NAME from os env...")
	configFileName := os.Getenv("CONFIG_FILE_NAME")
	if configFileName != "" {
		log.Infof("Import configFileName from os env: %v", configFileName)
		configPath = path.Join("./configs", configFileName)
	}
	log.Infof("config path: %v", configPath)
}

func newApp(logger log.Logger, hs *http.Server, cs *cron.Server) *kratos.App {
	app := kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			cs,
		),
	)

	return app
}

// 设置全局trace
func initTracer(url, serviceNameKey string) error {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return err
	}

	tp := tracesdk.NewTracerProvider(
		// 将基于父span的采样率设置为100%
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1.0))),
		// 始终确保在生产中批量处理
		tracesdk.WithBatcher(exp),
		// tracesdk.WithSpanProcessor(),
		// 在资源中记录有关此应用程序的信息
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(serviceNameKey),
			attribute.String("exporter", "jaeger"),
		)),
	)
	otel.SetTracerProvider(tp)
	return nil
}

func main() {
	logger := log.With(
		Logger(),
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)

	log.Infof("load config from: %v", configPath)
	c := config.New(
		config.WithSource(
			file.NewSource(configPath),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(fmt.Errorf("c.load err: %v", err))
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	err := initTracer(bc.Tracer.JaegerUrl, bc.Tracer.ServiceNameKey)
	if err != nil {
		panic(fmt.Errorf("initTracer err: %v", err))
	}

	app, cleanup, err := wireApp(bc.Server,
		bc.Data,
		bc.Tracer,
		logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
