// Copyright 2024, Pulumi Corporation.
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

package internal

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/docker/buildx/builder"
	"github.com/docker/buildx/store/storeutil"
	"github.com/docker/cli/cli/command"
	cfgtypes "github.com/docker/cli/cli/config/types"
)

// host contains a host-level Docker CLI as well as a cache of initialized
// builders. Operations on the host are serialized.
type host struct {
	mu       sync.Mutex
	cli      command.Cli
	config   *Config
	builders map[string]*cachedBuilder
	auths    map[string]cfgtypes.AuthConfig
}

func newHost(config *Config) (*host, error) {
	docker, err := newDockerCLI(config)
	if err != nil {
		return nil, err
	}
	// Load existing credentials into memory.
	auths, err := docker.ConfigFile().GetAllCredentials()
	if err != nil {
		return nil, err
	}
	h := &host{
		cli:      docker,
		config:   config,
		builders: map[string]*cachedBuilder{},
		auths:    auths,
	}
	return h, err
}

// builderFor ensures a builder is available and running. This is guarded by a
// mutex to ensure other resources don't attempt to use the builder until it's
// ready.
//
// If the build doesn't specify a builder by name, we will iterate through all
// available builders until we find one that we can connect to.
func (h *host) builderFor(build Build) (*cachedBuilder, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	opts := build.BuildOptions()

	if b, ok := h.builders[opts.Builder]; ok {
		return b, nil
	}

	txn, release, err := storeutil.GetStore(h.cli)
	if err != nil {
		return nil, err
	}
	defer release()

	contextPathHash := opts.ContextPath
	if absContextPath, err := filepath.Abs(contextPathHash); err == nil {
		contextPathHash = absContextPath
	}
	b, err := builder.New(h.cli,
		builder.WithName(opts.Builder),
		builder.WithContextPathHash(contextPathHash),
		builder.WithStore(txn),
	)
	if err != nil {
		return nil, err
	}

	// If we didn't request a particular builder, and we loaded a default
	// builder with an unsupported (docker) driver, then look for a builder we
	// do support.
	if b.Driver == "" && opts.Builder == "" {
		builders, err := builder.GetBuilders(h.cli, txn)
		if err != nil {
			return nil, err
		}
	nextbuilder:
		for _, bb := range builders {
			if bb.Driver == "" {
				continue
			}
			if err := bb.Validate(); err != nil {
				continue
			}
			if bb.Err() != nil {
				continue
			}
			nodes, err := bb.LoadNodes(context.Background())
			if err != nil {
				continue
			}
			for _, n := range nodes {
				if n.Driver == nil {
					continue nextbuilder
				}
				if _, err := n.Driver.Dial(context.Background()); err != nil {
					continue nextbuilder
				}
				// TODO: Confirm the builder supports the requested platforms.
			}
			b = bb
			break
		}
	}

	if b.Driver == "" && opts.Builder == "" {
		// If we STILL don't have a builder, create a docker-container instance.
		b, err = builder.Create(context.Background(), txn, h.cli, builder.CreateOpts{Driver: "docker-container"})
		if err != nil {
			return nil, err
		}
	}

	// Attempt to load nodes in order to determine the builder's driver. Ignore
	// errors for "exec" builds because it's possible to request builders with
	// drivers that are unknown to us.
	nodes, err := b.LoadNodes(context.Background())
	if err != nil && !build.ShouldExec() {
		return nil, err
	}

	cached := &cachedBuilder{name: b.Name, driver: b.Driver, nodes: nodes}
	h.builders[opts.Builder] = cached

	return cached, nil
}

// cachedBuilder caches the builders we've loaded. Repeatedly fetching them can
// sometimes result in EOF errors from the daemon, especially when under load.
type cachedBuilder struct {
	name   string
	driver string
	nodes  []builder.Node
}
