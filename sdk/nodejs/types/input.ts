// *** WARNING: this file was generated by pulumi-language-nodejs. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as inputs from "../types/input";
import * as outputs from "../types/output";
import * as enums from "../types/enums";

import * as utilities from "../utilities";

export interface BuildContextArgs {
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
    location: pulumi.Input<string>;
    /**
     * Additional build contexts to use.
     *
     * These contexts are accessed with `FROM name` or `--from=name`
     * statements when using Dockerfile 1.4+ syntax.
     *
     * Values can be local paths, HTTP URLs, or  `docker-image://` images.
     */
    named?: pulumi.Input<{[key: string]: pulumi.Input<inputs.ContextArgs>}>;
}

export interface BuilderConfigArgs {
    /**
     * Name of an existing buildx builder to use.
     *
     * Only `docker-container`, `kubernetes`, or `remote` drivers are
     * supported. The legacy `docker` driver is not supported.
     *
     * Equivalent to Docker's `--builder` flag.
     */
    name?: pulumi.Input<string>;
}

export interface CacheFromArgs {
    /**
     * Upload build caches to Azure's blob storage service.
     */
    azblob?: pulumi.Input<inputs.CacheFromAzureBlobArgs>;
    /**
     * When `true` this entry will be excluded. Defaults to `false`.
     */
    disabled?: pulumi.Input<boolean>;
    /**
     * Recommended for use with GitHub Actions workflows.
     *
     * An action like `crazy-max/ghaction-github-runtime` is recommended to
     * expose appropriate credentials to your GitHub workflow.
     */
    gha?: pulumi.Input<inputs.CacheFromGitHubActionsArgs>;
    /**
     * A simple backend which caches images on your local filesystem.
     */
    local?: pulumi.Input<inputs.CacheFromLocalArgs>;
    /**
     * A raw string as you would provide it to the Docker CLI (e.g.,
     * `type=inline`).
     */
    raw?: pulumi.Input<string>;
    /**
     * Upload build caches to remote registries.
     */
    registry?: pulumi.Input<inputs.CacheFromRegistryArgs>;
    /**
     * Upload build caches to AWS S3 or an S3-compatible services such as
     * MinIO.
     */
    s3?: pulumi.Input<inputs.CacheFromS3Args>;
}
/**
 * cacheFromArgsProvideDefaults sets the appropriate defaults for CacheFromArgs
 */
export function cacheFromArgsProvideDefaults(val: CacheFromArgs): CacheFromArgs {
    return {
        ...val,
        gha: (val.gha ? pulumi.output(val.gha).apply(inputs.cacheFromGitHubActionsArgsProvideDefaults) : undefined),
        s3: (val.s3 ? pulumi.output(val.s3).apply(inputs.cacheFromS3ArgsProvideDefaults) : undefined),
    };
}

export interface CacheFromAzureBlobArgs {
    /**
     * Base URL of the storage account.
     */
    accountUrl?: pulumi.Input<string>;
    /**
     * The name of the cache image.
     */
    name: pulumi.Input<string>;
    /**
     * Blob storage account key.
     */
    secretAccessKey?: pulumi.Input<string>;
}

export interface CacheFromGitHubActionsArgs {
    /**
     * The scope to use for cache keys. Defaults to `buildkit`.
     *
     * This should be set if building and caching multiple images in one
     * workflow, otherwise caches will overwrite each other.
     */
    scope?: pulumi.Input<string>;
    /**
     * The GitHub Actions token to use. This is not a personal access tokens
     * and is typically generated automatically as part of each job.
     *
     * Defaults to `$ACTIONS_RUNTIME_TOKEN`, although a separate action like
     * `crazy-max/ghaction-github-runtime` is recommended to expose this
     * environment variable to your jobs.
     */
    token?: pulumi.Input<string>;
    /**
     * The cache server URL to use for artifacts.
     *
     * Defaults to `$ACTIONS_RUNTIME_URL`, although a separate action like
     * `crazy-max/ghaction-github-runtime` is recommended to expose this
     * environment variable to your jobs.
     */
    url?: pulumi.Input<string>;
}
/**
 * cacheFromGitHubActionsArgsProvideDefaults sets the appropriate defaults for CacheFromGitHubActionsArgs
 */
