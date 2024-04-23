// *** WARNING: this file was generated by pulumi-language-nodejs. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as inputs from "./types/input";
import * as outputs from "./types/output";
import * as enums from "./types/enums";
import * as utilities from "./utilities";

/**
 * A Docker image built using buildx -- Docker's interface to the improved
 * BuildKit backend.
 *
 * ## Stability
 *
 * **This resource is pre-1.0 and in public preview.**
 *
 * We will strive to keep APIs and behavior as stable as possible, but we
 * cannot guarantee stability until version 1.0.
 *
 * ## Migrating Pulumi Docker v3 and v4 Image resources
 *
 * This provider's `Image` resource provides a superset of functionality over the `Image` resources available in versions 3 and 4 of the Pulumi Docker provider.
 * Existing `Image` resources can be converted to the docker-build `Image` resources with minor modifications.
 *
 * ### Behavioral differences
 *
 * There are several key behavioral differences to keep in mind when transitioning images to the new `Image` resource.
 *
 * #### Previews
 *
 * Version `3.x` of the Pulumi Docker provider always builds images during preview operations.
 * This is helpful as a safeguard to prevent "broken" images from merging, but users found the behavior unnecessarily redundant when running previews and updates locally.
 *
 * Version `4.x` changed build-on-preview behavior to be opt-in.
 * By default, `v4.x` `Image` resources do _not_ build during previews, but this behavior can be toggled with the `buildOnPreview` option.
 * Several users reported outages due to the default behavior allowing bad images to accidentally sneak through CI.
 *
 * The default behavior of this provider's `Image` resource is similar to `3.x` and will build images during previews.
 * This behavior can be changed by specifying `buildOnPreview`.
 *
 * #### Push behavior
 *
 * Versions `3.x` and `4.x` of the Pulumi Docker provider attempt to push images to remote registries by default.
 * They expose a `skipPush: true` option to disable pushing.
 *
 * This provider's `Image` resource matches the Docker CLI's behavior and does not push images anywhere by default.
 *
 * To push images to a registry you can include `push: true` (equivalent to Docker's `--push` flag) or configure an `export` of type `registry` (equivalent to Docker's `--output type=registry`).
 * Like Docker, if an image is configured without exports you will see a warning with instructions for how to enable pushing, but the build will still proceed normally.
 *
 * #### Secrets
 *
 * Version `3.x` of the Pulumi Docker provider supports secrets by way of the `extraOptions` field.
 *
 * Version `4.x` of the Pulumi Docker provider does not support secrets.
 *
 * The `Image` resource supports secrets but does not require those secrets to exist on-disk or in environment variables.
 * Instead, they should be passed directly as values.
 * (Please be sure to familiarize yourself with Pulumi's [native secret handling](https://www.pulumi.com/docs/concepts/secrets/).)
 * Pulumi also provides [ESC](https://www.pulumi.com/product/esc/) to make it easier to share secrets across stacks and environments.
 *
 * #### Caching
 *
 * Version `3.x` of the Pulumi Docker provider exposes `cacheFrom: bool | { stages: [...] }`.
 * It builds targets individually and pushes them to separate images for caching.
 *
 * Version `4.x` exposes a similar parameter `cacheFrom: { images: [...] }` which pushes and pulls inline caches.
 *
 * Both versions 3 and 4 require specific environment variables to be set and deviate from Docker's native caching behavior.
 * This can result in inefficient builds due to unnecessary image pulls, repeated file transfers, etc.
 *
 * The `Image` resource delegates all caching behavior to Docker.
 * `cacheFrom` and `cacheTo` options (equivalent to Docker's `--cache-to` and `--cache-from`) are exposed and provide additional cache targets, such as local disk, S3 storage, etc.
 *
 * #### Outputs
 *
 * Versions `3.x` and `4.x` of the provider exposed a `repoDigest` output which was a fully qualified tag with digest.
 * In `4.x` this could also be a single sha256 hash if the image wasn't pushed.
 *
 * Unlike earlier providers the `Image` resource can push multiple tags.
 * As a convenience, it exposes a `ref` output consisting of a tag with digest as long as the image was pushed.
 * If multiple tags were pushed this uses one at random.
 *
 * If you need more control over tag references you can use the `digest` output, which is always a single sha256 hash as long as the image was exported somewhere.
 *
 * #### Tag deletion and refreshes
 *
 * Versions 3 and 4 of Pulumi Docker provider do not delete tags when the `Image` resource is deleted, nor do they confirm expected tags exist during `refresh` operations.
 *
 * The `buidx.Image` will query your registries during `refresh` to ensure the expected tags exist.
 * If any are missing a subsequent `update` will push them.
 *
 * When a `Image` is deleted, it will _attempt_ to also delete any pushed tags.
 * Deletion of remote tags is not guaranteed because not all registries support the manifest `DELETE` API (`docker.io` in particular).
 * Manifests are _not_ deleted in the same way during updates -- to do so safely would require a full build to determine whether a Pulumi operation should be an update or update-replace.
 *
 * Use the [`retainOnDelete: true`](https://www.pulumi.com/docs/concepts/options/retainondelete/) option if you do not want tags deleted.
 *
 * ### Example migration
 *
 * Examples of "fully-featured" `v3` and `v4` `Image` resources are shown below, along with an example `Image` resource showing how they would look after migration.
 *
 * The `v3` resource leverages `buildx` via a `DOCKER_BUILDKIT` environment variable and CLI flags passed in with `extraOption`.
 * After migration, the environment variable is no longer needed and CLI flags are now properties on the `Image`.
 * In almost all cases, properties of `Image` are named after the Docker CLI flag they correspond to.
 *
 * The `v4` resource is less functional than its `v3` counterpart because it lacks the flexibility of `extraOptions`.
 * It it is shown with parameters similar to the `v3` example for completeness.
 *
 * ## Example Usage
 * ### v3/v4 migration
 *
 * ```typescript
 *
 * // v3 Image
 * const v3 = new docker.Image("v3-image", {
 *   imageName: "myregistry.com/user/repo:latest",
 *   localImageName: "local-tag",
 *   skipPush: false,
 *   build: {
 *     dockerfile: "./Dockerfile",
 *     context: "../app",
 *     target: "mytarget",
 *     args: {
 *       MY_BUILD_ARG: "foo",
 *     },
 *     env: {
 *       DOCKER_BUILDKIT: "1",
 *     },
 *     extraOptions: [
 *       "--cache-from",
 *       "type=registry,myregistry.com/user/repo:cache",
 *       "--cache-to",
 *       "type=registry,myregistry.com/user/repo:cache",
 *       "--add-host",
 *       "metadata.google.internal:169.254.169.254",
 *       "--secret",
 *       "id=mysecret,src=/local/secret",
 *       "--ssh",
 *       "default=/home/runner/.ssh/id_ed25519",
 *       "--network",
 *       "host",
 *       "--platform",
 *       "linux/amd64",
 *     ],
 *   },
 *   registry: {
 *     server: "myregistry.com",
 *     username: "username",
 *     password: pulumi.secret("password"),
 *   },
 * });
 *
 * // v3 Image after migrating to docker-build.Image
 * const v3Migrated = new dockerbuild.Image("v3-to-buildx", {
 *     tags: ["myregistry.com/user/repo:latest", "local-tag"],
 *     push: true,
 *     dockerfile: {
 *         location: "./Dockerfile",
 *     },
 *     context: {
 *         location: "../app",
 *     },
 *     target: "mytarget",
 *     buildArgs: {
 *         MY_BUILD_ARG: "foo",
 *     },
 *     cacheFrom: [{ registry: { ref: "myregistry.com/user/repo:cache" } }],
 *     cacheTo: [{ registry: { ref: "myregistry.com/user/repo:cache" } }],
 *     secrets: {
 *         mysecret: "value",
 *     },
 *     addHosts: ["metadata.google.internal:169.254.169.254"],
 *     ssh: {
 *         default: ["/home/runner/.ssh/id_ed25519"],
 *     },
 *     network: "host",
 *     platforms: ["linux/amd64"],
 *     registries: [{
 *         address: "myregistry.com",
 *         username: "username",
 *         password: pulumi.secret("password"),
 *     }],
 * });
 *
 *
 * // v4 Image
 * const v4 = new docker.Image("v4-image", {
 *     imageName: "myregistry.com/user/repo:latest",
 *     skipPush: false,
 *     build: {
 *         dockerfile: "./Dockerfile",
 *         context: "../app",
 *         target: "mytarget",
 *         args: {
 *             MY_BUILD_ARG: "foo",
 *         },
 *         cacheFrom: {
 *             images: ["myregistry.com/user/repo:cache"],
 *         },
 *         addHosts: ["metadata.google.internal:169.254.169.254"],
 *         network: "host",
 *         platform: "linux/amd64",
 *     },
 *     buildOnPreview: true,
 *     registry: {
 *         server: "myregistry.com",
 *         username: "username",
 *         password: pulumi.secret("password"),
 *     },
 * });
 *
 * // v4 Image after migrating to docker-build.Image
 * const v4Migrated = new dockerbuild.Image("v4-to-buildx", {
 *     tags: ["myregistry.com/user/repo:latest"],
 *     push: true,
 *     dockerfile: {
 *         location: "./Dockerfile",
 *     },
 *     context: {
 *         location: "../app",
 *     },
 *     target: "mytarget",
 *     buildArgs: {
 *         MY_BUILD_ARG: "foo",
 *     },
 *     cacheFrom: [{ registry: { ref: "myregistry.com/user/repo:cache" } }],
 *     cacheTo: [{ registry: { ref: "myregistry.com/user/repo:cache" } }],
 *     addHosts: ["metadata.google.internal:169.254.169.254"],
 *     network: "host",
 *     platforms: ["linux/amd64"],
 *     registries: [{
 *         address: "myregistry.com",
 *         username: "username",
 *         password: pulumi.secret("password"),
 *     }],
 * });
 *
 * ```
 *
 * ## Example Usage
 * ### Push to AWS ECR with caching
 *
 * ```typescript
 * import * as pulumi from "@pulumi/pulumi";
 * import * as aws from "@pulumi/aws";
 * import * as docker_build from "@pulumi/docker-build";
 *
 * const ecrRepository = new aws.ecr.Repository("ecr-repository", {});
 * const authToken = aws.ecr.getAuthorizationTokenOutput({
 *     registryId: ecrRepository.registryId,
 * });
 * const myImage = new docker_build.Image("my-image", {
 *     cacheFrom: [{
 *         registry: {
 *             ref: pulumi.interpolate`${ecrRepository.repositoryUrl}:cache`,
 *         },
 *     }],
 *     cacheTo: [{
 *         registry: {
 *             imageManifest: true,
 *             ociMediaTypes: true,
 *             ref: pulumi.interpolate`${ecrRepository.repositoryUrl}:cache`,
 *         },
 *     }],
 *     context: {
 *         location: "./app",
 *     },
 *     push: true,
 *     registries: [{
 *         address: ecrRepository.repositoryUrl,
 *         password: authToken.apply(authToken => authToken.password),
 *         username: authToken.apply(authToken => authToken.userName),
 *     }],
 *     tags: [pulumi.interpolate`${ecrRepository.repositoryUrl}:latest`],
 * });
 * export const ref = myImage.ref;
 * ```
 * ### Multi-platform image
 *
 * ```typescript
 * import * as pulumi from "@pulumi/pulumi";
 * import * as docker_build from "@pulumi/docker-build";
 *
 * const image = new docker_build.Image("image", {
 *     context: {
 *         location: "app",
 *     },
 *     platforms: [
 *         docker_build.Platform.Plan9_amd64,
 *         docker_build.Platform.Plan9_386,
 *     ],
 * });
 * ```
 * ### Registry export
 *
 * ```typescript
 * import * as pulumi from "@pulumi/pulumi";
 * import * as docker_build from "@pulumi/docker-build";
 *
 * const image = new docker_build.Image("image", {
 *     context: {
 *         location: "app",
 *     },
 *     push: true,
 *     registries: [{
 *         address: "docker.io",
 *         password: dockerHubPassword,
 *         username: "pulumibot",
 *     }],
 *     tags: ["docker.io/pulumi/pulumi:3.107.0"],
 * });
 * export const ref = myImage.ref;
 * ```
 * ### Caching
 *
 * ```typescript
 * import * as pulumi from "@pulumi/pulumi";
 * import * as docker_build from "@pulumi/docker-build";
 *
 * const image = new docker_build.Image("image", {
 *     cacheFrom: [{
 *         local: {
 *             src: "tmp/cache",
 *         },
 *     }],
 *     cacheTo: [{
 *         local: {
 *             dest: "tmp/cache",
 *             mode: docker_build.CacheMode.Max,
 *         },
 *     }],
 *     context: {
 *         location: "app",
 *     },
 * });
 * ```
 * ### Docker Build Cloud
 *
 * ```typescript
 * import * as pulumi from "@pulumi/pulumi";
 * import * as docker_build from "@pulumi/docker-build";
 *
 * const image = new docker_build.Image("image", {
 *     builder: {
 *         name: "cloud-builder-name",
 *     },
 *     context: {
 *         location: "app",
 *     },
 *     exec: true,
 * });
 * ```
 * ### Build arguments
 *
 * ```typescript
 * import * as pulumi from "@pulumi/pulumi";
 * import * as docker_build from "@pulumi/docker-build";
 *
 * const image = new docker_build.Image("image", {
 *     buildArgs: {
 *         SET_ME_TO_TRUE: "true",
 *     },
 *     context: {
 *         location: "app",
 *     },
 * });
 * ```
 * ### Build target
 *
 * ```typescript
 * import * as pulumi from "@pulumi/pulumi";
 * import * as docker_build from "@pulumi/docker-build";
 *
 * const image = new docker_build.Image("image", {
 *     context: {
 *         location: "app",
 *     },
 *     target: "build-me",
 * });
 * ```
 * ### Named contexts
 *
 * ```typescript
 * import * as pulumi from "@pulumi/pulumi";
 * import * as docker_build from "@pulumi/docker-build";
 *
 * const image = new docker_build.Image("image", {context: {
 *     location: "app",
 *     named: {
 *         "golang:latest": {
 *             location: "docker-image://golang@sha256:b8e62cf593cdaff36efd90aa3a37de268e6781a2e68c6610940c48f7cdf36984",
 *         },
 *     },
 * }});
 * ```
 * ### Remote context
 *
 * ```typescript
 * import * as pulumi from "@pulumi/pulumi";
 * import * as docker_build from "@pulumi/docker-build";
 *
 * const image = new docker_build.Image("image", {context: {
 *     location: "https://raw.githubusercontent.com/pulumi/pulumi-docker/api-types/provider/testdata/Dockerfile",
 * }});
 * ```
 * ### Inline Dockerfile
 *
 * ```typescript
 * import * as pulumi from "@pulumi/pulumi";
 * import * as docker_build from "@pulumi/docker-build";
 *
 * const image = new docker_build.Image("image", {
 *     context: {
 *         location: "app",
 *     },
 *     dockerfile: {
 *         inline: `FROM busybox
 * COPY hello.c ./
 * `,
 *     },
 * });
 * ```
 * ### Remote context
 *
 * ```typescript
 * import * as pulumi from "@pulumi/pulumi";
 * import * as docker_build from "@pulumi/docker-build";
 *
 * const image = new docker_build.Image("image", {
 *     context: {
 *         location: "https://github.com/docker-library/hello-world.git",
 *     },
 *     dockerfile: {
 *         location: "app/Dockerfile",
 *     },
 * });
 * ```
 * ### Local export
 *
 * ```typescript
 * import * as pulumi from "@pulumi/pulumi";
 * import * as docker_build from "@pulumi/docker-build";
 *
 * const image = new docker_build.Image("image", {
 *     context: {
 *         location: "app",
 *     },
 *     exports: [{
 *         docker: {
 *             tar: true,
 *         },
 *     }],
 * });
 * ```
 */
