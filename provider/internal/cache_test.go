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
	"testing"

	"github.com/docker/buildx/util/buildflags"
	"github.com/stretchr/testify/assert"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func TestCacheString(t *testing.T) {
	t.Parallel()
	gzip := Gzip

	tests := []struct {
		name  string
		given fmt.Stringer
		want  string
	}{
		{
			name: "s3",
			given: CacheTo{S3: &CacheToS3{
				CacheFromS3: CacheFromS3{
					Region:          "us-west-2",
					Bucket:          "bucket-foo",
					Name:            "myname",
					EndpointURL:     "https://some.endpoint",
					BlobsPrefix:     "blob-prefix",
					ManifestsPrefix: "manifest-prefix",
					UsePathStyle:    pulumi.BoolRef(true),
					AccessKeyID:     "access-key-id",
					SecretAccessKey: "secret-key",
					SessionToken:    "session",
				},
			}},
			//nolint:lll // Taken from AWS reference docs.
			want: "type=s3,bucket=bucket-foo,name=myname,endpoint_url=https://some.endpoint,blobs_prefix=blob-prefix,manifests_prefix=manifest-prefix,use_path_type=true,access_key_id=access-key-id,secret_access_key=secret-key,session_token=session",
		},
		{
			name:  "gha",
			given: CacheTo{GHA: &CacheToGitHubActions{}},
			want:  "type=gha",
		},
		{
			name: "gha-with-url-and-token",
			given: CacheTo{GHA: &CacheToGitHubActions{
				CacheFromGitHubActions: CacheFromGitHubActions{
					URL:   "https://github.com/user/repo",
					Token: "token",
				},
			}},
			want: "type=gha,token=token,url=https://github.com/user/repo",
		},
		{
			name:  "from-local",
			given: CacheFrom{Local: &CacheFromLocal{Src: "/foo/bar"}},
			want:  "type=local,src=/foo/bar",
		},
		{
			name:  "to-local",
			given: CacheTo{Local: &CacheToLocal{Dest: "/foo/bar"}},
			want:  "type=local,dest=/foo/bar",
		},
		{
			name:  "inline",
			given: CacheTo{Inline: &CacheToInline{}},
			want:  "type=inline",
		},
		{
			name:  "raw",
			given: CacheTo{Raw: Raw("type=gha")},
			want:  "type=gha",
		},
		{
			name: "compression",
			given: CacheTo{Local: &CacheToLocal{
				Dest: "/foo",
				CacheWithCompression: CacheWithCompression{
					Compression:      &gzip,
					CompressionLevel: 100,
					ForceCompression: pulumi.BoolRef(true),
				},
			}},
			want: "type=local,dest=/foo,compression=gzip,compression-level=22,force-compression=true",
		},
		{
			name: "ignore-error",
			given: CacheTo{
				AZBlob: &CacheToAzureBlob{
					CacheWithIgnoreError: CacheWithIgnoreError{pulumi.BoolRef(true)},
				},
			},
			want: "type=azblob,ignore-error=true",
		},
		{
			name: "oci",
			given: CacheTo{
				Registry: &CacheToRegistry{
					CacheFromRegistry: CacheFromRegistry{Ref: "docker.io/foo/bar:baz"},
					CacheWithOCI: CacheWithOCI{
						OCI:           pulumi.BoolRef(true),
						ImageManifest: pulumi.BoolRef(true),
					},
				},
			},
			want: "type=registry,ref=docker.io/foo/bar:baz,oci-mediatypes=true,image-manifest=true",
		},
		{
			name: "disabled-to",
			given: CacheTo{
				Raw:      Raw("type=gha"),
				Disabled: true,
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := tt.given.String()
			assert.Equal(t, tt.want, actual)

			if tt.want != "" {
				// Our output should be parsable by Docker.
				_, err := buildflags.ParseCacheEntry([]string{actual})
				assert.NoError(t, err)
			}
		})
	}
}