export function cacheFromGitHubActionsArgsProvideDefaults(val: CacheFromGitHubActionsArgs): CacheFromGitHubActionsArgs {
    return {
        ...val,
        scope: (val.scope) ?? (utilities.getEnv("buildkit") || ""),
        token: (val.token) ?? (utilities.getEnv("ACTIONS_RUNTIME_TOKEN") || ""),
        url: (val.url) ?? (utilities.getEnv("ACTIONS_RUNTIME_URL") || ""),
    };
}

export interface CacheFromLocalArgs {
    /**
     * Digest of manifest to import.
     */
    digest?: pulumi.Input<string>;
    /**
     * Path of the local directory where cache gets imported from.
     */
    src: pulumi.Input<string>;
}

export interface CacheFromRegistryArgs {
    /**
     * Fully qualified name of the cache image to import.
     */
    ref: pulumi.Input<string>;
}

export interface CacheFromS3Args {
    /**
     * Defaults to `$AWS_ACCESS_KEY_ID`.
     */
    accessKeyId?: pulumi.Input<string>;
    /**
     * Prefix to prepend to blob filenames.
     */
    blobsPrefix?: pulumi.Input<string>;
    /**
     * Name of the S3 bucket.
     */
    bucket: pulumi.Input<string>;
    /**
     * Endpoint of the S3 bucket.
     */
    endpointUrl?: pulumi.Input<string>;
    /**
     * Prefix to prepend on manifest filenames.
     */
    manifestsPrefix?: pulumi.Input<string>;
    /**
     * Name of the cache image.
     */
    name?: pulumi.Input<string>;
    /**
     * The geographic location of the bucket. Defaults to `$AWS_REGION`.
     */
    region: pulumi.Input<string>;
    /**
     * Defaults to `$AWS_SECRET_ACCESS_KEY`.
     */
    secretAccessKey?: pulumi.Input<string>;
    /**
     * Defaults to `$AWS_SESSION_TOKEN`.
     */
    sessionToken?: pulumi.Input<string>;
    /**
     * Uses `bucket` in the URL instead of hostname when `true`.
     */
    usePathStyle?: pulumi.Input<boolean>;
}
/**
 * cacheFromS3ArgsProvideDefaults sets the appropriate defaults for CacheFromS3Args
 */
export function cacheFromS3ArgsProvideDefaults(val: CacheFromS3Args): CacheFromS3Args {
    return {
        ...val,
        accessKeyId: (val.accessKeyId) ?? (utilities.getEnv("AWS_ACCESS_KEY_ID") || ""),
        region: (val.region) ?? (utilities.getEnv("AWS_REGION") || ""),
        secretAccessKey: (val.secretAccessKey) ?? (utilities.getEnv("AWS_SECRET_ACCESS_KEY") || ""),
        sessionToken: (val.sessionToken) ?? (utilities.getEnv("AWS_SESSION_TOKEN") || ""),
    };
}

export interface CacheToArgs {
    /**
     * Push cache to Azure's blob storage service.
     */
    azblob?: pulumi.Input<inputs.CacheToAzureBlobArgs>;
    /**
     * When `true` this entry will be excluded. Defaults to `false`.
     */
    disabled?: pulumi.Input<boolean>;
    /**
     * Recommended for use with GitHub Actions workflows.
     *
     * An action like `crazy-max/ghaction-github-runtime` is recommended to
     * expose appropriate credentials to your GitHub workflow.
     */
    gha?: pulumi.Input<inputs.CacheToGitHubActionsArgs>;
    /**
     * The inline cache storage backend is the simplest implementation to get
     * started with, but it does not handle multi-stage builds. Consider the
     * `registry` cache backend instead.
     */
    inline?: pulumi.Input<inputs.CacheToInlineArgs>;
    /**
     * A simple backend which caches imagines on your local filesystem.
     */
    local?: pulumi.Input<inputs.CacheToLocalArgs>;
    /**
     * A raw string as you would provide it to the Docker CLI (e.g.,
     * `type=inline`)
     */
    raw?: pulumi.Input<string>;
    /**
     * Push caches to remote registries. Incompatible with the `docker` build
     * driver.
     */
    registry?: pulumi.Input<inputs.CacheToRegistryArgs>;
    /**
     * Push cache to AWS S3 or S3-compatible services such as MinIO.
     */
    s3?: pulumi.Input<inputs.CacheToS3Args>;
}
/**
 * cacheToArgsProvideDefaults sets the appropriate defaults for CacheToArgs
 */
