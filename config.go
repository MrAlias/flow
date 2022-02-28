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

const (
	// DefaultListenPort is the port the HTTP server listens on if not
	// configured with the WithListenAddress option.
	DefaultListenPort = 41820
	// DefaultListenAddress is the listen address of the HTTP server if not
	// configured with the WithListenAddress option.
	DefaultListenAddress = ":41820"
)

type config struct {
	// address is the listen address for the HTTP server.
	address string
}

func newConfig(options []Option) config {
	c := config{
		address: DefaultListenAddress,
	}

	for _, opt := range options {
		c = opt.apply(c)
	}

	return c
}

type Option interface {
	apply(config) config
}

type addressOpt string

func (o addressOpt) apply(c config) config {
	c.address = string(o)
	return c
}

// WithListenAddresss sets the listen address of the HTTP server.
func WithListenAddress(addr string) Option {
	return addressOpt(addr)
}
