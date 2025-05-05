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

package flow_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/MrAlias/flow"
	"go.opentelemetry.io/otel/sdk/trace"
)

func ExampleWithBatcher() {
	sdk := trace.NewTracerProvider(flow.WithBatcher(exporter{}))
	defer func() { _ = sdk.Shutdown(context.Background()) }()
}

func ExampleWithSpanProcessor() {
	spanProcessor := trace.NewSimpleSpanProcessor(exporter{})
	sdk := trace.NewTracerProvider(flow.WithSpanProcessor(spanProcessor))
	defer func() { _ = sdk.Shutdown(context.Background()) }()
}

func Example() {
	ctx := context.TODO()
	sdk := trace.NewTracerProvider(flow.WithBatcher(exporter{}))
	defer func() { _ = sdk.Shutdown(ctx) }()

	_, span := sdk.Tracer("flow-example").Start(ctx, "example")
	fmt.Println("started span")
	printSpansTotal()

	span.End()
	fmt.Println("ended span")
	printSpansTotal()

	// Output:
	// started span
	// spans_total{state="started"} 1
	// ended span
	// spans_total{state="ended"} 1
	// spans_total{state="started"} 1
	// exported: example
}

type exporter struct{}

func (e exporter) ExportSpans(_ context.Context, spans []trace.ReadOnlySpan) error {
	for _, span := range spans {
		fmt.Println("exported:", span.Name())
	}
	return nil
}

func (e exporter) Shutdown(ctx context.Context) error { return nil }

func printSpansTotal() {
	addr := fmt.Sprintf("http://localhost:%d/metrics", flow.DefaultListenPort)
	resp, err := http.Get(addr)
	if err != nil {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	for _, line := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(line, "spans_total") {
			fmt.Println(string(line))
		}
	}
}
