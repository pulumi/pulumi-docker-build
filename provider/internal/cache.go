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
	"strings"

	controllerapi "github.com/docker/buildx/controller/pb"
	"github.com/docker/buildx/util/buildflags"

	"github.com/pulumi/pulumi-go-provider/infer"
)

var (
	_ fmt.Stringer                = (*CacheFrom)(nil)
	_ fmt.Stringer                = (*CacheFromAzureBlob)(nil)
	_ fmt.Stringer                = (*CacheFromGitHubActions)(nil)
	_ fmt.Stringer                = (*CacheFromLocal)(nil)
	_ fmt.Stringer                = (*CacheFromRegistry)(nil)
	_ fmt.Stringer                = (*CacheFromS3)(nil)
	_ fmt.Stringer                = (*CacheTo)(nil)
	_ fmt.Stringer                = (*CacheToAzureBlob)(nil)
	_ fmt.Stringer                = (*CacheToGitHubActions)(nil)
	_ fmt.Stringer                = (*CacheToInline)(nil)
	_ fmt.Stringer                = (*CacheToLocal)(nil)
	_ fmt.Stringer                = (*CacheToRegistry)(nil)
	_ fmt.Stringer                = (*CacheToS3)(nil)
	_ fmt.Stringer                = CacheWithCompression{}
	_ fmt.Stringer                = CacheWithIgnoreError{}
	_ fmt.Stringer                = CacheWithMode{}
	_ fmt.Stringer                = CacheWithOCI{}
	_ infer.Annotated             = (*CacheFrom)(nil)
	_ infer.Annotated             = (*CacheFromAzureBlob)(nil)
	_ infer.Annotated             = (*CacheFromGitHubActions)(nil)
	_ infer.Annotated             = (*CacheFromLocal)(nil)
	_ infer.Annotated             = (*CacheFromRegistry)(nil)
	_ infer.Annotated             = (*CacheFromS3)(nil)
	_ infer.Annotated             = (*CacheTo)(nil)
	_ infer.Annotated             = (*CacheToInline)(nil)
	_ infer.Annotated             = (*CacheToLocal)(nil)
	_ infer.Annotated             = (*CacheWithCompression)(nil)
	_ infer.Annotated             = (*CacheWithIgnoreError)(nil)
	_ infer.Annotated             = (*CacheWithMode)(nil)
	_ infer.Annotated             = (*CacheWithOCI)(nil)
	_ infer.Enum[CacheMode]       = (*CacheMode)(nil)
	_ infer.Enum[CompressionType] = (*CompressionType)(nil)
)

// CacheFromLocal pulls cache manifests from a local directory.
type CacheFromLocal struct {
	Src    string `pulumi:"src"`
	Digest string `pulumi:"digest,optional"`
}

// String returns the CLI-encoded value of these cache options, or an empty
// string if the receiver is nil.
func (c *CacheFromLocal) String() string {
	if c == nil {
		return ""
	}
	parts := []string{"type=local"}
	if c.Src != "" {
		parts = append(parts, "src="+c.Src)
	}
	if c.Digest != "" {
		parts = append(parts, "digest="+c.Digest)
	}
	return strings.Join(parts, ",")
}

// Annotate sets docstrings on CacheFromLocal.
func (c *CacheFromLocal) Annotate(a infer.Annotator) {
	a.Describe(&c.Src, "Path of the local directory where cache gets imported from.")
	a.Describe(&c.Digest, "Digest of manifest to import.")
}

// CacheFromRegistry pulls cache manifests from a registry ref.
type CacheFromRegistry struct {
	Ref string `pulumi:"ref"`
}

// Annotate sets docstrings on CacheFromRegistry.
func (c *CacheFromRegistry) Annotate(a infer.Annotator) {
	a.Describe(&c.Ref, "Fully qualified name of the cache image to import.")
}

// String returns the CLI-encoded value of these cache options, or an empty
// string if the receiver is nil.
func (c *CacheFromRegistry) String() string {
	if c == nil {
		return ""
	}
	return "type=registry,ref=" + c.Ref
}

// CacheWithOCI exposes OCI media type options.
type CacheWithOCI struct {
	OCI           *bool `pulumi:"ociMediaTypes,optional"`
	ImageManifest *bool `pulumi:"imageManifest,optional"`
}

