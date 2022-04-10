module github.com/MrAlias/flow/example

go 1.16

require (
	github.com/MrAlias/flow v0.1.1
	go.opentelemetry.io/otel v1.6.3
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.6.1
	go.opentelemetry.io/otel/sdk v1.6.3
)

replace github.com/MrAlias/flow => ../
