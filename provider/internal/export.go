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
	"fmt"
	"os"
	"slices"
	"strings"

	controllerapi "github.com/docker/buildx/controller/pb"
	"github.com/moby/buildkit/client"
	"github.com/tonistiigi/go-csvvalue"

	"github.com/pulumi/pulumi-go-provider/infer"
)

var (
	_ fmt.Stringer    = (*Export)(nil)
	_ fmt.Stringer    = (*ExportDocker)(nil)
	_ fmt.Stringer    = (*ExportImage)(nil)
	_ fmt.Stringer    = (*ExportLocal)(nil)
	_ fmt.Stringer    = (*ExportOCI)(nil)
	_ fmt.Stringer    = (*ExportRegistry)(nil)
	_ fmt.Stringer    = (*ExportTar)(nil)
	_ fmt.Stringer    = ExportWithAnnotations{}
	_ fmt.Stringer    = ExportWithCompression{}
	_ fmt.Stringer    = ExportWithNames{}
	_ fmt.Stringer    = ExportWithOCI{}
	_ infer.Annotated = (*Export)(nil)
	_ infer.Annotated = (*ExportDocker)(nil)
	_ infer.Annotated = (*ExportImage)(nil)
	_ infer.Annotated = (*ExportLocal)(nil)
	_ infer.Annotated = (*ExportOCI)(nil)
	_ infer.Annotated = (*ExportRegistry)(nil)
	_ infer.Annotated = (*ExportTar)(nil)
)

// Export is a "union" type for all of our available `--output` options.
type Export struct {
	Tar       *ExportTar       `pulumi:"tar,optional"`
	Local     *ExportLocal     `pulumi:"local,optional"`
	Registry  *ExportRegistry  `pulumi:"registry,optional"`
	Image     *ExportImage     `pulumi:"image,optional"`
	OCI       *ExportOCI       `pulumi:"oci,optional"`
	Docker    *ExportDocker    `pulumi:"docker,optional"`
	CacheOnly *ExportCacheOnly `pulumi:"cacheonly,optional"`
	Raw       Raw              `pulumi:"raw,optional"`

	Disabled bool `pulumi:"disabled,optional"`
}

// Annotate sets docstrings on Export.
func (e *Export) Annotate(a infer.Annotator) {
	a.Describe(&e.Tar, dedent(`
		Export to a local directory as a tarball.`,
	))
	a.Describe(&e.Local, dedent(`
		Export to a local directory as files and directories.`,
	))
	a.Describe(&e.Registry, dedent(`
		Identical to the Image exporter, but pushes by default.`,
	))
	a.Describe(&e.Image, dedent(`
		Outputs the build result into a container image format.`,
	))
	a.Describe(&e.OCI, dedent(`
		Identical to the Docker exporter but uses OCI media types by default.`,
	))
	a.Describe(&e.Docker, dedent(`
		Export as a Docker image layout.`,
	))
	a.Describe(&e.Raw, dedent(`
		A raw string as you would provide it to the Docker CLI (e.g.,
		"type=docker")`,
	))
	a.Describe(&e.CacheOnly, dedent(`
		A no-op export. Helpful for silencing the 'no exports' warning if you
		just want to populate caches.
	`))

	a.Describe(&e.Disabled, dedent(`
		When "true" this entry will be excluded. Defaults to "false".
	`))
}

// String returns a CLI-encoded value for this `--output` entry, or an empty
// string if disabled. `validate` should be called to ensure only one entry was
// set.
func (e Export) String() string {
	if e.Disabled {
		return ""
	}
	return join(e.Tar, e.Local, e.Registry, e.Image, e.OCI, e.Docker, e.CacheOnly, e.Raw)
}

// pushed returns true if the export would result in a registry push.
func (e Export) pushed() bool {
	if e.Raw != "" {
		exp, err := parseExports([]string{e.Raw.String()})
		if err != nil {
			return false
		}
		return exp[0].Attrs["push"] == "true"
	}
	if e.Registry != nil {
		return e.Registry.Push == nil || *e.Registry.Push
	}
	if e.Image != nil {
		return e.Image.Push != nil && *e.Image.Push
	}
	return false
}

