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
	"fmt"
)

// keeper decides whether an element should be included for a preview
// operation, optionally returning a mutated copy of that element.
type keeper[T any] interface {
	keep(T) bool
}

// filter applies a keeper to each element, returning a new slice.
func filter[T any](k keeper[T], elems ...T) []T {
	if elems == nil {
		return nil
	}
	result := make([]T, 0, len(elems))
	for _, e := range elems {
		if !k.keep(e) {
			continue
		}
		result = append(result, e)
	}
	return result
}

// stringKeeper preserves any non-empty string values for preview.
type stringKeeper struct{ preview bool }

func (k stringKeeper) keep(s string) bool {
	if !k.preview {
		return true
	}
	return s != ""
}

//nolint:structcheck // False positive due to generics.
type stringerKeeper[T fmt.Stringer] struct{ preview bool }

//nolint:unused // False positive due to generics.
func (k stringerKeeper[T]) keep(t T) bool {
	if !k.preview {
		return true
	}
	return stringKeeper(k).keep(t.String())
}

// registryKeeper preserves any registries with known values for address and
// password. This is imprecise and doesn't permit alternative auth strategies
// like registry tokens, email, etc.
type registryKeeper struct{ preview bool }

//nolint:unused // False positive due to generics.
func (k registryKeeper) keep(r Registry) bool {
	if !k.preview {
		return true
	}
	return r.Password != "" && r.Address != ""
}

// mapKeeper preserves map elements with known keys and values.
type mapKeeper struct{ preview bool }

func (k mapKeeper) keep(m map[string]string) map[string]string {
	if !k.preview || len(m) == 0 {
		return m
	}
	kk := stringKeeper(k)
	filtered := make(map[string]string)
	for key, val := range m {
		if !kk.keep(key) {
			continue
		}
		if !kk.keep(val) {
			continue
		}
		filtered[key] = val
	}
	return filtered
}

type contextKeeper struct{ preview bool }

func (k contextKeeper) keep(bc *BuildContext) *BuildContext {
	if !k.preview || bc == nil || len(bc.Named) == 0 {
		return bc
	}

	named := NamedContexts{}
	sk := stringKeeper(k)
	for k, v := range bc.Named {
		if !sk.keep(k) || !sk.keep(v.Location) {
			continue
		}
		named[k] = v
	}

	return &BuildContext{
		Context: Context{bc.Location},
		Named:   named,
	}
}