export function cacheToArgsProvideDefaults(val: CacheToArgs): CacheToArgs {
    return {
        ...val,
        azblob: (val.azblob ? pulumi.output(val.azblob).apply(inputs.cacheToAzureBlobArgsProvideDefaults) : undefined),
        gha: (val.gha ? pulumi.output(val.gha).apply(inputs.cacheToGitHubActionsArgsProvideDefaults) : undefined),
        local: (val.local ? pulumi.output(val.local).apply(inputs.cacheToLocalArgsProvideDefaults) : undefined),
        registry: (val.registry ? pulumi.output(val.registry).apply(inputs.cacheToRegistryArgsProvideDefaults) : undefined),
        s3: (val.s3 ? pulumi.output(val.s3).apply(inputs.cacheToS3ArgsProvideDefaults) : undefined),
    };
}

export interface CacheToAzureBlobArgs {
    /**
     * Base URL of the storage account.
     */
    accountUrl?: pulumi.Input<string>;
    /**
     * Ignore errors caused by failed cache exports.
     */
    ignoreError?: pulumi.Input<boolean>;
    /**
     * The cache mode to use. Defaults to `min`.
     */
    mode?: pulumi.Input<enums.CacheMode>;
    /**
     * The name of the cache image.
     */
    name: pulumi.Input<string>;
    /**
     * Blob storage account key.
     */
    secretAccessKey?: pulumi.Input<string>;
}
/**
 * cacheToAzureBlobArgsProvideDefaults sets the appropriate defaults for CacheToAzureBlobArgs
 */
export function cacheToAzureBlobArgsProvideDefaults(val: CacheToAzureBlobArgs): CacheToAzureBlobArgs {
    return {
        ...val,
        ignoreError: (val.ignoreError) ?? false,
        mode: (val.mode) ?? "min",
    };
}

export interface CacheToGitHubActionsArgs {
    /**
     * Ignore errors caused by failed cache exports.
     */
    ignoreError?: pulumi.Input<boolean>;
    /**
     * The cache mode to use. Defaults to `min`.
     */
    mode?: pulumi.Input<enums.CacheMode>;
    /**
     * The scope to use for cache keys. Defaults to `buildkit`.
     *
     * This should be set if building and caching multiple images in one
     * workflow, otherwise caches will overwrite each other.
     */
    scope?: pulumi.Input<string>;
    /**
     * The GitHub Actions token to use. This is not a personal access tokens
     * and is typically generated automatically as part of each job.
     *
     * Defaults to `$ACTIONS_RUNTIME_TOKEN`, although a separate action like
     * `crazy-max/ghaction-github-runtime` is recommended to expose this
     * environment variable to your jobs.
     */
    token?: pulumi.Input<string>;
    /**
     * The cache server URL to use for artifacts.
     *
     * Defaults to `$ACTIONS_RUNTIME_URL`, although a separate action like
     * `crazy-max/ghaction-github-runtime` is recommended to expose this
     * environment variable to your jobs.
     */
    url?: pulumi.Input<string>;
}
/**
 * cacheToGitHubActionsArgsProvideDefaults sets the appropriate defaults for CacheToGitHubActionsArgs
 */
