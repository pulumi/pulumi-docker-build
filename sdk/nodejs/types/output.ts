// *** WARNING: this file was generated by pulumi-language-nodejs. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as inputs from "../types/input";
import * as outputs from "../types/output";
import * as enums from "../types/enums";

import * as utilities from "../utilities";

export interface BuildContext {
    /**
     * Resources to use for build context.
     *
     * The location can be:
     * * A relative or absolute path to a local directory (`.`, `./app`,
     *   `/app`, etc.).
     * * A remote URL of a Git repository, tarball, or plain text file
     *   (`https://github.com/user/myrepo.git`, `http://server/context.tar.gz`,
     *   etc.).
     */
    location: string;
    /**
     * Additional build contexts to use.
     *
     * These contexts are accessed with `FROM name` or `--from=name`
     * statements when using Dockerfile 1.4+ syntax.
     *
     * Values can be local paths, HTTP URLs, or  `docker-image://` images.
     */
    named?: {[key: string]: outputs.Context};
}

export interface BuilderConfig {
    /**
     * Name of an existing buildx builder to use.
     *
     * Only `docker-container`, `kubernetes`, or `remote` drivers are
     * supported. The legacy `docker` driver is not supported.
     *
     * Equivalent to Docker's `--builder` flag.
     */
    name?: string;
}

export interface CacheFrom {
    /**
     * Upload build caches to Azure's blob storage service.
     */
    azblob?: outputs.CacheFromAzureBlob;
    /**
     * When `true` this entry will be excluded. Defaults to `false`.
     */
    disabled?: boolean;
    /**
     * Recommended for use with GitHub Actions workflows.
     *
     * An action like `crazy-max/ghaction-github-runtime` is recommended to
     * expose appropriate credentials to your GitHub workflow.
     */
    gha?: outputs.CacheFromGitHubActions;
    /**
     * A simple backend which caches images on your local filesystem.
     */
    local?: outputs.CacheFromLocal;
    /**
     * A raw string as you would provide it to the Docker CLI (e.g.,
     * `type=inline`).
     */
    raw?: string;
    /**
     * Upload build caches to remote registries.
     */
    registry?: outputs.CacheFromRegistry;
    /**
     * Upload build caches to AWS S3 or an S3-compatible services such as
     * MinIO.
     */
    s3?: outputs.CacheFromS3;
}
/**
 * cacheFromProvideDefaults sets the appropriate defaults for CacheFrom
 */
export function cacheFromProvideDefaults(val: CacheFrom): CacheFrom {
    return {
        ...val,
        gha: (val.gha ? outputs.cacheFromGitHubActionsProvideDefaults(val.gha) : undefined),
        s3: (val.s3 ? outputs.cacheFromS3ProvideDefaults(val.s3) : undefined),
    };
}

export interface CacheFromAzureBlob {
    /**
     * Base URL of the storage account.
     */
    accountUrl?: string;
    /**
     * The name of the cache image.
     */
    name: string;
    /**
     * Blob storage account key.
     */
    secretAccessKey?: string;
}

export interface CacheFromGitHubActions {
    /**
     * The scope to use for cache keys. Defaults to `buildkit`.
     *
     * This should be set if building and caching multiple images in one
     * workflow, otherwise caches will overwrite each other.
     */
    scope?: string;
    /**
     * The GitHub Actions token to use. This is not a personal access tokens
     * and is typically generated automatically as part of each job.
     *
     * Defaults to `$ACTIONS_RUNTIME_TOKEN`, although a separate action like
     * `crazy-max/ghaction-github-runtime` is recommended to expose this
     * environment variable to your jobs.
     */
    token?: string;
    /**
     * The cache server URL to use for artifacts.
     *
     * Defaults to `$ACTIONS_CACHE_URL`, although a separate action like
     * `crazy-max/ghaction-github-runtime` is recommended to expose this
     * environment variable to your jobs.
     */
    url?: string;
}
/**
 * cacheFromGitHubActionsProvideDefaults sets the appropriate defaults for CacheFromGitHubActions
 */
export function cacheFromGitHubActionsProvideDefaults(val: CacheFromGitHubActions): CacheFromGitHubActions {
    return {
        ...val,
        scope: (val.scope) ?? "buildkit",
        token: (val.token) ?? (utilities.getEnv("ACTIONS_RUNTIME_TOKEN") || ""),
        url: (val.url) ?? (utilities.getEnv("ACTIONS_CACHE_URL") || ""),
    };
}