// parseExports is forked from docker/buildx@v0.18.0 from util/buildflags/export.go
// to maintain the old logic. This is to get a working version of the provider with
// the latest buildx while maintaining the old behaviour.
//
// TODO: Remove this fork and update existing logic/tests.
func parseExports(inp []string) ([]*controllerapi.ExportEntry, error) {
	if len(inp) == 0 {
		return nil, nil
	}
	outs := make([]*controllerapi.ExportEntry, 0, len(inp))

	for _, s := range inp {
		fields, err := csvvalue.Fields(s, nil)
		if err != nil {
			return nil, err
		}

		out := controllerapi.ExportEntry{
			Attrs: map[string]string{},
		}
		if len(fields) == 1 && fields[0] == s && !strings.HasPrefix(s, "type=") {
			if s != "-" {
				outs = append(outs, &controllerapi.ExportEntry{
					Type:        client.ExporterLocal,
					Destination: s,
				})
				continue
			}
			out = controllerapi.ExportEntry{
				Type:        client.ExporterTar,
				Destination: s,
			}
		}

		if out.Type == "" {
			for _, field := range fields {
				parts := strings.SplitN(field, "=", 2)
				if len(parts) != 2 {
					return nil, fmt.Errorf("invalid value %s", field)
				}
				key := strings.TrimSpace(strings.ToLower(parts[0]))
				value := parts[1]
				switch key {
				case "type":
					out.Type = value
				default:
					out.Attrs[key] = value
				}
			}
		}
		if out.Type == "" {
			return nil, errors.New("type is required for output")
		}

		if out.Type == "registry" {
			out.Type = client.ExporterImage
			if _, ok := out.Attrs["push"]; !ok {
				out.Attrs["push"] = "true"
			}
		}

		if dest, ok := out.Attrs["dest"]; ok {
			out.Destination = dest
			delete(out.Attrs, "dest")
		}

		outs = append(outs, &out)
	}
	return outs, nil
}

func (e Export) validate(preview bool, tags []string) (*controllerapi.ExportEntry, error) {
	if strings.Count(e.String(), "type=") > 1 {
		return nil, errors.New("exports should only specify one export type")
	}
	ee, err := parseExports([]string{e.String()})
	if err != nil {
		return nil, err
	}

	exp := ee[0]
	if len(tags) == 0 && isRegistryPush(exp) && exp.Attrs["name"] == "" {
		return nil, errors.New(
			"at least one tag or export name is needed when pushing to a registry",
		)
	}
	if !preview {
		return exp, nil
	}

	// Don't perform registry pushes during previews.
	if exp.Type == "image" {
		exp.Attrs["push"] = "false"
	}
	return exp, nil
}

// ExportCacheOnly is a dummy/no-op --cache-to entry. It exists only to help
// silence the "no exports configured" warning. By using this the user signals
// that they intentionally do not want exports, and only caches will be
// populated as a result.
type ExportCacheOnly struct{}

// String returns the CLI-encoded value of these export options, or an empty
// string if the receiver is nil.
func (e *ExportCacheOnly) String() string {
	if e == nil {
		return ""
	}
	return "type=cacheonly"
}

// ExportDocker pushes the final image to the local build daemon.
type ExportDocker struct {
	ExportWithOCI
	ExportWithCompression
	ExportWithAnnotations
	ExportWithNames

	Dest string `pulumi:"dest,optional"`
	Tar  *bool  `pulumi:"tar,optional"`
}

// Annotate sets docstrings and defaults on ExportDocker.
func (e *ExportDocker) Annotate(a infer.Annotator) {
	a.SetDefault(&e.Tar, true)

	a.Describe(&e.Dest, "The local export path.")
	a.Describe(&e.Tar, "Bundle the output into a tarball layout.")
}

// String returns the CLI-encoded value of these export options, or an empty
// string if the receiver is nil.
func (e *ExportDocker) String() string {
	if e == nil {
		return ""
	}
	parts := []string{}
	if e.Dest != "" {
		parts = append(parts, "dest="+e.Dest)
	}
	if e.Tar != nil {
		parts = append(parts, fmt.Sprintf("tar=%t", *e.Tar))
	}

	return join(
		Raw("type=docker"),
		Raw(strings.Join(parts, ",")),
		e.ExportWithOCI,
		e.ExportWithCompression,
		e.ExportWithAnnotations,
		e.ExportWithNames,
	)
}

// ExportOCI is a cache that defaults to using OCI media types.
type ExportOCI struct {
	ExportDocker
}

// Annotate sets docstrings and defaults on ExportOCI.
func (e *ExportOCI) Annotate(a infer.Annotator) {
	a.SetDefault(&e.OCI, true)
	a.Describe(&e.OCI, "Use OCI media types in exporter manifests.")
}

