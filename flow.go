// Copyright 2022 Tyler Yahn (MrAlias)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package flow provides an OpenTelemetry SpanProcessor that reports telemetry
// flow as Prometheus metrics.
//
// To start using, replace the TracerProviderOption from the default
// OpenTelemetry SDK with the ones provided here. For example:
//   sdk := trace.NewTracerProvider(trace.WithBatcher(exporter{}))
//
// Can be replaced with:
//   sdk := trace.NewTracerProvider(flow.WithBatcher(exporter{}))
//
// Additionally, any custom span processor can be wrapped into a
// TracerProviderOption. For example:
//   spanProcessor := trace.NewSimpleSpanProcessor(exporter{})
//   sdk := trace.NewTracerProvider(flow.WithSpanProcessor(spanProcessor))
package flow

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
)

const (
	startedState = "started"
	endedState   = "ended"
)

type spanProcessor struct {
	wrapped trace.SpanProcessor

	idleConnsClosed chan struct{}
	server          *http.Server
	spanCounter     *prometheus.CounterVec
}

// Wrap returns a wrapped version of the downstream SpanProcessor with
// telemetry flow reporting. All calls to the returned SpanProcessor will
// introspected for telemetry data and then forwarded to downstream.
func Wrap(downstream trace.SpanProcessor, options ...Option) trace.SpanProcessor {
	mux := http.NewServeMux()
	registry := prometheus.NewRegistry()
	mux.Handle("/metrics", promhttp.InstrumentMetricHandler(
		registry,
		promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
	))

	c := newConfig(options)
	sp := &spanProcessor{
		wrapped:         downstream,
		idleConnsClosed: make(chan struct{}),
		server:          &http.Server{Addr: c.address, Handler: mux},
		spanCounter: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "spans_total",
			Help: "The total number of processed spans",
		}, []string{"state"}),
	}
	registry.MustRegister(sp.spanCounter)

	go func() {
		switch err := sp.server.ListenAndServe(); err {
		case nil, http.ErrServerClosed:
		default:
			otel.Handle(err)
		}
		close(sp.idleConnsClosed)
	}()

	return sp
}

// WithSpanProcessor returns an option that registers spanProcessor with a
// TracerProvider after wrapping it to report telemetry flow metrics.
func WithSpanProcessor(spanProcessor trace.SpanProcessor, options ...Option) trace.TracerProviderOption {
	return trace.WithSpanProcessor(Wrap(spanProcessor, options...))
}

// WithBatcher returns an option that registers exporter using a
// BatchSpanProcessor with a TracerProvider after wrapping it to report
// telemetry flow metrics.
//
// If configuration of the flow span processor is needed, use
// WithSpanProcessor or Wrap directly.
func WithBatcher(exporter trace.SpanExporter, options ...trace.BatchSpanProcessorOption) trace.TracerProviderOption {
	spanProcessor := trace.NewBatchSpanProcessor(exporter, options...)
	return trace.WithSpanProcessor(Wrap(spanProcessor))
}

// OnStart is called when a span is started.
func (sp *spanProcessor) OnStart(parent context.Context, s trace.ReadWriteSpan) {
	sp.spanCounter.WithLabelValues(startedState).Inc()
	sp.wrapped.OnStart(parent, s)
}

// OnEnd is called when span is finished.
func (sp *spanProcessor) OnEnd(s trace.ReadOnlySpan) {
	sp.spanCounter.WithLabelValues(endedState).Inc()
	sp.wrapped.OnEnd(s)
}

// Shutdown is called when the SDK shuts down. The telemetry reporting process
// will be halted when this is called.
func (sp *spanProcessor) Shutdown(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- sp.wrapped.Shutdown(ctx)
	}()

	err := sp.server.Shutdown(ctx)
	select {
	case <-ctx.Done():
		// Abandon idle conns if context has expired.
		if err == nil {
			return ctx.Err()
		}
		return err
	case <-sp.idleConnsClosed:
	}

	// Downstream honors ctx timeout, no need to include in select above.
	if e := <-errCh; e != nil {
		// Prioritize downstream error over server shutdown error.
		err = e
	}
	return err
}

// ForceFlush forwards call to wrapped SpanProcessor.
func (sp *spanProcessor) ForceFlush(ctx context.Context) error {
	return sp.wrapped.ForceFlush(ctx)
}
