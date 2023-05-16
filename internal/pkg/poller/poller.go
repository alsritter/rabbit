package poller

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Poller struct {
	interval time.Duration
	tracer   trace.Tracer
}

func NewPoller(defaultInterval time.Duration, tracer trace.Tracer) *Poller {
	return &Poller{
		interval: defaultInterval,
		tracer:   tracer,
	}
}

func (p *Poller) Poll(ctx context.Context, fn func(ctx context.Context) bool, interval ...time.Duration) (err error) {
	inv := p.interval
	if len(interval) > 0 {
		inv = interval[0]
	}

	ticker := time.NewTicker(inv)
	defer ticker.Stop()
	ctx, span := p.tracer.Start(ctx, "Poller.Poll")
	defer span.End()

	defer func() {
		if r := recover(); r != nil {
			span.SetAttributes(attribute.Bool("error", true))
			span.RecordError(r.(error))
			err = r.(error)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ctx, childSpan := p.tracer.Start(ctx, "Poller.PollLoop", trace.WithSpanKind(trace.SpanKindClient))
			defer childSpan.End()

			if fn(ctx) {
				return
			}
		}
	}
}

func (p *Poller) PollWithRetry(ctx context.Context, fn func(ctx context.Context) bool, retry int, interval ...time.Duration) (err error) {
	inv := p.interval
	if len(interval) > 0 {
		inv = interval[0]
	}

	ticker := time.NewTicker(inv)
	defer ticker.Stop()
	ctx, span := p.tracer.Start(ctx, "Poller.Poll")
	defer span.End()

	defer func() {
		if r := recover(); r != nil {
			span.SetAttributes(attribute.Bool("error", true))
			span.RecordError(r.(error))
			err = r.(error)
		}
	}()

	for i := 0; i < retry; i++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ctx, childSpan := p.tracer.Start(ctx, "Poller.PollLoop", trace.WithSpanKind(trace.SpanKindClient))
			defer childSpan.End()

			if fn(ctx) {
				return
			}
		}
	}

	return fmt.Errorf("Poller.PollWithRetry failed after %d retries", retry)
}