func (e *ExportOCI) String() string {
	if e == nil {
		return ""
	}
	return strings.Replace(e.ExportDocker.String(), "type=docker", "type=oci", 1)
}

// ExportImage can push the final image to remote registries.
type ExportImage struct {
	ExportWithOCI
	ExportWithCompression
	ExportWithNames
	ExportWithAnnotations

	Push               *bool  `pulumi:"push,optional"`
	PushByDigest       *bool  `pulumi:"pushByDigest,optional"`
	Insecure           *bool  `pulumi:"insecure,optional"`
	DanglingNamePrefix string `pulumi:"danglingNamePrefix,optional"`
	NameCanonical      *bool  `pulumi:"nameCanonical,optional"`
	Unpack             *bool  `pulumi:"unpack,optional"`
	Store              *bool  `pulumi:"store,optional"`
}

// Annotate sets docstrings and defaults on ExportImage.
func (e *ExportImage) Annotate(a infer.Annotator) {
	a.SetDefault(&e.Store, true)

	a.Describe(&e.Store, dedent(`
		Store resulting images to the worker's image store and ensure all of
		its blobs are in the content store.

		Defaults to "true".

		Ignored if the worker doesn't have image store (when using OCI workers,
		for example).
	`))
	a.Describe(&e.Push, dedent(`
		Push after creating the image. Defaults to "false".
	`))
	a.Describe(&e.DanglingNamePrefix, dedent(`
		Name image with "prefix@<digest>", used for anonymous images.
	`))
	a.Describe(&e.NameCanonical, dedent(`
		Add additional canonical name ("name@<digest>").
	`))
	a.Describe(&e.Insecure, dedent(`
		Allow pushing to an insecure registry.
	`))
	a.Describe(&e.PushByDigest, dedent(`
		Push image without name.
	`))
	a.Describe(&e.Unpack, dedent(`
		Unpack image after creation (for use with containerd). Defaults to
		"false".
	`))
}

// String returns the CLI-encoded value of these export options, or an empty
// string if the receiver is nil.
func (e *ExportImage) String() string {
	if e == nil {
		return ""
	}
	parts := []string{}
	if e.Push != nil {
		parts = append(parts, fmt.Sprintf("push=%t", *e.Push))
	}
	if e.PushByDigest != nil {
		parts = append(parts, fmt.Sprintf("push-by-digest=%t", *e.PushByDigest))
	}
	if e.Insecure != nil {
		parts = append(parts, fmt.Sprintf("insecure=%t", *e.Insecure))
	}
	if e.DanglingNamePrefix != "" {
		parts = append(parts, "dangling-name-prefix="+e.DanglingNamePrefix)
	}
	if e.NameCanonical != nil {
		parts = append(parts, fmt.Sprintf("name-canonical=%t", *e.NameCanonical))
	}
	if e.Unpack != nil {
		parts = append(parts, fmt.Sprintf("unpack=%t", *e.Unpack))
	}
	if e.Store != nil {
		parts = append(parts, fmt.Sprintf("store=%t", *e.Store))
	}
	return join(
		Raw("type=image"),
		Raw(strings.Join(parts, ",")),
		e.ExportWithOCI,
		e.ExportWithCompression,
		e.ExportWithNames,
		e.ExportWithAnnotations,
	)
}

// ExportRegistry is equivalent to ExportImage but defaults to push=true.
type ExportRegistry struct {
	ExportImage
}

// Annotate sets docstrings and defaults on ExportRegistry.
func (e *ExportRegistry) Annotate(a infer.Annotator) {
	a.Describe(&e.Push, dedent(`
		Push after creating the image. Defaults to "true".
	`))
	a.SetDefault(&e.Push, true)
}

// String returns the CLI-encoded value of these export options, or an empty
// string if the receiver is nil.
func (e *ExportRegistry) String() string {
	if e == nil {
		return ""
	}
	return strings.Replace(e.ExportImage.String(), "type=image", "type=registry", 1)
}

// ExportLocal writes the final image to disk.
type ExportLocal struct {
	Dest string `pulumi:"dest"`
}

// String returns the CLI-encoded value of these export options, or an empty
// string if the receiver is nil.
func (e *ExportLocal) String() string {
	if e == nil {
		return ""
	}
	return "type=local,dest=" + e.Dest
}

// Annotate sets docstrings on ExportLocal.
func (e *ExportLocal) Annotate(a infer.Annotator) {
	a.Describe(&e.Dest, "Output path.")
}

