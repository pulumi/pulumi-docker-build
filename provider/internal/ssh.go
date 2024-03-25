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
	"strings"

	"github.com/pulumi/pulumi-go-provider/infer"
)

// SSH is an SSH option.
type SSH struct {
	ID    string   `pulumi:"id"`
	Paths []string `pulumi:"paths,optional"`
}

// Annotate sets docstrings on SSH.
func (s *SSH) Annotate(a infer.Annotator) {
	a.Describe(&s.ID, dedent(`
		Useful for distinguishing different servers that are part of the same
		build.

		A value of "default" is appropriate if only dealing with a single host.
	`))
	a.Describe(&s.Paths, dedent(`
		SSH agent socket or private keys to expose to the build under the given
		identifier.

		Defaults to "[$SSH_AUTH_SOCK]".

		Note that your keys are **not** automatically added when using an
		agent. Run "ssh-add -l" locally to confirm which public keys are
		visible to the agent; these will be exposed to your build.
	`))
}

// String returns a CLI-encoded value for the SSH option, or an empty string if
// its ID is not known.
func (s SSH) String() string {
	if s.ID == "" {
		return ""
	}

	r := s.ID

	if len(s.Paths) > 0 {
		r += "=" + strings.Join(s.Paths, ",")
	}

	return r
}
