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
	"github.com/pulumi/pulumi-go-provider/infer"
)

var _ infer.Annotated = (*BuilderConfig)(nil)

// BuilderConfig configures the builder to use for an image build.
type BuilderConfig struct {
	Name string `pulumi:"name,optional"`
}

// Annotate sets docstrings on BuilderConfig.
func (b *BuilderConfig) Annotate(a infer.Annotator) {
	a.Describe(&b.Name, dedent(`
		Name of an existing buildx builder to use.

		Only "docker-container", "kubernetes", or "remote" drivers are
		supported. The legacy "docker" driver is not supported.

		Equivalent to Docker's "--builder" flag.
	`))
}
