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
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	buildx "github.com/docker/buildx/build"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"

	"github.com/pulumi/pulumi-go-provider/infer"
)

// Dockerfile references a local, remote, or inline Dockerfile.
type Dockerfile struct {
	Location string `pulumi:"location,optional"`
	Inline   string `pulumi:"inline,optional"`
}

// Annotate sets docstrings on Dockerfile.
func (d *Dockerfile) Annotate(a infer.Annotator) {
	a.Describe(&d.Location, dedent(`
        Location of the Dockerfile to use.

        Can be a relative or absolute path to a local file, or a remote URL.

        Defaults to "${context.location}/Dockerfile" if context is on-disk.

        Conflicts with "inline".
    `))
	a.Describe(&d.Inline, dedent(`
        Raw Dockerfile contents.

        Conflicts with "location".

        Equivalent to invoking Docker with "-f -".
    `))
}

func (d *Dockerfile) validate(preview bool, c *Context) error {
	if d.Location != "" && d.Inline != "" {
		return newCheckFailure(
			errors.New(`only specify "file" or "inline", not both`),
			"dockerfile",
		)
	}

	if d.Location != "" {
		if buildx.IsRemoteURL(d.Location) {
			return nil
		}
		abs, err := filepath.Abs(d.Location)
		if err != nil {
			return err
		}
		f, err := os.Open(filepath.Clean(abs))
		if err != nil {
			return newCheckFailure(err, "dockerfile.location")
		}
		if err := parseDockerfile(f); err != nil {
			return newCheckFailure(err, "dockerfile.location")
		}
		return nil
	}

	if d.Inline != "" {
		err := parseDockerfile(strings.NewReader(d.Inline))
		if err != nil {
			return newCheckFailure(err, "dockerfile.inline")
		}
		return nil
	}

	if !preview && c != nil && !buildx.IsRemoteURL(c.Location) {
		return newCheckFailure(errors.New("missing 'location' or 'inline'"), "dockerfile")
	}

	return nil
}

func parseDockerfile(r io.Reader) error {
	parsed, err := parser.Parse(r)
	if err != nil {
		return newCheckFailure(err, "dockerfile")
	}
	_, _, err = instructions.Parse(parsed.AST)
	if err != nil {
		return err
	}
	return nil
}