export interface CacheFromLocal {
    /**
     * Digest of manifest to import.
     */
    digest?: string;
    /**
     * Path of the local directory where cache gets imported from.
     */
    src: string;
}

export interface CacheFromRegistry {
    /**
     * Fully qualified name of the cache image to import.
     */
    ref: string;
}

export interface CacheFromS3 {
    /**
     * Defaults to `$AWS_ACCESS_KEY_ID`.
     */
    accessKeyId?: string;
    /**
     * Prefix to prepend to blob filenames.
     */
    blobsPrefix?: string;
    /**
     * Name of the S3 bucket.
     */
    bucket: string;
    /**
     * Endpoint of the S3 bucket.
     */
    endpointUrl?: string;
    /**
     * Prefix to prepend on manifest filenames.
     */
    manifestsPrefix?: string;
    /**
     * Name of the cache image.
     */
    name?: string;
    /**
     * The geographic location of the bucket. Defaults to `$AWS_REGION`.
     */
    region: string;
    /**
     * Defaults to `$AWS_SECRET_ACCESS_KEY`.
     */
    secretAccessKey?: string;
    /**
     * Defaults to `$AWS_SESSION_TOKEN`.
     */
    sessionToken?: string;
    /**
     * Uses `bucket` in the URL instead of hostname when `true`.
     */
    usePathStyle?: boolean;
}
/**
 * cacheFromS3ProvideDefaults sets the appropriate defaults for CacheFromS3
 */
export function cacheFromS3ProvideDefaults(val: CacheFromS3): CacheFromS3 {
    return {
        ...val,
        accessKeyId: (val.accessKeyId) ?? (utilities.getEnv("AWS_ACCESS_KEY_ID") || ""),
        region: (val.region) ?? (utilities.getEnv("AWS_REGION") || ""),
        secretAccessKey: (val.secretAccessKey) ?? (utilities.getEnv("AWS_SECRET_ACCESS_KEY") || ""),
        sessionToken: (val.sessionToken) ?? (utilities.getEnv("AWS_SESSION_TOKEN") || ""),
    };
}

export interface CacheTo {
    /**
     * Push cache to Azure's blob storage service.
     */
    azblob?: outputs.CacheToAzureBlob;
    /**
     * When `true` this entry will be excluded. Defaults to `false`.
     */
    disabled?: boolean;
    /**
     * Recommended for use with GitHub Actions workflows.
     *
     * An action like `crazy-max/ghaction-github-runtime` is recommended to
     * expose appropriate credentials to your GitHub workflow.
     */
    gha?: outputs.CacheToGitHubActions;
    /**
     * The inline cache storage backend is the simplest implementation to get
     * started with, but it does not handle multi-stage builds. Consider the
     * `registry` cache backend instead.
     */
    inline?: outputs.CacheToInline;
    /**
     * A simple backend which caches imagines on your local filesystem.
     */
    local?: outputs.CacheToLocal;
    /**
     * A raw string as you would provide it to the Docker CLI (e.g.,
     * `type=inline`)
     */
    raw?: string;
    /**
     * Push caches to remote registries. Incompatible with the `docker` build
     * driver.
     */
    registry?: outputs.CacheToRegistry;
    /**
     * Push cache to AWS S3 or S3-compatible services such as MinIO.
     */
    s3?: outputs.CacheToS3;
}
/**
 * cacheToProvideDefaults sets the appropriate defaults for CacheTo
 */
export function cacheToProvideDefaults(val: CacheTo): CacheTo {
    return {
        ...val,
        azblob: (val.azblob ? outputs.cacheToAzureBlobProvideDefaults(val.azblob) : undefined),
        gha: (val.gha ? outputs.cacheToGitHubActionsProvideDefaults(val.gha) : undefined),
        local: (val.local ? outputs.cacheToLocalProvideDefaults(val.local) : undefined),
        registry: (val.registry ? outputs.cacheToRegistryProvideDefaults(val.registry) : undefined),
        s3: (val.s3 ? outputs.cacheToS3ProvideDefaults(val.s3) : undefined),
    };
}