// Annotate sets docstrings on CacheWithOCI.
func (c *CacheWithOCI) Annotate(a infer.Annotator) {
	a.Describe(&c.OCI, dedent(`
		Whether to use OCI media types in exported manifests. Defaults to
		"true".
	`))
	a.Describe(&c.ImageManifest, dedent(`
		Export cache manifest as an OCI-compatible image manifest instead of a
		manifest list. Requires "ociMediaTypes" to also be "true".

		Some registries like AWS ECR will not work with caching if this is
		"false".

		Defaults to "false" to match Docker's default behavior.
	`))

	a.SetDefault(&c.OCI, true)
	a.SetDefault(&c.ImageManifest, false)
}

// String returns the CLI-encoded value of these cache options, or an empty
// string if unknown.
func (c CacheWithOCI) String() string {
	if c.OCI == nil {
		return ""
	}
	parts := []string{fmt.Sprintf("oci-mediatypes=%t", *c.OCI)}
	if c.ImageManifest != nil {
		parts = append(parts, fmt.Sprintf("image-manifest=%t", *c.ImageManifest))
	}
	return strings.Join(parts, ",")
}

// CacheFromGitHubActions pulls cache manifests from the GitHub actions cache.
type CacheFromGitHubActions struct {
	URL   string `pulumi:"url,optional"`
	Token string `pulumi:"token,optional" provider:"secret"`
	Scope string `pulumi:"scope,optional"`
}

// Annotate sets docstrings on CacheFromGitHubActions.
func (c *CacheFromGitHubActions) Annotate(a infer.Annotator) {
	a.SetDefault(&c.URL, "", "ACTIONS_RUNTIME_URL")
	a.SetDefault(&c.Token, "", "ACTIONS_RUNTIME_TOKEN")
	a.SetDefault(&c.Scope, "", "buildkit")

	a.Describe(&c.URL, dedent(`
		The cache server URL to use for artifacts.

		Defaults to "$ACTIONS_RUNTIME_URL", although a separate action like
		"crazy-max/ghaction-github-runtime" is recommended to expose this
		environment variable to your jobs.
	`))
	a.Describe(&c.Token, dedent(`
		The GitHub Actions token to use. This is not a personal access tokens
		and is typically generated automatically as part of each job.

		Defaults to "$ACTIONS_RUNTIME_TOKEN", although a separate action like
		"crazy-max/ghaction-github-runtime" is recommended to expose this
		environment variable to your jobs.

	`))
	a.Describe(&c.Scope, dedent(`
		The scope to use for cache keys. Defaults to "buildkit".

		This should be set if building and caching multiple images in one
		workflow, otherwise caches will overwrite each other.
	`))
}

func (c *CacheFromGitHubActions) String() string {
	if c == nil {
		return ""
	}
	parts := []string{"type=gha"}
	if c.Scope != "" {
		parts = append(parts, "scope="+c.Scope)
	}
	if c.Token != "" {
		parts = append(parts, "token="+c.Token)
	}
	if c.URL != "" {
		parts = append(parts, "url="+c.URL)
	}
	return strings.Join(parts, ",")
}

// CacheFromAzureBlob pulls cache manifests from Azure
// blob storage.
type CacheFromAzureBlob struct {
	Name            string `pulumi:"name"`
	AccountURL      string `pulumi:"accountUrl,optional"`
	SecretAccessKey string `pulumi:"secretAccessKey,optional" provider:"secret"`
}

// String returns the CLI-encoded value of these cache options, or an empty
// string if the receiver is nil.
func (c *CacheFromAzureBlob) String() string {
	if c == nil {
		return ""
	}
	parts := []string{"type=azblob"}
	if c.Name != "" {
		parts = append(parts, "name="+c.Name)
	}
	if c.AccountURL != "" {
		parts = append(parts, "account_url="+c.AccountURL)
	}
	if c.SecretAccessKey != "" {
		parts = append(parts, "secret_access_key="+c.SecretAccessKey)
	}
	return strings.Join(parts, ",")
}

// Annotate sets docstrings on CacheFromAzureBlob.
func (c *CacheFromAzureBlob) Annotate(a infer.Annotator) {
	a.Describe(&c.Name, "The name of the cache image.")
	a.Describe(&c.AccountURL, "Base URL of the storage account.")
	a.Describe(&c.SecretAccessKey, "Blob storage account key.")
}

