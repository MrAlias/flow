# flow

[![Go Reference](https://pkg.go.dev/badge/github.com/MrAlias/flow.svg)](https://pkg.go.dev/github.com/MrAlias/flow)

An OpenTelemetry [`SpanProcessor`] reporting tracing flow metrics.

## Getting Started

Assuming you have working code using the OpenTelemetry SDK, update the
registration of your exporter to use a wrapped [`SpanProcessor`].

Update your exporter registration with a [`BatchSpanProcessor`] to use the
equivalent `flow` [`TracerProviderOption`].

```go
import (
	"github.com/MrAlias/flow"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	sdk := trace.NewTracerProvider(flow.WithBatcher(exporter{}))
	/* ... */
}
```

More generically, all [`SpanProcessor`]s can be wrapped directly.

```go
import (
	"github.com/MrAlias/flow"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	spanProcessor := trace.NewSimpleSpanProcessor(exporter{})
	sdk := trace.NewTracerProvider(flow.WithSpanProcessor(spanProcessor))
	/* ... */
}
```

See the included [example](./example) for an end-to-end illustration of
functionality.

## Produced Metrics

The `flow` [`SpanProcessor`] will report `spans_total` metrics as a counter.
They are exposed at `localhost:41820` by default (this can be changed using the
`WithListenAddress` option).

```sh
$ curl -s http://localhost:41820/metrics | grep 'spans_total'

# HELP spans_total The total number of processed spans
# TYPE spans_total counter
spans_total{state="ended"} 762
spans_total{state="started"} 762
```

Configure a locally running [Prometheus] or [OpenTelemetry Collector] instance
to scrape these using a scrape target similar to this.

```yaml
scrape_configs:
- job_name: myapp
  static_configs:
  - targets:
    - 'localhost:41820'
```

[`SpanProcessor`]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#SpanProcessor
[`BatchSpanProcessor`]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#NewBatchSpanProcessor
[`TracerProviderOption`]: https://pkg.go.dev/go.opentelemetry.io/otel/sdk/trace#TracerProviderOption
[Prometheus]: https://prometheus.io/docs/prometheus/latest/configuration/configuration/
[OpenTelemetry Collector]: https://opentelemetry.io/docs/collector/configuration/#receivers