// Exists returns true if the exported directory exists and is populated with
// some files.
func (e *ExportLocal) Exists() bool {
	if e == nil {
		return true // Degenerate case.
	}
	files, err := os.ReadDir(e.Dest)
	if err != nil {
		return false
	}
	return len(files) > 0
}

// ExportTar is an export that uses the tar format for exporting.
type ExportTar struct {
	ExportLocal
}

func (e *ExportTar) String() string {
	if e == nil {
		return ""
	}
	return "type=tar,dest=" + e.Dest
}

// Exists returns true if a file exists at the expected path.
func (e *ExportTar) Exists() bool {
	if e == nil {
		return true // Degenerate case.
	}
	_, err := os.Stat(e.Dest)
	return err == nil
}

// ExportWithOCI is an export that support OCI media types.
type ExportWithOCI struct {
	OCI *bool `pulumi:"ociMediaTypes,optional"`
}

// Annotate sets defaults on ExportWithOCI.
func (c *ExportWithOCI) Annotate(a infer.Annotator) {
	a.SetDefault(&c.OCI, false)
	a.Describe(&c.OCI, "Use OCI media types in exporter manifests.")
}

func (c ExportWithOCI) String() string {
	if c.OCI == nil {
		return ""
	}
	return fmt.Sprintf("oci-mediatypes=%t", *c.OCI)
}

// ExportWithCompression is an export with options to configure compression
// settings.
type ExportWithCompression struct {
	Compression      *CompressionType `pulumi:"compression,optional"`
	CompressionLevel int              `pulumi:"compressionLevel,optional"`
	ForceCompression *bool            `pulumi:"forceCompression,optional"`
}

// Annotate sets docstrings and defaults on ExportWithCompression.
func (e *ExportWithCompression) Annotate(a infer.Annotator) {
	a.SetDefault(&e.Compression, Gzip)
	a.SetDefault(&e.CompressionLevel, 0)
	a.SetDefault(&e.ForceCompression, false)

	a.Describe(&e.Compression, "The compression type to use.")
	a.Describe(&e.CompressionLevel, "Compression level from 0 to 22.")
	a.Describe(&e.ForceCompression, "Forcefully apply compression.")
}

func (e ExportWithCompression) String() string {
	if e.CompressionLevel == 0 {
		return ""
	}
	parts := []string{}
	if e.Compression != nil {
		parts = append(parts, fmt.Sprintf("compression=%s", *e.Compression))
	}
	if e.CompressionLevel > 0 {
		cl := e.CompressionLevel
		if cl > 22 {
			cl = 22
		}
		parts = append(parts, fmt.Sprintf("compression-level=%d", cl))
	}
	if e.ForceCompression != nil {
		parts = append(parts, fmt.Sprintf("force-compression=%t", *e.ForceCompression))
	}
	return strings.Join(parts, ",")
}

// ExportWithNames is an export with configurable names (tags).
type ExportWithNames struct {
	Names []string `pulumi:"names,optional"`
}

func (e ExportWithNames) String() string {
	parts := []string{}
	for _, n := range e.Names {
		parts = append(parts, "name="+n)
	}
	return strings.Join(parts, ",")
}

// Annotate sets docstrings on ExportWithNames.
func (e *ExportWithNames) Annotate(a infer.Annotator) {
	a.Describe(
		&e.Names,
		"Specify images names to export. This is overridden if tags are already specified.",
	)
}

// ExportWithAnnotations is an export with configurable annotations.
type ExportWithAnnotations struct {
	Annotations map[string]string `pulumi:"annotations,optional"`
}

func (e ExportWithAnnotations) String() string {
	parts := []string{}
	for k, v := range e.Annotations {
		parts = append(parts, fmt.Sprintf("annotation.%s=%s", k, v))
	}
	slices.Sort(parts)
	return strings.Join(parts, ",")
}

// Annotate sets docstrings on ExportWithAnnotations.
func (e *ExportWithAnnotations) Annotate(a infer.Annotator) {
	a.Describe(&e.Annotations, dedent(`
		Attach an arbitrary key/value annotation to the image.
	`))
}

// isRegistryPush returns true if the ExportEntry results in an image pushed to
// a registry.
func isRegistryPush(export *controllerapi.ExportEntry) bool {
	// "type=registry" is shorthand for "type=image,push=true" so we only need
	// to check "image" types.
	return export.Type == "image" && export.Attrs["push"] == "true"
}
