module github.com/MrAlias/flow/example

go 1.16

require (
	github.com/MrAlias/flow v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v1.4.1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.4.1
	go.opentelemetry.io/otel/sdk v1.4.1
)

replace github.com/MrAlias/flow => ../