export function cacheToGitHubActionsArgsProvideDefaults(val: CacheToGitHubActionsArgs): CacheToGitHubActionsArgs {
    return {
        ...val,
        ignoreError: (val.ignoreError) ?? false,
        mode: (val.mode) ?? "min",
        scope: (val.scope) ?? (utilities.getEnv("buildkit") || ""),
        token: (val.token) ?? (utilities.getEnv("ACTIONS_RUNTIME_TOKEN") || ""),
        url: (val.url) ?? (utilities.getEnv("ACTIONS_RUNTIME_URL") || ""),
    };
}

/**
 * Include an inline cache with the exported image.
 */
export interface CacheToInlineArgs {
}

export interface CacheToLocalArgs {
    /**
     * The compression type to use.
     */
    compression?: pulumi.Input<enums.CompressionType>;
    /**
     * Compression level from 0 to 22.
     */
    compressionLevel?: pulumi.Input<number>;
    /**
     * Path of the local directory to export the cache.
     */
    dest: pulumi.Input<string>;
    /**
     * Forcefully apply compression.
     */
    forceCompression?: pulumi.Input<boolean>;
    /**
     * Ignore errors caused by failed cache exports.
     */
    ignoreError?: pulumi.Input<boolean>;
    /**
     * The cache mode to use. Defaults to `min`.
     */
    mode?: pulumi.Input<enums.CacheMode>;
}
/**
 * cacheToLocalArgsProvideDefaults sets the appropriate defaults for CacheToLocalArgs
 */
export function cacheToLocalArgsProvideDefaults(val: CacheToLocalArgs): CacheToLocalArgs {
    return {
        ...val,
        compression: (val.compression) ?? "gzip",
        compressionLevel: (val.compressionLevel) ?? 0,
        forceCompression: (val.forceCompression) ?? false,
        ignoreError: (val.ignoreError) ?? false,
        mode: (val.mode) ?? "min",
    };
}

export interface CacheToRegistryArgs {
    /**
     * The compression type to use.
     */
    compression?: pulumi.Input<enums.CompressionType>;
    /**
     * Compression level from 0 to 22.
     */
    compressionLevel?: pulumi.Input<number>;
    /**
     * Forcefully apply compression.
     */
    forceCompression?: pulumi.Input<boolean>;
    /**
     * Ignore errors caused by failed cache exports.
     */
    ignoreError?: pulumi.Input<boolean>;
    /**
     * Export cache manifest as an OCI-compatible image manifest instead of a
     * manifest list. Requires `ociMediaTypes` to also be `true`.
     *
     * Some registries like AWS ECR will not work with caching if this is
     * `false`.
     *
     * Defaults to `false` to match Docker's default behavior.
     */
    imageManifest?: pulumi.Input<boolean>;
    /**
     * The cache mode to use. Defaults to `min`.
     */
    mode?: pulumi.Input<enums.CacheMode>;
    /**
     * Whether to use OCI media types in exported manifests. Defaults to
     * `true`.
     */
    ociMediaTypes?: pulumi.Input<boolean>;
    /**
     * Fully qualified name of the cache image to import.
     */
    ref: pulumi.Input<string>;
}
/**
 * cacheToRegistryArgsProvideDefaults sets the appropriate defaults for CacheToRegistryArgs
 */