export interface CacheToAzureBlob {
    /**
     * Base URL of the storage account.
     */
    accountUrl?: string;
    /**
     * Ignore errors caused by failed cache exports.
     */
    ignoreError?: boolean;
    /**
     * The cache mode to use. Defaults to `min`.
     */
    mode?: enums.CacheMode;
    /**
     * The name of the cache image.
     */
    name: string;
    /**
     * Blob storage account key.
     */
    secretAccessKey?: string;
}
/**
 * cacheToAzureBlobProvideDefaults sets the appropriate defaults for CacheToAzureBlob
 */
export function cacheToAzureBlobProvideDefaults(val: CacheToAzureBlob): CacheToAzureBlob {
    return {
        ...val,
        ignoreError: (val.ignoreError) ?? false,
        mode: (val.mode) ?? "min",
    };
}

export interface CacheToGitHubActions {
    /**
     * Ignore errors caused by failed cache exports.
     */
    ignoreError?: boolean;
    /**
     * The cache mode to use. Defaults to `min`.
     */
    mode?: enums.CacheMode;
    /**
     * The scope to use for cache keys. Defaults to `buildkit`.
     *
     * This should be set if building and caching multiple images in one
     * workflow, otherwise caches will overwrite each other.
     */
    scope?: string;
    /**
     * The GitHub Actions token to use. This is not a personal access tokens
     * and is typically generated automatically as part of each job.
     *
     * Defaults to `$ACTIONS_RUNTIME_TOKEN`, although a separate action like
     * `crazy-max/ghaction-github-runtime` is recommended to expose this
     * environment variable to your jobs.
     */
    token?: string;
    /**
     * The cache server URL to use for artifacts.
     *
     * Defaults to `$ACTIONS_CACHE_URL`, although a separate action like
     * `crazy-max/ghaction-github-runtime` is recommended to expose this
     * environment variable to your jobs.
     */
    url?: string;
}
/**
 * cacheToGitHubActionsProvideDefaults sets the appropriate defaults for CacheToGitHubActions
 */
export function cacheToGitHubActionsProvideDefaults(val: CacheToGitHubActions): CacheToGitHubActions {
    return {
        ...val,
        ignoreError: (val.ignoreError) ?? false,
        mode: (val.mode) ?? "min",
        scope: (val.scope) ?? "buildkit",
        token: (val.token) ?? (utilities.getEnv("ACTIONS_RUNTIME_TOKEN") || ""),
        url: (val.url) ?? (utilities.getEnv("ACTIONS_CACHE_URL") || ""),
    };
}

/**
 * Include an inline cache with the exported image.
 */
export interface CacheToInline {
}

export interface CacheToLocal {
    /**
     * The compression type to use.
     */
    compression?: enums.CompressionType;
    /**
     * Compression level from 0 to 22.
     */
    compressionLevel?: number;
    /**
     * Path of the local directory to export the cache.
     */
    dest: string;
    /**
     * Forcefully apply compression.
     */
    forceCompression?: boolean;
    /**
     * Ignore errors caused by failed cache exports.
     */
    ignoreError?: boolean;
    /**
     * The cache mode to use. Defaults to `min`.
     */
    mode?: enums.CacheMode;
}
/**
 * cacheToLocalProvideDefaults sets the appropriate defaults for CacheToLocal
 */
export function cacheToLocalProvideDefaults(val: CacheToLocal): CacheToLocal {
    return {
        ...val,
        compression: (val.compression) ?? "gzip",
        compressionLevel: (val.compressionLevel) ?? 0,
        forceCompression: (val.forceCompression) ?? false,
        ignoreError: (val.ignoreError) ?? false,
        mode: (val.mode) ?? "min",
    };
}

export interface CacheToRegistry {
    /**
     * The compression type to use.
     */
    compression?: enums.CompressionType;
    /**
     * Compression level from 0 to 22.
     */
    compressionLevel?: number;
    /**
     * Forcefully apply compression.
     */
    forceCompression?: boolean;
    /**
     * Ignore errors caused by failed cache exports.
     */
    ignoreError?: boolean;
    /**
     * Export cache manifest as an OCI-compatible image manifest instead of a
     * manifest list. Requires `ociMediaTypes` to also be `true`.
     *
     * Some registries like AWS ECR will not work with caching if this is
     * `false`.
     *
     * Defaults to `false` to match Docker's default behavior.
     */
    imageManifest?: boolean;
    /**
     * The cache mode to use. Defaults to `min`.
     */
    mode?: enums.CacheMode;
    /**
     * Whether to use OCI media types in exported manifests. Defaults to
     * `true`.
     */
    ociMediaTypes?: boolean;
    /**
     * Fully qualified name of the cache image to import.
     */
    ref: string;
}
/**
 * cacheToRegistryProvideDefaults sets the appropriate defaults for CacheToRegistry
 */