export class Image extends pulumi.CustomResource {
    /**
     * Get an existing Image resource's state with the given name, ID, and optional extra
     * properties used to qualify the lookup.
     *
     * @param name The _unique_ name of the resulting resource.
     * @param id The _unique_ provider ID of the resource to lookup.
     * @param opts Optional settings to control the behavior of the CustomResource.
     */
    public static get(name: string, id: pulumi.Input<pulumi.ID>, opts?: pulumi.CustomResourceOptions): Image {
        return new Image(name, undefined as any, { ...opts, id: id });
    }

    /** @internal */
    public static readonly __pulumiType = 'docker-build:index:Image';

    /**
     * Returns true if the given object is an instance of Image.  This is designed to work even
     * when multiple copies of the Pulumi SDK have been loaded into the same process.
     */
    public static isInstance(obj: any): obj is Image {
        if (obj === undefined || obj === null) {
            return false;
        }
        return obj['__pulumiType'] === Image.__pulumiType;
    }

    /**
     * Custom `host:ip` mappings to use during the build.
     *
     * Equivalent to Docker's `--add-host` flag.
     */
    public readonly addHosts!: pulumi.Output<string[] | undefined>;
    /**
     * `ARG` names and values to set during the build.
     *
     * These variables are accessed like environment variables inside `RUN`
     * instructions.
     *
     * Build arguments are persisted in the image, so you should use `secrets`
     * if these arguments are sensitive.
     *
     * Equivalent to Docker's `--build-arg` flag.
     */
    public readonly buildArgs!: pulumi.Output<{[key: string]: string} | undefined>;
    /**
     * Setting this to `false` will always skip image builds during previews,
     * and setting it to `true` will always build images during previews.
     *
     * Images built during previews are never exported to registries, however
     * cache manifests are still exported.
     *
     * On-disk Dockerfiles are always validated for syntactic correctness
     * regardless of this setting.
     *
     * Defaults to `true` as a safeguard against broken images merging as part
     * of CI pipelines.
     */
    public readonly buildOnPreview!: pulumi.Output<boolean | undefined>;
    /**
     * Builder configuration.
     */
    public readonly builder!: pulumi.Output<outputs.BuilderConfig | undefined>;
    /**
     * Cache export configuration.
     *
     * Equivalent to Docker's `--cache-from` flag.
     */
    public readonly cacheFrom!: pulumi.Output<outputs.CacheFrom[] | undefined>;
    /**
     * Cache import configuration.
     *
     * Equivalent to Docker's `--cache-to` flag.
     */
    public readonly cacheTo!: pulumi.Output<outputs.CacheTo[] | undefined>;
    /**
     * Build context settings.
     *
     * Equivalent to Docker's `PATH | URL | -` positional argument.
     */
    public readonly context!: pulumi.Output<outputs.BuildContext | undefined>;
    /**
     * A preliminary hash of the image's build context.
     *
     * Pulumi uses this to determine if an image _may_ need to be re-built.
     */
    public /*out*/ readonly contextHash!: pulumi.Output<string>;
    /**
     * A SHA256 digest of the image if it was exported to a registry or
     * elsewhere.
     *
     * Empty if the image was not exported.
     *
     * Registry images can be referenced precisely as `<tag>@<digest>`. The
     * `ref` output provides one such reference as a convenience.
     */
    public /*out*/ readonly digest!: pulumi.Output<string>;
    /**
     * Dockerfile settings.
     *
     * Equivalent to Docker's `--file` flag.
     */
    public readonly dockerfile!: pulumi.Output<outputs.Dockerfile | undefined>;
    /**
     * Use `exec` mode to build this image.
     *
     * By default the provider embeds a v25 Docker client with v0.12 buildx
     * support. This helps ensure consistent behavior across environments and
     * is compatible with alternative build backends (e.g. `buildkitd`), but
     * it may not be desirable if you require a specific version of buildx.
     * For example you may want to run a custom `docker-buildx` binary with
     * support for [Docker Build
     * Cloud](https://docs.docker.com/build/cloud/setup/) (DBC).
     *
     * When this is set to `true` the provider will instead execute the
     * `docker-buildx` binary directly to perform its operations. The user is
     * responsible for ensuring this binary exists, with correct permissions
     * and pre-configured builders, at a path Docker expects (e.g.
     * `~/.docker/cli-plugins`).
     *
     * Debugging `exec` mode may be more difficult as Pulumi will not be able
     * to surface fine-grained errors and warnings. Additionally credentials
     * are temporarily written to disk in order to provide them to the
     * `docker-buildx` binary.
     */
    public readonly exec!: pulumi.Output<boolean | undefined>;
    /**
     * Controls where images are persisted after building.
     *
     * Images are only stored in the local cache unless `exports` are
     * explicitly configured.
     *
     * Exporting to multiple destinations requires a daemon running BuildKit
     * 0.13 or later.
     *
     * Equivalent to Docker's `--output` flag.
     */
    public readonly exports!: pulumi.Output<outputs.Export[] | undefined>;
    /**
     * Attach arbitrary key/value metadata to the image.
     *
     * Equivalent to Docker's `--label` flag.
     */
    public readonly labels!: pulumi.Output<{[key: string]: string} | undefined>;
    /**
     * When `true` the build will automatically include a `docker` export.
     *
     * Defaults to `false`.
     *
     * Equivalent to Docker's `--load` flag.
     */
    public readonly load!: pulumi.Output<boolean | undefined>;
    /**
     * Set the network mode for `RUN` instructions. Defaults to `default`.
     *
     * For custom networks, configure your builder with `--driver-opt network=...`.
     *
     * Equivalent to Docker's `--network` flag.
     */
    public readonly network!: pulumi.Output<enums.NetworkMode | undefined>;
    /**
     * Do not import cache manifests when building the image.
     *
     * Equivalent to Docker's `--no-cache` flag.
     */
    public readonly noCache!: pulumi.Output<boolean | undefined>;
    /**
     * Set target platform(s) for the build. Defaults to the host's platform.
     *
     * Equivalent to Docker's `--platform` flag.
     */
    public readonly platforms!: pulumi.Output<enums.Platform[] | undefined>;
    /**
     * Always pull referenced images.
     *
     * Equivalent to Docker's `--pull` flag.
     */
    public readonly pull!: pulumi.Output<boolean | undefined>;
    /**
     * When `true` the build will automatically include a `registry` export.
     *
     * Defaults to `false`.
     *
     * Equivalent to Docker's `--push` flag.
     */
    public readonly push!: pulumi.Output<boolean | undefined>;
    /**
     * If the image was pushed to any registries then this will contain a
     * single fully-qualified tag including the build's digest.
     *
     * If the image had tags but was not exported, this will take on a value
     * of one of those tags.
     *
     * This will be empty if the image had no exports and no tags.
     *
     * This is only for convenience and may not be appropriate for situations
     * where multiple tags or registries are involved. In those cases this
     * output is not guaranteed to be stable.
     *
     * For more control over tags consumed by downstream resources you should
     * use the `digest` output.
     */
    public /*out*/ readonly ref!: pulumi.Output<string>;
    /**
     * Registry credentials. Required if reading or exporting to private
     * repositories.
     *
     * Credentials are kept in-memory and do not pollute pre-existing
     * credentials on the host.
     *
     * Similar to `docker login`.
     */
    public readonly registries!: pulumi.Output<outputs.Registry[] | undefined>;
    /**
     * A mapping of secret names to their corresponding values.
     *
     * Unlike the Docker CLI, these can be passed by value and do not need to
     * exist on-disk or in environment variables.
     *
     * Build arguments and environment variables are persistent in the final
     * image, so you should use this for sensitive values.
     *
     * Similar to Docker's `--secret` flag.
     */
    public readonly secrets!: pulumi.Output<{[key: string]: string} | undefined>;
    /**
     * SSH agent socket or keys to expose to the build.
     *
     * Equivalent to Docker's `--ssh` flag.
     */
    public readonly ssh!: pulumi.Output<outputs.SSH[] | undefined>;
    /**
     * Name and optionally a tag (format: `name:tag`).
     *
     * If exporting to a registry, the name should include the fully qualified
     * registry address (e.g. `docker.io/pulumi/pulumi:latest`).
     *
     * Equivalent to Docker's `--tag` flag.
     */
    public readonly tags!: pulumi.Output<string[] | undefined>;
    /**
     * Set the target build stage(s) to build.
     *
     * If not specified all targets will be built by default.
     *
     * Equivalent to Docker's `--target` flag.
     */
    public readonly target!: pulumi.Output<string | undefined>;

