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
	"testing"

	controllerapi "github.com/docker/buildx/controller/pb"
	"github.com/docker/buildx/util/buildflags"
	"github.com/stretchr/testify/assert"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func TestValidateExport(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		e         Export
		givenTags []string
		preview   bool

		wantExp *controllerapi.ExportEntry
		wantErr string
	}{
		{
			name:      "raw - no push on preview",
			preview:   true,
			e:         Export{Raw: "type=registry"},
			givenTags: []string{"docker.io/foo/bar"},
			wantExp: &controllerapi.ExportEntry{
				Type:  "image",
				Attrs: map[string]string{"push": "false"},
			},
		},
		{
			name:    "raw - push requires tags",
			e:       Export{Raw: "type=registry"},
			wantErr: "tag or export name is needed",
		},
		{
			name:      "registry - no push on preview",
			preview:   true,
			e:         Export{Registry: &ExportRegistry{}},
			givenTags: []string{"docker.io/foo/bar"},
			wantExp: &controllerapi.ExportEntry{
				Type:  "image",
				Attrs: map[string]string{"push": "false"},
			},
		},
		{
			name:    "registry - push requires tags",
			e:       Export{Registry: &ExportRegistry{}},
			wantErr: "tag or export name is needed",
		},
		{
			name:    "over-specified",
			e:       Export{Raw: "type=registry", Registry: &ExportRegistry{}},
			wantErr: "specify one export type",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e, err := tt.e.validate(tt.preview, tt.givenTags)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.wantErr)
			}
			if tt.wantExp != nil {
				assert.Equal(t, tt.wantExp.Type, e.Type)
				assert.Equal(t, tt.wantExp.Attrs, e.Attrs)
			}
		})
	}
}

func TestExportString(t *testing.T) {
	t.Parallel()
	gzip := Gzip
	tests := []struct {
		name  string
		given Export
		want  string
	}{
		{
			name:  "tar",
			given: Export{Tar: &ExportTar{ExportLocal: ExportLocal{Dest: "/foo"}}},
			want:  "type=tar,dest=/foo",
		},
		{
			name:  "local",
			given: Export{Local: &ExportLocal{Dest: "/bar"}},
			want:  "type=local,dest=/bar",
		},
		{
			name: "registry-with-compression",
			given: Export{Registry: &ExportRegistry{
				ExportImage: ExportImage{
					ExportWithCompression: ExportWithCompression{
						Compression:      &gzip,
						CompressionLevel: 100,
						ForceCompression: pulumi.BoolRef(true),
					},
				},
			}},
			want: "type=registry,compression=gzip,compression-level=22,force-compression=true",
		},
		{
			name: "registry-without-push",
			given: Export{Registry: &ExportRegistry{
				ExportImage: ExportImage{
					Push: pulumi.BoolRef(false),
				},
			}},
			want: "type=registry,push=false",
		},
		{
			name: "image",
			given: Export{
				Image: &ExportImage{
					Push:               pulumi.BoolRef(true),
					PushByDigest:       pulumi.BoolRef(true),
					Insecure:           pulumi.BoolRef(true),
					DanglingNamePrefix: "prefix",
					Unpack:             pulumi.BoolRef(true),
					Store:              pulumi.BoolRef(false),
				},
			},
			want: "type=image,push=true,push-by-digest=true,insecure=true,dangling-name-prefix=prefix,unpack=true,store=false",
		},
		{
			name: "oci-with-names",
			given: Export{OCI: &ExportOCI{
				ExportDocker: ExportDocker{
					ExportWithNames: ExportWithNames{
						Names: []string{"foo", "bar"},
					},
				},
			}},
			want: "type=oci,name=foo,name=bar",
		},
		{
			name: "docker-with-annotations",
			given: Export{Docker: &ExportDocker{
				ExportWithAnnotations: ExportWithAnnotations{
					Annotations: map[string]string{
						"foo": "bar",
						"boo": "baz",
					},
				},
			}},
			want: "type=docker,annotation.boo=baz,annotation.foo=bar",
		},
		{
			name:  "raw",
			given: Export{Raw: Raw("type=docker")},
			want:  "type=docker",
		},
		{
			name:  "disabled",
			given: Export{Raw: Raw("type=docker"), Disabled: true},
			want:  "",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.given.String()
			assert.Equal(t, tt.want, tt.given.String())

			if tt.want != "" {
				// Our output should be parsable by Docker.
				_, err := buildflags.ParseExports([]string{actual})
				assert.NoError(t, err)
			}
		})
	}
}

func TestExportPushed(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		e    Export
		want bool
	}{
		{
			name: "raw registry",
			e:    Export{Raw: "type=registry"},
			want: true,
		},
		{
			name: "raw image",
			e:    Export{Raw: "type=image"},
			want: false,
		},
		{
			name: "registry with no push",
			e:    Export{Registry: &ExportRegistry{}},
			want: true,
		},
		{
			name: "registry with explicit push",
			e:    Export{Registry: &ExportRegistry{ExportImage{Push: pulumi.BoolRef(false)}}},
			want: false,
		},
		{
			name: "image with explicit push",
			e:    Export{Image: &ExportImage{Push: pulumi.BoolRef(true)}},
			want: true,
		},
		{
			name: "local",
			e:    Export{Local: &ExportLocal{}},
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			actual := tt.e.pushed()
			assert.Equal(t, tt.want, actual)
		})
	}
}