export function cacheToRegistryArgsProvideDefaults(val: CacheToRegistryArgs): CacheToRegistryArgs {
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

export interface CacheToS3Args {
    /**
     * Defaults to `$AWS_ACCESS_KEY_ID`.
     */
    accessKeyId?: pulumi.Input<string>;
    /**
     * Prefix to prepend to blob filenames.
     */
    blobsPrefix?: pulumi.Input<string>;
    /**
     * Name of the S3 bucket.
     */
    bucket: pulumi.Input<string>;
    /**
     * Endpoint of the S3 bucket.
     */
    endpointUrl?: pulumi.Input<string>;
    /**
     * Ignore errors caused by failed cache exports.
     */
    ignoreError?: pulumi.Input<boolean>;
    /**
     * Prefix to prepend on manifest filenames.
     */
    manifestsPrefix?: pulumi.Input<string>;
    /**
     * The cache mode to use. Defaults to `min`.
     */
    mode?: pulumi.Input<enums.CacheMode>;
    /**
     * Name of the cache image.
     */
    name?: pulumi.Input<string>;
    /**
     * The geographic location of the bucket. Defaults to `$AWS_REGION`.
     */
    region: pulumi.Input<string>;
    /**
     * Defaults to `$AWS_SECRET_ACCESS_KEY`.
     */
    secretAccessKey?: pulumi.Input<string>;
    /**
     * Defaults to `$AWS_SESSION_TOKEN`.
     */
    sessionToken?: pulumi.Input<string>;
    /**
     * Uses `bucket` in the URL instead of hostname when `true`.
     */
    usePathStyle?: pulumi.Input<boolean>;
}
/**
 * cacheToS3ArgsProvideDefaults sets the appropriate defaults for CacheToS3Args
 */
export function cacheToS3ArgsProvideDefaults(val: CacheToS3Args): CacheToS3Args {
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

export interface ContextArgs {
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
    location: pulumi.Input<string>;
}

export interface DockerfileArgs {
    /**
     * Raw Dockerfile contents.
     *
     * Conflicts with `location`.
     *
     * Equivalent to invoking Docker with `-f -`.
     */
    inline?: pulumi.Input<string>;
    /**
     * Location of the Dockerfile to use.
     *
     * Can be a relative or absolute path to a local file, or a remote URL.
     *
     * Defaults to `${context.location}/Dockerfile` if context is on-disk.
     *
     * Conflicts with `inline`.
     */
    location?: pulumi.Input<string>;
}

export interface ExportArgs {
    /**
     * A no-op export. Helpful for silencing the 'no exports' warning if you
     * just want to populate caches.
     */
    cacheonly?: pulumi.Input<inputs.ExportCacheOnlyArgs>;
    /**
     * When `true` this entry will be excluded. Defaults to `false`.
     */
    disabled?: pulumi.Input<boolean>;
    /**
     * Export as a Docker image layout.
     */
    docker?: pulumi.Input<inputs.ExportDockerArgs>;
    /**
     * Outputs the build result into a container image format.
     */
    image?: pulumi.Input<inputs.ExportImageArgs>;
    /**
     * Export to a local directory as files and directories.
     */
    local?: pulumi.Input<inputs.ExportLocalArgs>;
    /**
     * Identical to the Docker exporter but uses OCI media types by default.
     */
    oci?: pulumi.Input<inputs.ExportOCIArgs>;
    /**
     * A raw string as you would provide it to the Docker CLI (e.g.,
     * `type=docker`)
     */
    raw?: pulumi.Input<string>;
    /**
     * Identical to the Image exporter, but pushes by default.
     */
    registry?: pulumi.Input<inputs.ExportRegistryArgs>;
    /**
     * Export to a local directory as a tarball.
     */
    tar?: pulumi.Input<inputs.ExportTarArgs>;
}
/**
 * exportArgsProvideDefaults sets the appropriate defaults for ExportArgs
 */
export function exportArgsProvideDefaults(val: ExportArgs): ExportArgs {
    return {
        ...val,
        docker: (val.docker ? pulumi.output(val.docker).apply(inputs.exportDockerArgsProvideDefaults) : undefined),
        image: (val.image ? pulumi.output(val.image).apply(inputs.exportImageArgsProvideDefaults) : undefined),
        oci: (val.oci ? pulumi.output(val.oci).apply(inputs.exportOCIArgsProvideDefaults) : undefined),
        registry: (val.registry ? pulumi.output(val.registry).apply(inputs.exportRegistryArgsProvideDefaults) : undefined),
    };
}

export interface ExportCacheOnlyArgs {
}

export interface ExportDockerArgs {
    /**
     * Attach an arbitrary key/value annotation to the image.
     */
    annotations?: pulumi.Input<{[key: string]: pulumi.Input<string>}>;
    /**
     * The compression type to use.
     */
    compression?: pulumi.Input<enums.CompressionType>;
    /**
     * Compression level from 0 to 22.
     */
    compressionLevel?: pulumi.Input<number>;
    /**
     * The local export path.
     */
    dest?: pulumi.Input<string>;
    /**
     * Forcefully apply compression.
     */
    forceCompression?: pulumi.Input<boolean>;
    /**
     * Specify images names to export. This is overridden if tags are already specified.
     */
    names?: pulumi.Input<pulumi.Input<string>[]>;
    /**
     * Use OCI media types in exporter manifests.
     */
    ociMediaTypes?: pulumi.Input<boolean>;
    /**
     * Bundle the output into a tarball layout.
     */
    tar?: pulumi.Input<boolean>;
}
/**
 * exportDockerArgsProvideDefaults sets the appropriate defaults for ExportDockerArgs
 */
export function exportDockerArgsProvideDefaults(val: ExportDockerArgs): ExportDockerArgs {
    return {
        ...val,
        compression: (val.compression) ?? "gzip",
        compressionLevel: (val.compressionLevel) ?? 0,
        forceCompression: (val.forceCompression) ?? false,
        ociMediaTypes: (val.ociMediaTypes) ?? false,
        tar: (val.tar) ?? true,
    };
}

export interface ExportImageArgs {
    /**
     * Attach an arbitrary key/value annotation to the image.
     */
    annotations?: pulumi.Input<{[key: string]: pulumi.Input<string>}>;
    /**
     * The compression type to use.
     */
    compression?: pulumi.Input<enums.CompressionType>;
    /**
     * Compression level from 0 to 22.
     */
    compressionLevel?: pulumi.Input<number>;
    /**
     * Name image with `prefix@<digest>`, used for anonymous images.
     */
    danglingNamePrefix?: pulumi.Input<string>;
    /**
     * Forcefully apply compression.
     */
    forceCompression?: pulumi.Input<boolean>;
    /**
     * Allow pushing to an insecure registry.
     */
    insecure?: pulumi.Input<boolean>;
    /**
     * Add additional canonical name (`name@<digest>`).
     */
    nameCanonical?: pulumi.Input<boolean>;
    /**
     * Specify images names to export. This is overridden if tags are already specified.
     */
    names?: pulumi.Input<pulumi.Input<string>[]>;
    /**
     * Use OCI media types in exporter manifests.
     */
    ociMediaTypes?: pulumi.Input<boolean>;
    /**
     * Push after creating the image. Defaults to `false`.
     */
    push?: pulumi.Input<boolean>;
    /**
     * Push image without name.
     */
    pushByDigest?: pulumi.Input<boolean>;
    /**
     * Store resulting images to the worker's image store and ensure all of
     * its blobs are in the content store.
     *
     * Defaults to `true`.
     *
     * Ignored if the worker doesn't have image store (when using OCI workers,
     * for example).
     */
    store?: pulumi.Input<boolean>;
    /**
     * Unpack image after creation (for use with containerd). Defaults to
     * `false`.
     */
    unpack?: pulumi.Input<boolean>;
}
/**
 * exportImageArgsProvideDefaults sets the appropriate defaults for ExportImageArgs
 */
export function exportImageArgsProvideDefaults(val: ExportImageArgs): ExportImageArgs {
    return {
        ...val,
        compression: (val.compression) ?? "gzip",
        compressionLevel: (val.compressionLevel) ?? 0,
        forceCompression: (val.forceCompression) ?? false,
        ociMediaTypes: (val.ociMediaTypes) ?? false,
        store: (val.store) ?? true,
    };
}

export interface ExportLocalArgs {
    /**
     * Output path.
     */
    dest: pulumi.Input<string>;
}

export interface ExportOCIArgs {
    /**
     * Attach an arbitrary key/value annotation to the image.
     */
    annotations?: pulumi.Input<{[key: string]: pulumi.Input<string>}>;
    /**
     * The compression type to use.
     */
    compression?: pulumi.Input<enums.CompressionType>;
    /**
     * Compression level from 0 to 22.
     */
    compressionLevel?: pulumi.Input<number>;
    /**
     * The local export path.
     */
    dest?: pulumi.Input<string>;
    /**
     * Forcefully apply compression.
     */
    forceCompression?: pulumi.Input<boolean>;
    /**
     * Specify images names to export. This is overridden if tags are already specified.
     */
    names?: pulumi.Input<pulumi.Input<string>[]>;
    /**
     * Use OCI media types in exporter manifests.
     */
    ociMediaTypes?: pulumi.Input<boolean>;
    /**
     * Bundle the output into a tarball layout.
     */
    tar?: pulumi.Input<boolean>;
}
/**
 * exportOCIArgsProvideDefaults sets the appropriate defaults for ExportOCIArgs
 */
export function exportOCIArgsProvideDefaults(val: ExportOCIArgs): ExportOCIArgs {
    return {
        ...val,
        compression: (val.compression) ?? "gzip",
        compressionLevel: (val.compressionLevel) ?? 0,
        forceCompression: (val.forceCompression) ?? false,
        ociMediaTypes: (val.ociMediaTypes) ?? true,
        tar: (val.tar) ?? true,
    };
}

export interface ExportRegistryArgs {
    /**
     * Attach an arbitrary key/value annotation to the image.
     */
    annotations?: pulumi.Input<{[key: string]: pulumi.Input<string>}>;
    /**
     * The compression type to use.
     */
    compression?: pulumi.Input<enums.CompressionType>;
    /**
     * Compression level from 0 to 22.
     */
    compressionLevel?: pulumi.Input<number>;
    /**
     * Name image with `prefix@<digest>`, used for anonymous images.
     */
    danglingNamePrefix?: pulumi.Input<string>;
    /**
     * Forcefully apply compression.
     */
    forceCompression?: pulumi.Input<boolean>;
    /**
     * Allow pushing to an insecure registry.
     */
    insecure?: pulumi.Input<boolean>;
    /**
     * Add additional canonical name (`name@<digest>`).
     */
    nameCanonical?: pulumi.Input<boolean>;
    /**
     * Specify images names to export. This is overridden if tags are already specified.
     */
    names?: pulumi.Input<pulumi.Input<string>[]>;
    /**
     * Use OCI media types in exporter manifests.
     */
    ociMediaTypes?: pulumi.Input<boolean>;
    /**
     * Push after creating the image. Defaults to `true`.
     */
    push?: pulumi.Input<boolean>;
    /**
     * Push image without name.
     */
    pushByDigest?: pulumi.Input<boolean>;
    /**
     * Store resulting images to the worker's image store and ensure all of
     * its blobs are in the content store.
     *
     * Defaults to `true`.
     *
     * Ignored if the worker doesn't have image store (when using OCI workers,
     * for example).
     */
    store?: pulumi.Input<boolean>;
    /**
     * Unpack image after creation (for use with containerd). Defaults to
     * `false`.
     */
    unpack?: pulumi.Input<boolean>;
}
/**
 * exportRegistryArgsProvideDefaults sets the appropriate defaults for ExportRegistryArgs
 */
export function exportRegistryArgsProvideDefaults(val: ExportRegistryArgs): ExportRegistryArgs {
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

export interface ExportTarArgs {
    /**
     * Output path.
     */
    dest: pulumi.Input<string>;
}

export interface RegistryArgs {
    /**
     * The registry's address (e.g. "docker.io").
     */
    address: pulumi.Input<string>;
    /**
     * Password or token for the registry.
     */
    password?: pulumi.Input<string>;
    /**
     * Username for the registry.
     */
    username?: pulumi.Input<string>;
}

export interface SSHArgs {
    /**
     * Useful for distinguishing different servers that are part of the same
     * build.
     *
     * A value of `default` is appropriate if only dealing with a single host.
     */
    id: pulumi.Input<string>;
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
    paths?: pulumi.Input<pulumi.Input<string>[]>;
}
