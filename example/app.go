// Copyright 2022 Tyler Yahn (MrAlias)
// Copyright The OpenTelemetry Authors
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

package main

import (
	"context"
	"math/rand"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// name is the Tracer name used to identify this instrumentation library.
const name = "fib"

// App is an Fibonacci computation application.
type App struct{}

// NewApp returns a new App.
func NewApp() *App {
	return &App{}
}

// Run starts computing random Fibonacci numbers.
func (a *App) Run(ctx context.Context) error {
	for {
		delay := time.Duration(rand.Intn(71)+30) * time.Millisecond
		select {
		case <-time.Tick(delay):
		case <-ctx.Done():
			return ctx.Err()
		}

		// Each execution of the run loop, we should get a new "root" span and context.
		newCtx, span := otel.Tracer(name).Start(ctx, "Run")
		a.Rand(newCtx)
		span.End()
	}
}

// Rand computes a random Fibonacci number between [0, 100].
func (a *App) Rand(ctx context.Context) {
	_, span := otel.Tracer(name).Start(ctx, "Rand")
	defer span.End()

	n := rand.Intn(101)
	span.SetAttributes(attribute.Int("request.n", n))
	_, err := Fibonacci(uint(n))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}