export function cacheToRegistryProvideDefaults(val: CacheToRegistry): CacheToRegistry {
    return {
        ...val,
        compression: (val.compression) ?? "gzip",
        compressionLevel: (val.compressionLevel) ?? 0,
        forceCompression: (val.forceCompression) ?? false,
        ignoreError: (val.ignoreError) ?? false,
        imageManifest: (val.imageManifest) ?? false,
        mode: (val.mode) ?? "min",
        ociMediaTypes: (val.ociMediaTypes) ?? true,
    };
}

export interface CacheToS3 {
    /**
     * Defaults to `$AWS_ACCESS_KEY_ID`.
     */
    accessKeyId?: string;
    /**
     * Prefix to prepend to blob filenames.
     */
    blobsPrefix?: string;
    /**
     * Name of the S3 bucket.
     */
    bucket: string;
    /**
     * Endpoint of the S3 bucket.
     */
    endpointUrl?: string;
    /**
     * Ignore errors caused by failed cache exports.
     */
    ignoreError?: boolean;
    /**
     * Prefix to prepend on manifest filenames.
     */
    manifestsPrefix?: string;
    /**
     * The cache mode to use. Defaults to `min`.
     */
    mode?: enums.CacheMode;
    /**
     * Name of the cache image.
     */
    name?: string;
    /**
     * The geographic location of the bucket. Defaults to `$AWS_REGION`.
     */
    region: string;
    /**
     * Defaults to `$AWS_SECRET_ACCESS_KEY`.
     */
    secretAccessKey?: string;
    /**
     * Defaults to `$AWS_SESSION_TOKEN`.
     */
    sessionToken?: string;
    /**
     * Uses `bucket` in the URL instead of hostname when `true`.
     */
    usePathStyle?: boolean;
}
/**
 * cacheToS3ProvideDefaults sets the appropriate defaults for CacheToS3
 */
export function cacheToS3ProvideDefaults(val: CacheToS3): CacheToS3 {
    return {
        ...val,
        accessKeyId: (val.accessKeyId) ?? (utilities.getEnv("AWS_ACCESS_KEY_ID") || ""),
        ignoreError: (val.ignoreError) ?? false,
        mode: (val.mode) ?? "min",
        region: (val.region) ?? (utilities.getEnv("AWS_REGION") || ""),
        secretAccessKey: (val.secretAccessKey) ?? (utilities.getEnv("AWS_SECRET_ACCESS_KEY") || ""),
        sessionToken: (val.sessionToken) ?? (utilities.getEnv("AWS_SESSION_TOKEN") || ""),
    };
}

export interface Context {
    /**
     * Resources to use for build context.
     *
     * The location can be:
     * * A relative or absolute path to a local directory (`.`, `./app`,
     *   `/app`, etc.).
     * * A remote URL of a Git repository, tarball, or plain text file
     *   (`https://github.com/user/myrepo.git`, `http://server/context.tar.gz`,
     *   etc.).
     */
    location: string;
}

export interface Dockerfile {
    /**
     * Raw Dockerfile contents.
     *
     * Conflicts with `location`.
     *
     * Equivalent to invoking Docker with `-f -`.
     */
    inline?: string;
    /**
     * Location of the Dockerfile to use.
     *
     * Can be a relative or absolute path to a local file, or a remote URL.
     *
     * Defaults to `${context.location}/Dockerfile` if context is on-disk.
     *
     * Conflicts with `inline`.
     */
    location?: string;
}

