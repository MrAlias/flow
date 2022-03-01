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

package flow

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/trace"
)

type recorder struct {
	onStartN, onEndN, shutdownN, forceFlushN int
}

func (r *recorder) OnStart(context.Context, trace.ReadWriteSpan) {
	r.onStartN++
}

func (r *recorder) OnEnd(trace.ReadOnlySpan) {
	r.onEndN++
}

func (r *recorder) Shutdown(ctx context.Context) error {
	r.shutdownN++
	return nil
}

func (r *recorder) ForceFlush(ctx context.Context) error {
	r.forceFlushN++
	return nil
}

func TestDownstreamSpanProcessorCalled(t *testing.T) {
	r := new(recorder)
	sp := Wrap(r, WithListenAddress("localhost:0"))
	sp.OnStart(nil, nil)
	sp.OnEnd(nil)
	sp.ForceFlush(nil)
	sp.Shutdown(context.Background())

	assert.Equal(t, 1, r.onStartN, "wrong number of calls to OnStart")
	assert.Equal(t, 1, r.onEndN, "wrong number of calls to OnEnd")
	assert.Equal(t, 1, r.forceFlushN, "wrong number of calls to ForceFlush")
	assert.Equal(t, 1, r.shutdownN, "wrong number of calls to Shutdown")
}