// CacheToAzureBlob pushes cache manifests to Azure blob storage.
type CacheToAzureBlob struct {
	CacheWithMode
	CacheWithIgnoreError

	CacheFromAzureBlob
}

func (c *CacheToAzureBlob) String() string {
	if c == nil {
		return ""
	}
	return join(&c.CacheFromAzureBlob, c.CacheWithMode, c.CacheWithIgnoreError)
}

// CacheFromS3 pulls cache manifests from S3-compatible APIs.
type CacheFromS3 struct {
	Region          string `pulumi:"region"`
	Bucket          string `pulumi:"bucket"`
	Name            string `pulumi:"name,optional"`
	EndpointURL     string `pulumi:"endpointUrl,optional"`
	BlobsPrefix     string `pulumi:"blobsPrefix,optional"`
	ManifestsPrefix string `pulumi:"manifestsPrefix,optional"`
	UsePathStyle    *bool  `pulumi:"usePathStyle,optional"`
	AccessKeyID     string `pulumi:"accessKeyId,optional"`
	SecretAccessKey string `pulumi:"secretAccessKey,optional" provider:"secret"`
	SessionToken    string `pulumi:"sessionToken,optional"    provider:"secret"`
}

// Annotate sets docstrings and defaults on CacheFromS3.
func (c *CacheFromS3) Annotate(a infer.Annotator) {
	a.SetDefault(&c.Region, "", "AWS_REGION")
	a.SetDefault(&c.AccessKeyID, "", "AWS_ACCESS_KEY_ID")
	a.SetDefault(&c.SecretAccessKey, "", "AWS_SECRET_ACCESS_KEY")
	a.SetDefault(&c.SessionToken, "", "AWS_SESSION_TOKEN")

	a.Describe(&c.Bucket, dedent(`
		Name of the S3 bucket.
	`))
	a.Describe(&c.Region, dedent(`
		The geographic location of the bucket. Defaults to "$AWS_REGION".
	`))
	a.Describe(&c.AccessKeyID, dedent(`
		Defaults to "$AWS_ACCESS_KEY_ID".
	`))
	a.Describe(&c.SecretAccessKey, dedent(`
		Defaults to "$AWS_SECRET_ACCESS_KEY".
	`))
	a.Describe(&c.SessionToken, dedent(`
		Defaults to "$AWS_SESSION_TOKEN".
	`))
	a.Describe(&c.BlobsPrefix, dedent(`
		Prefix to prepend to blob filenames.
	`))
	a.Describe(&c.EndpointURL, dedent(`
		Endpoint of the S3 bucket.
	`))
	a.Describe(&c.ManifestsPrefix, dedent(`
		Prefix to prepend on manifest filenames.
	`))
	a.Describe(&c.Name, dedent(`
		Name of the cache image.
	`))
	a.Describe(&c.UsePathStyle, dedent(`
		Uses "bucket" in the URL instead of hostname when "true".
	`))
}

// String returns the CLI-encoded value of these cache options, or an empty
// string if the receiver is nil.
func (c *CacheFromS3) String() string {
	if c == nil {
		return ""
	}
	parts := []string{"type=s3"}
	if c.Bucket != "" {
		parts = append(parts, "bucket="+c.Bucket)
	}
	if c.Name != "" {
		parts = append(parts, "name="+c.Name)
	}
	if c.EndpointURL != "" {
		parts = append(parts, "endpoint_url="+c.EndpointURL)
	}
	if c.BlobsPrefix != "" {
		parts = append(parts, "blobs_prefix="+c.BlobsPrefix)
	}
	if c.ManifestsPrefix != "" {
		parts = append(parts, "manifests_prefix="+c.ManifestsPrefix)
	}
	if c.UsePathStyle != nil {
		parts = append(parts, fmt.Sprintf("use_path_type=%t", *c.UsePathStyle))
	}
	if c.AccessKeyID != "" {
		parts = append(parts, "access_key_id="+c.AccessKeyID)
	}
	if c.SecretAccessKey != "" {
		parts = append(parts, "secret_access_key="+c.SecretAccessKey)
	}
	if c.SessionToken != "" {
		parts = append(parts, "session_token="+c.SessionToken)
	}

	return strings.Join(parts, ",")
}

// CacheWithMode is a cache that can configure its mode.
type CacheWithMode struct {
	Mode *CacheMode `pulumi:"mode,optional"`
}