export interface Export {
    /**
     * A no-op export. Helpful for silencing the 'no exports' warning if you
     * just want to populate caches.
     */
    cacheonly?: outputs.ExportCacheOnly;
    /**
     * When `true` this entry will be excluded. Defaults to `false`.
     */
    disabled?: boolean;
    /**
     * Export as a Docker image layout.
     */
    docker?: outputs.ExportDocker;
    /**
     * Outputs the build result into a container image format.
     */
    image?: outputs.ExportImage;
    /**
     * Export to a local directory as files and directories.
     */
    local?: outputs.ExportLocal;
    /**
     * Identical to the Docker exporter but uses OCI media types by default.
     */
    oci?: outputs.ExportOCI;
    /**
     * A raw string as you would provide it to the Docker CLI (e.g.,
     * `type=docker`)
     */
    raw?: string;
    /**
     * Identical to the Image exporter, but pushes by default.
     */
    registry?: outputs.ExportRegistry;
    /**
     * Export to a local directory as a tarball.
     */
    tar?: outputs.ExportTar;
}
/**
 * exportProvideDefaults sets the appropriate defaults for Export
 */
export function exportProvideDefaults(val: Export): Export {
    return {
        ...val,
        docker: (val.docker ? outputs.exportDockerProvideDefaults(val.docker) : undefined),
        image: (val.image ? outputs.exportImageProvideDefaults(val.image) : undefined),
        oci: (val.oci ? outputs.exportOCIProvideDefaults(val.oci) : undefined),
        registry: (val.registry ? outputs.exportRegistryProvideDefaults(val.registry) : undefined),
    };
}

export interface ExportCacheOnly {
}

export interface ExportDocker {
    /**
     * Attach an arbitrary key/value annotation to the image.
     */
    annotations?: {[key: string]: string};
    /**
     * The compression type to use.
     */
    compression?: enums.CompressionType;
    /**
     * Compression level from 0 to 22.
     */
    compressionLevel?: number;
    /**
     * The local export path.
     */
    dest?: string;
    /**
     * Forcefully apply compression.
     */
    forceCompression?: boolean;
    /**
     * Specify images names to export. This is overridden if tags are already specified.
     */
    names?: string[];
    /**
     * Use OCI media types in exporter manifests.
     */
    ociMediaTypes?: boolean;
    /**
     * Bundle the output into a tarball layout.
     */
    tar?: boolean;
}
/**
 * exportDockerProvideDefaults sets the appropriate defaults for ExportDocker
 */
export function exportDockerProvideDefaults(val: ExportDocker): ExportDocker {
    return {
        ...val,
        compression: (val.compression) ?? "gzip",
        compressionLevel: (val.compressionLevel) ?? 0,
        forceCompression: (val.forceCompression) ?? false,
        ociMediaTypes: (val.ociMediaTypes) ?? false,
        tar: (val.tar) ?? true,
    };
}

export interface ExportImage {
    /**
     * Attach an arbitrary key/value annotation to the image.
     */
    annotations?: {[key: string]: string};
    /**
     * The compression type to use.
     */
    compression?: enums.CompressionType;
    /**
     * Compression level from 0 to 22.
     */
    compressionLevel?: number;
    /**
     * Name image with `prefix@<digest>`, used for anonymous images.
     */
    danglingNamePrefix?: string;
    /**
     * Forcefully apply compression.
     */
    forceCompression?: boolean;
    /**
     * Allow pushing to an insecure registry.
     */
    insecure?: boolean;
    /**
     * Add additional canonical name (`name@<digest>`).
     */
    nameCanonical?: boolean;
    /**
     * Specify images names to export. This is overridden if tags are already specified.
     */
    names?: string[];
    /**
     * Use OCI media types in exporter manifests.
     */
    ociMediaTypes?: boolean;
    /**
     * Push after creating the image. Defaults to `false`.
     */
    push?: boolean;
    /**
     * Push image without name.
     */
    pushByDigest?: boolean;
    /**
     * Store resulting images to the worker's image store and ensure all of
     * its blobs are in the content store.
     *
     * Defaults to `true`.
     *
     * Ignored if the worker doesn't have image store (when using OCI workers,
     * for example).
     */
    store?: boolean;
    /**
     * Unpack image after creation (for use with containerd). Defaults to
     * `false`.
     */
    unpack?: boolean;
}
/**
 * exportImageProvideDefaults sets the appropriate defaults for ExportImage
 */
export function exportImageProvideDefaults(val: ExportImage): ExportImage {
    return {
        ...val,
        compression: (val.compression) ?? "gzip",
        compressionLevel: (val.compressionLevel) ?? 0,
        forceCompression: (val.forceCompression) ?? false,
        ociMediaTypes: (val.ociMediaTypes) ?? false,
        store: (val.store) ?? true,
    };
}