    /**
     * Create a Image resource with the given unique name, arguments, and options.
     *
     * @param name The _unique_ name of the resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param opts A bag of options that control this resource's behavior.
     */
    constructor(name: string, args?: ImageArgs, opts?: pulumi.CustomResourceOptions) {
        let resourceInputs: pulumi.Inputs = {};
        opts = opts || {};
        if (!opts.id) {
            resourceInputs["addHosts"] = args ? args.addHosts : undefined;
            resourceInputs["buildArgs"] = args ? args.buildArgs : undefined;
            resourceInputs["buildOnPreview"] = (args ? args.buildOnPreview : undefined) ?? true;
            resourceInputs["builder"] = args ? args.builder : undefined;
            resourceInputs["cacheFrom"] = args ? args.cacheFrom : undefined;
            resourceInputs["cacheTo"] = args ? args.cacheTo : undefined;
            resourceInputs["context"] = args ? args.context : undefined;
            resourceInputs["dockerfile"] = args ? args.dockerfile : undefined;
            resourceInputs["exec"] = args ? args.exec : undefined;
            resourceInputs["exports"] = args ? args.exports : undefined;
            resourceInputs["labels"] = args ? args.labels : undefined;
            resourceInputs["load"] = args ? args.load : undefined;
            resourceInputs["network"] = (args ? args.network : undefined) ?? "default";
            resourceInputs["noCache"] = args ? args.noCache : undefined;
            resourceInputs["platforms"] = args ? args.platforms : undefined;
            resourceInputs["pull"] = args ? args.pull : undefined;
            resourceInputs["push"] = args ? args.push : undefined;
            resourceInputs["registries"] = args ? args.registries : undefined;
            resourceInputs["secrets"] = args ? args.secrets : undefined;
            resourceInputs["ssh"] = args ? args.ssh : undefined;
            resourceInputs["tags"] = args ? args.tags : undefined;
            resourceInputs["target"] = args ? args.target : undefined;
            resourceInputs["contextHash"] = undefined /*out*/;
            resourceInputs["digest"] = undefined /*out*/;
            resourceInputs["ref"] = undefined /*out*/;
        } else {
            resourceInputs["addHosts"] = undefined /*out*/;
            resourceInputs["buildArgs"] = undefined /*out*/;
            resourceInputs["buildOnPreview"] = undefined /*out*/;
            resourceInputs["builder"] = undefined /*out*/;
            resourceInputs["cacheFrom"] = undefined /*out*/;
            resourceInputs["cacheTo"] = undefined /*out*/;
            resourceInputs["context"] = undefined /*out*/;
            resourceInputs["contextHash"] = undefined /*out*/;
            resourceInputs["digest"] = undefined /*out*/;
            resourceInputs["dockerfile"] = undefined /*out*/;
            resourceInputs["exec"] = undefined /*out*/;
            resourceInputs["exports"] = undefined /*out*/;
            resourceInputs["labels"] = undefined /*out*/;
            resourceInputs["load"] = undefined /*out*/;
            resourceInputs["network"] = undefined /*out*/;
            resourceInputs["noCache"] = undefined /*out*/;
            resourceInputs["platforms"] = undefined /*out*/;
            resourceInputs["pull"] = undefined /*out*/;
            resourceInputs["push"] = undefined /*out*/;
            resourceInputs["ref"] = undefined /*out*/;
            resourceInputs["registries"] = undefined /*out*/;
            resourceInputs["secrets"] = undefined /*out*/;
            resourceInputs["ssh"] = undefined /*out*/;
            resourceInputs["tags"] = undefined /*out*/;
            resourceInputs["target"] = undefined /*out*/;
        }
        opts = pulumi.mergeOptions(utilities.resourceOptsDefaults(), opts);
        super(Image.__pulumiType, name, resourceInputs, opts);
    }
}