// Annotate sets docstrings and defaults on CacheWithMode.
func (c *CacheWithMode) Annotate(a infer.Annotator) {
	m := Min
	a.SetDefault(&c.Mode, &m)
	a.Describe(&c.Mode, dedent(`
		The cache mode to use. Defaults to "min".
	`))
}

func (c CacheWithMode) String() string {
	if c.Mode == nil {
		return ""
	}
	return fmt.Sprintf("mode=%s", *c.Mode)
}

// CacheWithIgnoreError exposes an option to ignore errors during caching.
type CacheWithIgnoreError struct {
	IgnoreError *bool `pulumi:"ignoreError,optional"`
}

// Annotate sets docstrings and defaults on CacheWithIgnoreError.
func (c *CacheWithIgnoreError) Annotate(a infer.Annotator) {
	a.SetDefault(&c.IgnoreError, false)
	a.Describe(&c.IgnoreError, "Ignore errors caused by failed cache exports.")
}

func (c CacheWithIgnoreError) String() string {
	if c.IgnoreError == nil {
		return ""
	}
	return fmt.Sprintf("ignore-error=%t", *c.IgnoreError)
}

// CacheToS3 pushes cache manifests to an S3-compatible API.
type CacheToS3 struct {
	CacheWithMode
	CacheWithIgnoreError

	CacheFromS3
}

// String returns the CLI-encoded value of these cache options, or an empty
// string if the receiver is nil.
func (c *CacheToS3) String() string {
	if c == nil {
		return ""
	}
	return join(&c.CacheFromS3, c.CacheWithMode, c.CacheWithIgnoreError)
}

// Raw is a CLI-encoded cache entry appropriate for passing directly to the
// CLI. Useful if the Docker backend supports cache types not captured by our
// API, or if the user just prefers "type=..." inputs.
type Raw string

// String return the raw string as-is. This can be empty during previews where
// the user has provided a value but it is unknown.
func (c Raw) String() string {
	return string(c)
}

// CacheFrom is a "union" type for all of our available `--cache-from` options.
type CacheFrom struct {
	Local    *CacheFromLocal         `pulumi:"local,optional"`
	Registry *CacheFromRegistry      `pulumi:"registry,optional"`
	GHA      *CacheFromGitHubActions `pulumi:"gha,optional"`
	AZBlob   *CacheFromAzureBlob     `pulumi:"azblob,optional"`
	S3       *CacheFromS3            `pulumi:"s3,optional"`
	Raw      Raw                     `pulumi:"raw,optional"`

	Disabled bool `pulumi:"disabled,optional"`
}

// Annotate sets docstrings and defaults on CacheFrom.
func (c *CacheFrom) Annotate(a infer.Annotator) {
	a.Describe(&c.Local, dedent(`
		A simple backend which caches images on your local filesystem.
	`))
	a.Describe(&c.Registry, dedent(`
		Upload build caches to remote registries.
	`))
	a.Describe(&c.GHA, dedent(`
		Recommended for use with GitHub Actions workflows.

		An action like "crazy-max/ghaction-github-runtime" is recommended to
		expose appropriate credentials to your GitHub workflow.
	`))
	a.Describe(&c.AZBlob, dedent(`
		Upload build caches to Azure's blob storage service.
	`))
	a.Describe(&c.S3, dedent(`
		Upload build caches to AWS S3 or an S3-compatible services such as
		MinIO.
	`))
	a.Describe(&c.Raw, dedent(`
		A raw string as you would provide it to the Docker CLI (e.g.,
		"type=inline").
	`))

	a.Describe(&c.Disabled, dedent(`
		When "true" this entry will be excluded. Defaults to "false".
	`))
}

// String returns a CLI-encoded value for this `--cache-from` entry, or an
// empty string if disabled. `validate` should be called to ensure only one
// entry was set.
func (c CacheFrom) String() string {
	if c.Disabled {
		return ""
	}
	return join(c.Local, c.Registry, c.GHA, c.AZBlob, c.S3, c.Raw)
}

func (c CacheFrom) validate(preview bool) (*controllerapi.CacheOptionsEntry, error) {
	if strings.Count(c.String(), "type=") > 1 {
		return nil, errors.New("cacheFrom should only specify one cache type")
	}
	parsed, err := buildflags.ParseCacheEntry([]string{c.String()})
	if err != nil {
		return nil, err
	}
	if len(parsed) == 0 {
		// This can happen for example if we have a GHA cache but no GitHub
		// environment variables set.
		// Shouldn't happen...
		return nil, nil
	}
	return parsed[0], nil
}