export interface ExportLocal {
    /**
     * Output path.
     */
    dest: string;
}

export interface ExportOCI {
    /**
     * Attach an arbitrary key/value annotation to the image.
     */
    annotations?: {[key: string]: string};
    /**
     * The compression type to use.
     */
    compression?: enums.CompressionType;
    /**
     * Compression level from 0 to 22.
     */
    compressionLevel?: number;
    /**
     * The local export path.
     */
    dest?: string;
    /**
     * Forcefully apply compression.
     */
    forceCompression?: boolean;
    /**
     * Specify images names to export. This is overridden if tags are already specified.
     */
    names?: string[];
    /**
     * Use OCI media types in exporter manifests.
     */
    ociMediaTypes?: boolean;
    /**
     * Bundle the output into a tarball layout.
     */
    tar?: boolean;
}
/**
 * exportOCIProvideDefaults sets the appropriate defaults for ExportOCI
 */
export function exportOCIProvideDefaults(val: ExportOCI): ExportOCI {
    return {
        ...val,
        compression: (val.compression) ?? "gzip",
        compressionLevel: (val.compressionLevel) ?? 0,
        forceCompression: (val.forceCompression) ?? false,
        ociMediaTypes: (val.ociMediaTypes) ?? true,
        tar: (val.tar) ?? true,
    };
}

export interface ExportRegistry {
    /**
     * Attach an arbitrary key/value annotation to the image.
     */
    annotations?: {[key: string]: string};
    /**
     * The compression type to use.
     */
    compression?: enums.CompressionType;
    /**
     * Compression level from 0 to 22.
     */
    compressionLevel?: number;
    /**
     * Name image with `prefix@<digest>`, used for anonymous images.
     */
    danglingNamePrefix?: string;
    /**
     * Forcefully apply compression.
     */
    forceCompression?: boolean;
    /**
     * Allow pushing to an insecure registry.
     */
    insecure?: boolean;
    /**
     * Add additional canonical name (`name@<digest>`).
     */
    nameCanonical?: boolean;
    /**
     * Specify images names to export. This is overridden if tags are already specified.
     */
    names?: string[];
    /**
     * Use OCI media types in exporter manifests.
     */
    ociMediaTypes?: boolean;
    /**
     * Push after creating the image. Defaults to `true`.
     */
    push?: boolean;
    /**
     * Push image without name.
     */
    pushByDigest?: boolean;
    /**
     * Store resulting images to the worker's image store and ensure all of
     * its blobs are in the content store.
     *
     * Defaults to `true`.
     *
     * Ignored if the worker doesn't have image store (when using OCI workers,
     * for example).
     */
    store?: boolean;
    /**
     * Unpack image after creation (for use with containerd). Defaults to
     * `false`.
     */
    unpack?: boolean;
}
/**
 * exportRegistryProvideDefaults sets the appropriate defaults for ExportRegistry
 */
export function exportRegistryProvideDefaults(val: ExportRegistry): ExportRegistry {
    return {
        ...val,
        compression: (val.compression) ?? "gzip",
        compressionLevel: (val.compressionLevel) ?? 0,
        forceCompression: (val.forceCompression) ?? false,
        ociMediaTypes: (val.ociMediaTypes) ?? false,
        push: (val.push) ?? true,
        store: (val.store) ?? true,
    };
}

export interface ExportTar {
    /**
     * Output path.
     */
    dest: string;
}

export interface Registry {
    /**
     * The registry's address (e.g. "docker.io").
     */
    address: string;
    /**
     * Password or token for the registry.
     */
    password?: string;
    /**
     * Username for the registry.
     */
    username?: string;
}

export interface SSH {
    /**
     * Useful for distinguishing different servers that are part of the same
     * build.
     *
     * A value of `default` is appropriate if only dealing with a single host.
     */
    id: string;
    /**
     * SSH agent socket or private keys to expose to the build under the given
     * identifier.
     *
     * Defaults to `[$SSH_AUTH_SOCK]`.
     *
     * Note that your keys are **not** automatically added when using an
     * agent. Run `ssh-add -l` locally to confirm which public keys are
     * visible to the agent; these will be exposed to your build.
     */
    paths?: string[];
}