/**
 * The set of arguments for constructing a Image resource.
 */
export interface ImageArgs {
    /**
     * Custom `host:ip` mappings to use during the build.
     *
     * Equivalent to Docker's `--add-host` flag.
     */
    addHosts?: pulumi.Input<pulumi.Input<string>[]>;
    /**
     * `ARG` names and values to set during the build.
     *
     * These variables are accessed like environment variables inside `RUN`
     * instructions.
     *
     * Build arguments are persisted in the image, so you should use `secrets`
     * if these arguments are sensitive.
     *
     * Equivalent to Docker's `--build-arg` flag.
     */
    buildArgs?: pulumi.Input<{[key: string]: pulumi.Input<string>}>;
    /**
     * Setting this to `false` will always skip image builds during previews,
     * and setting it to `true` will always build images during previews.
     *
     * Images built during previews are never exported to registries, however
     * cache manifests are still exported.
     *
     * On-disk Dockerfiles are always validated for syntactic correctness
     * regardless of this setting.
     *
     * Defaults to `true` as a safeguard against broken images merging as part
     * of CI pipelines.
     */
    buildOnPreview?: pulumi.Input<boolean>;
    /**
     * Builder configuration.
     */
    builder?: pulumi.Input<inputs.BuilderConfigArgs>;
    /**
     * Cache export configuration.
     *
     * Equivalent to Docker's `--cache-from` flag.
     */
    cacheFrom?: pulumi.Input<pulumi.Input<inputs.CacheFromArgs>[]>;
    /**
     * Cache import configuration.
     *
     * Equivalent to Docker's `--cache-to` flag.
     */
    cacheTo?: pulumi.Input<pulumi.Input<inputs.CacheToArgs>[]>;
    /**
     * Build context settings.
     *
     * Equivalent to Docker's `PATH | URL | -` positional argument.
     */
    context?: pulumi.Input<inputs.BuildContextArgs>;
    /**
     * Dockerfile settings.
     *
     * Equivalent to Docker's `--file` flag.
     */
    dockerfile?: pulumi.Input<inputs.DockerfileArgs>;
    /**
     * Use `exec` mode to build this image.
     *
     * By default the provider embeds a v25 Docker client with v0.12 buildx
     * support. This helps ensure consistent behavior across environments and
     * is compatible with alternative build backends (e.g. `buildkitd`), but
     * it may not be desirable if you require a specific version of buildx.
     * For example you may want to run a custom `docker-buildx` binary with
     * support for [Docker Build
     * Cloud](https://docs.docker.com/build/cloud/setup/) (DBC).
     *
     * When this is set to `true` the provider will instead execute the
     * `docker-buildx` binary directly to perform its operations. The user is
     * responsible for ensuring this binary exists, with correct permissions
     * and pre-configured builders, at a path Docker expects (e.g.
     * `~/.docker/cli-plugins`).
     *
     * Debugging `exec` mode may be more difficult as Pulumi will not be able
     * to surface fine-grained errors and warnings. Additionally credentials
     * are temporarily written to disk in order to provide them to the
     * `docker-buildx` binary.
     */
    exec?: pulumi.Input<boolean>;
    /**
     * Controls where images are persisted after building.
     *
     * Images are only stored in the local cache unless `exports` are
     * explicitly configured.
     *
     * Exporting to multiple destinations requires a daemon running BuildKit
     * 0.13 or later.
     *
     * Equivalent to Docker's `--output` flag.
     */
    exports?: pulumi.Input<pulumi.Input<inputs.ExportArgs>[]>;
    /**
     * Attach arbitrary key/value metadata to the image.
     *
     * Equivalent to Docker's `--label` flag.
     */
    labels?: pulumi.Input<{[key: string]: pulumi.Input<string>}>;
    /**
     * When `true` the build will automatically include a `docker` export.
     *
     * Defaults to `false`.
     *
     * Equivalent to Docker's `--load` flag.
     */
    load?: pulumi.Input<boolean>;
    /**
     * Set the network mode for `RUN` instructions. Defaults to `default`.
     *
     * For custom networks, configure your builder with `--driver-opt network=...`.
     *
     * Equivalent to Docker's `--network` flag.
     */
    network?: pulumi.Input<enums.NetworkMode>;
    /**
     * Do not import cache manifests when building the image.
     *
     * Equivalent to Docker's `--no-cache` flag.
     */
    noCache?: pulumi.Input<boolean>;
    /**
     * Set target platform(s) for the build. Defaults to the host's platform.
     *
     * Equivalent to Docker's `--platform` flag.
     */
    platforms?: pulumi.Input<pulumi.Input<enums.Platform>[]>;
    /**
     * Always pull referenced images.
     *
     * Equivalent to Docker's `--pull` flag.
     */
    pull?: pulumi.Input<boolean>;
    /**
     * When `true` the build will automatically include a `registry` export.
     *
     * Defaults to `false`.
     *
     * Equivalent to Docker's `--push` flag.
     */
    push?: pulumi.Input<boolean>;
    /**
     * Registry credentials. Required if reading or exporting to private
     * repositories.
     *
     * Credentials are kept in-memory and do not pollute pre-existing
     * credentials on the host.
     *
     * Similar to `docker login`.
     */
    registries?: pulumi.Input<pulumi.Input<inputs.RegistryArgs>[]>;
    /**
     * A mapping of secret names to their corresponding values.
     *
     * Unlike the Docker CLI, these can be passed by value and do not need to
     * exist on-disk or in environment variables.
     *
     * Build arguments and environment variables are persistent in the final
     * image, so you should use this for sensitive values.
     *
     * Similar to Docker's `--secret` flag.
     */
    secrets?: pulumi.Input<{[key: string]: pulumi.Input<string>}>;
    /**
     * SSH agent socket or keys to expose to the build.
     *
     * Equivalent to Docker's `--ssh` flag.
     */
    ssh?: pulumi.Input<pulumi.Input<inputs.SSHArgs>[]>;
    /**
     * Name and optionally a tag (format: `name:tag`).
     *
     * If exporting to a registry, the name should include the fully qualified
     * registry address (e.g. `docker.io/pulumi/pulumi:latest`).
     *
     * Equivalent to Docker's `--tag` flag.
     */
    tags?: pulumi.Input<pulumi.Input<string>[]>;
    /**
     * Set the target build stage(s) to build.
     *
     * If not specified all targets will be built by default.
     *
     * Equivalent to Docker's `--target` flag.
     */
    target?: pulumi.Input<string>;
}