// CacheToInline embeds cache information directly into an image.
type CacheToInline struct{}

// String returns the CLI-encoded value of these cache options, or an empty
// string if unknown.
func (c *CacheToInline) String() string {
	if c == nil {
		return ""
	}
	return "type=inline"
}

// Annotate sets docstrings on CacheToInline.
func (c *CacheToInline) Annotate(a infer.Annotator) {
	a.Describe(&c, dedent(`
		Include an inline cache with the exported image.
	`))
}

// CacheToLocal writes cache manifests to a local directory.
type CacheToLocal struct {
	CacheWithCompression
	CacheWithIgnoreError
	CacheWithMode

	Dest string `pulumi:"dest"`
}

// Annotate sets docstrings on CacheToLocal.
func (c *CacheToLocal) Annotate(a infer.Annotator) {
	a.Describe(&c.Dest, dedent(`
		Path of the local directory to export the cache.
	`))
}

// String returns the CLI-encoded value of these cache options, or an empty
// string if the receiver is nil.
func (c *CacheToLocal) String() string {
	if c == nil {
		return ""
	}
	return join(
		Raw("type=local,dest="+c.Dest),
		c.CacheWithCompression,
		c.CacheWithIgnoreError,
	)
}

// CacheToRegistry pushes cache manifests to a remote registry.
type CacheToRegistry struct {
	CacheWithMode
	CacheWithIgnoreError
	CacheWithOCI
	CacheWithCompression

	CacheFromRegistry
}

// String returns the CLI-encoded value of these cache options, or an empty
// string if the receiver is nil.
func (c *CacheToRegistry) String() string {
	if c == nil {
		return ""
	}
	return join(
		&c.CacheFromRegistry,
		c.CacheWithMode,
		c.CacheWithIgnoreError,
		c.CacheWithOCI,
		c.CacheWithCompression,
	)
}

// CacheWithCompression is a cache with options to configure compression
// settings.
type CacheWithCompression struct {
	Compression      *CompressionType `pulumi:"compression,optional"`
	CompressionLevel int              `pulumi:"compressionLevel,optional"`
	ForceCompression *bool            `pulumi:"forceCompression,optional"`
}

// Annotate sets docstrings and defaults on CacheWithCompression.
func (c *CacheWithCompression) Annotate(a infer.Annotator) {
	gz := Gzip
	a.SetDefault(&c.Compression, &gz)
	a.SetDefault(&c.CompressionLevel, 0)
	a.SetDefault(&c.ForceCompression, false)

	a.Describe(&c.Compression, "The compression type to use.")
	a.Describe(&c.CompressionLevel, "Compression level from 0 to 22.")
	a.Describe(&c.ForceCompression, "Forcefully apply compression.")
}

// String returns the CLI-encoded value of these cache options, or an empty
// string if the receiver is nil.
func (c CacheWithCompression) String() string {
	if c.CompressionLevel == 0 {
		return ""
	}
	parts := []string{}
	if c.Compression != nil {
		parts = append(parts, fmt.Sprintf("compression=%s", *c.Compression))
	}
	if c.CompressionLevel > 0 {
		cl := c.CompressionLevel
		if cl > 22 {
			cl = 22
		}
		parts = append(parts, fmt.Sprintf("compression-level=%d", cl))
	}
	if c.ForceCompression != nil {
		parts = append(parts, fmt.Sprintf("force-compression=%t", *c.ForceCompression))
	}
	return strings.Join(parts, ",")
}

// CacheToGitHubActions pushes cache manifests to the GitHub Actions cache
// backend.
type CacheToGitHubActions struct {
	CacheWithMode
	CacheWithIgnoreError

	CacheFromGitHubActions
}

// String returns the CLI-encoded value of these cache options, or an empty
// string if the receiver is nil.
func (c *CacheToGitHubActions) String() string {
	if c == nil {
		return ""
	}
	return join(&c.CacheFromGitHubActions, c.CacheWithMode, c.CacheWithIgnoreError)
}

