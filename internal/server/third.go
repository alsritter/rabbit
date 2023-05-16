package server

import (
	"time"

	"alsritter.icu/rabbit-template/internal/conf"
	"alsritter.icu/rabbit-template/internal/pkg/httpclient"
	"alsritter.icu/rabbit-template/internal/pkg/poller"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func NewHttpClient(tracer trace.Tracer) *httpclient.Client {
	return httpclient.New(tracer)
}

func NewTracer(cf *conf.Tracer) trace.Tracer {
	return otel.Tracer(cf.ServiceNameKey)
}

func NewPoller(tracer trace.Tracer) *poller.Poller {
	return poller.NewPoller(time.Second, tracer)
}
