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

import "github.com/pulumi/pulumi-go-provider/infer"

// Registry contains credentials for authenticating with a remote registry.
type Registry struct {
	Address  string `pulumi:"address"`
	Password string `pulumi:"password,optional" provider:"secret"`
	Username string `pulumi:"username,optional"`
}

// Annotate sets docstrings on Registry.
func (r *Registry) Annotate(a infer.Annotator) {
	a.Describe(&r.Address, `The registry's address (e.g. "docker.io").`)
	a.Describe(&r.Username, `Username for the registry.`)
	a.Describe(&r.Password, `Password or token for the registry.`)
}