// CacheTo is a "union" type for all of our available `--cache-to` options.
type CacheTo struct {
	Inline   *CacheToInline        `pulumi:"inline,optional"`
	Local    *CacheToLocal         `pulumi:"local,optional"`
	Registry *CacheToRegistry      `pulumi:"registry,optional"`
	GHA      *CacheToGitHubActions `pulumi:"gha,optional"`
	AZBlob   *CacheToAzureBlob     `pulumi:"azblob,optional"`
	S3       *CacheToS3            `pulumi:"s3,optional"`
	Raw      Raw                   `pulumi:"raw,optional"`

	Disabled bool `pulumi:"disabled,optional"`
}

// Annotate sets docstrings and defaults on CacheTo.
func (c *CacheTo) Annotate(a infer.Annotator) {
	a.Describe(&c.Inline, dedent(`
		The inline cache storage backend is the simplest implementation to get
		started with, but it does not handle multi-stage builds. Consider the
		"registry" cache backend instead.
	`))
	a.Describe(&c.Local, dedent(`
		A simple backend which caches imagines on your local filesystem.
	`))
	a.Describe(&c.Registry, dedent(`
		Push caches to remote registries. Incompatible with the "docker" build
		driver.
	`))
	a.Describe(&c.GHA, dedent(`
		Recommended for use with GitHub Actions workflows.

		An action like "crazy-max/ghaction-github-runtime" is recommended to
		expose appropriate credentials to your GitHub workflow.
	`))
	a.Describe(&c.AZBlob, dedent(`
		Push cache to Azure's blob storage service.
	`))
	a.Describe(&c.S3, dedent(`
		Push cache to AWS S3 or S3-compatible services such as MinIO.
	`))
	a.Describe(&c.Raw, dedent(`
		A raw string as you would provide it to the Docker CLI (e.g.,
		"type=inline")`,
	))

	a.Describe(&c.Disabled, dedent(`
		When "true" this entry will be excluded. Defaults to "false".
	`))
}

// String returns a CLI-encoded value for this `--cache-to` entry, or an
// empty string if disabled. `validate` should be called to ensure only one
// entry was set.
func (c CacheTo) String() string {
	if c.Disabled {
		return ""
	}
	return join(c.Inline, c.Local, c.Registry, c.GHA, c.AZBlob, c.S3, c.Raw)
}

func (c CacheTo) validate(preview bool) (*controllerapi.CacheOptionsEntry, error) {
	if strings.Count(c.String(), "type=") > 1 {
		return nil, errors.New("cacheTo should only specify one cache type")
	}
	parsed, err := buildflags.ParseCacheEntry([]string{c.String()})
	if err != nil {
		return nil, err
	}
	if len(parsed) == 0 {
		// This can happen for example if we have a GHA cache but no GitHub
		// environment variables set.
		// Shouldn't happen...
		return nil, nil
	}
	return parsed[0], nil
}

// CacheMode controls the complexity of exported cache manifests.
type CacheMode string

const (
	Min CacheMode = "min" // Min cache mode.
	Max CacheMode = "max" // Max cache mode.
)

// Values returns all valid CacheMode values for SDK generation.
func (CacheMode) Values() []infer.EnumValue[CacheMode] {
	return []infer.EnumValue[CacheMode]{
		{
			Value:       Min,
			Description: "Only layers that are exported into the resulting image are cached.",
		},
		{
			Value:       Max,
			Description: "All layers are cached, even those of intermediate steps.",
		},
	}
}

// CompressionType is the algorithm used for compressing blobs.
type CompressionType string

const (
	Gzip    CompressionType = "gzip"    // Gzip compression.
	Estargz CompressionType = "estargz" // Estargz compression.
	Zstd    CompressionType = "zstd"    // Zstd compression.
)

// Values returns all valid CompressionType values for SDK generation.
func (CompressionType) Values() []infer.EnumValue[CompressionType] {
	return []infer.EnumValue[CompressionType]{
		{Value: Gzip, Description: "Use `gzip` for compression."},
		{Value: Estargz, Description: "Use `estargz` for compression."},
		{Value: Zstd, Description: "Use `zstd` for compression."},
	}
}

type joiner struct{ sep string }

func (j joiner) join(ss ...fmt.Stringer) string {
	parts := []string{}
	for _, s := range ss {
		p := s.String()
		if p == "" {
			continue
		}
		parts = append(parts, p)
	}
	return strings.Join(parts, j.sep)
}

func join(ss ...fmt.Stringer) string {
	return joiner{","}.join(ss...)
}
