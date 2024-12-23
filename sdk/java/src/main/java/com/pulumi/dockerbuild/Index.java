package com.pulumi.dockerbuild;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Export;
import com.pulumi.core.annotations.ResourceType;
import com.pulumi.core.internal.Codegen;
import com.pulumi.dockerbuild.IndexArgs;
import com.pulumi.dockerbuild.Utilities;
import com.pulumi.dockerbuild.outputs.Registry;
import java.lang.Boolean;
import java.lang.String;
import java.util.List;
import java.util.Optional;
import javax.annotation.Nullable;

/**
 * A wrapper around `docker buildx imagetools create` to create an index
 * (or manifest list) referencing one or more existing images.
 * 
 * In most cases you do not need an `Index` to build a multi-platform
 * image -- specifying multiple platforms on the `Image` will handle this
 * for you automatically.
 * 
 * However, as of April 2024, building multi-platform images _with
 * caching_ will only export a cache for one platform at a time (see [this
 * discussion](https://github.com/docker/buildx/discussions/1382) for more
 * details).
 * 
 * Therefore this resource can be helpful if you are building
 * multi-platform images with caching: each platform can be built and
 * cached separately, and an `Index` can join them all together. An
 * example of this is shown below.
 * 
 * This resource creates an OCI image index or a Docker manifest list
 * depending on the media types of the source images.
 * 
 * ## Example Usage
 * ### Multi-platform registry caching
 * <pre>
 * {@code
 * package generated_program;
 * 
 * import com.pulumi.Context;
 * import com.pulumi.Pulumi;
 * import com.pulumi.core.Output;
 * import com.pulumi.dockerbuild.Image;
 * import com.pulumi.dockerbuild.ImageArgs;
 * import com.pulumi.dockerbuild.inputs.CacheFromArgs;
 * import com.pulumi.dockerbuild.inputs.CacheFromRegistryArgs;
 * import com.pulumi.dockerbuild.inputs.CacheToArgs;
 * import com.pulumi.dockerbuild.inputs.CacheToRegistryArgs;
 * import com.pulumi.dockerbuild.inputs.BuildContextArgs;
 * import com.pulumi.dockerbuild.Index;
 * import com.pulumi.dockerbuild.IndexArgs;
 * import java.util.List;
 * import java.util.ArrayList;
 * import java.util.Map;
 * import java.io.File;
 * import java.nio.file.Files;
 * import java.nio.file.Paths;
 * 
 * public class App {
 *     public static void main(String[] args) {
 *         Pulumi.run(App::stack);
 *     }
 * 
 *     public static void stack(Context ctx) {
 *         var amd64 = new Image("amd64", ImageArgs.builder()
 *             .cacheFrom(CacheFromArgs.builder()
 *                 .registry(CacheFromRegistryArgs.builder()
 *                     .ref("docker.io/pulumi/pulumi:cache-amd64")
 *                     .build())
 *                 .build())
 *             .cacheTo(CacheToArgs.builder()
 *                 .registry(CacheToRegistryArgs.builder()
 *                     .mode("max")
 *                     .ref("docker.io/pulumi/pulumi:cache-amd64")
 *                     .build())
 *                 .build())
 *             .context(BuildContextArgs.builder()
 *                 .location("app")
 *                 .build())
 *             .platforms("linux/amd64")
 *             .tags("docker.io/pulumi/pulumi:3.107.0-amd64")
 *             .build());
 * 
 *         var arm64 = new Image("arm64", ImageArgs.builder()
 *             .cacheFrom(CacheFromArgs.builder()
 *                 .registry(CacheFromRegistryArgs.builder()
 *                     .ref("docker.io/pulumi/pulumi:cache-arm64")
 *                     .build())
 *                 .build())
 *             .cacheTo(CacheToArgs.builder()
 *                 .registry(CacheToRegistryArgs.builder()
 *                     .mode("max")
 *                     .ref("docker.io/pulumi/pulumi:cache-arm64")
 *                     .build())
 *                 .build())
 *             .context(BuildContextArgs.builder()
 *                 .location("app")
 *                 .build())
 *             .platforms("linux/arm64")
 *             .tags("docker.io/pulumi/pulumi:3.107.0-arm64")
 *             .build());
 * 
 *         var index = new Index("index", IndexArgs.builder()
 *             .sources(            
 *                 amd64.ref(),
 *                 arm64.ref())
 *             .tag("docker.io/pulumi/pulumi:3.107.0")
 *             .build());
 * 
 *         ctx.export("ref", index.ref());
 *     }
 * }
 * }
 * </pre>
 * 
 */
@ResourceType(type="docker-build:index:Index")
public class Index extends com.pulumi.resources.CustomResource {
    /**
     * If true, push the index to the target registry.
     * 
     * Defaults to `true`.
     * 
     */
    @Export(name="push", refs={Boolean.class}, tree="[0]")
    private Output</* @Nullable */ Boolean> push;

    /**
     * @return If true, push the index to the target registry.
     * 
     * Defaults to `true`.
     * 
     */
    public Output<Optional<Boolean>> push() {
        return Codegen.optional(this.push);
    }
    /**
     * The pushed tag with digest.
     * 
     * Identical to the tag if the index was not pushed.
     * 
     */
    @Export(name="ref", refs={String.class}, tree="[0]")
    private Output<String> ref;

    /**
     * @return The pushed tag with digest.
     * 
     * Identical to the tag if the index was not pushed.
     * 
     */
    public Output<String> ref() {
        return this.ref;
    }
    /**
     * Authentication for the registry where the tagged index will be pushed.
     * 
     * Credentials can also be included with the provider&#39;s configuration.
     * 
     */
    @Export(name="registry", refs={Registry.class}, tree="[0]")
    private Output</* @Nullable */ Registry> registry;

    /**
     * @return Authentication for the registry where the tagged index will be pushed.
     * 
     * Credentials can also be included with the provider&#39;s configuration.
     * 
     */
    public Output<Optional<Registry>> registry() {
        return Codegen.optional(this.registry);
    }
    /**
     * Existing images to include in the index.
     * 
     */
    @Export(name="sources", refs={List.class,String.class}, tree="[0,1]")
    private Output<List<String>> sources;

    /**
     * @return Existing images to include in the index.
     * 
     */
    public Output<List<String>> sources() {
        return this.sources;
    }
    /**
     * The tag to apply to the index.
     * 
     */
    @Export(name="tag", refs={String.class}, tree="[0]")
    private Output<String> tag;

    /**
     * @return The tag to apply to the index.
     * 
     */
    public Output<String> tag() {
        return this.tag;
    }

    /**
     *
     * @param name The _unique_ name of the resulting resource.
     */
    public Index(java.lang.String name) {
        this(name, IndexArgs.Empty);
    }
    /**
     *
     * @param name The _unique_ name of the resulting resource.
     * @param args The arguments to use to populate this resource's properties.
     */
    public Index(java.lang.String name, IndexArgs args) {
        this(name, args, null);
    }
    /**
     *
     * @param name The _unique_ name of the resulting resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param options A bag of options that control this resource's behavior.
     */
    public Index(java.lang.String name, IndexArgs args, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        super("docker-build:index:Index", name, makeArgs(args, options), makeResourceOptions(options, Codegen.empty()), false);
    }

    private Index(java.lang.String name, Output<java.lang.String> id, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        super("docker-build:index:Index", name, null, makeResourceOptions(options, id), false);
    }

    private static IndexArgs makeArgs(IndexArgs args, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        if (options != null && options.getUrn().isPresent()) {
            return null;
        }
        return args == null ? IndexArgs.Empty : args;
    }

    private static com.pulumi.resources.CustomResourceOptions makeResourceOptions(@Nullable com.pulumi.resources.CustomResourceOptions options, @Nullable Output<java.lang.String> id) {
        var defaultOptions = com.pulumi.resources.CustomResourceOptions.builder()
            .version(Utilities.getVersion())
            .build();
        return com.pulumi.resources.CustomResourceOptions.merge(defaultOptions, options, id);
    }

    /**
     * Get an existing Host resource's state with the given name, ID, and optional extra
     * properties used to qualify the lookup.
     *
     * @param name The _unique_ name of the resulting resource.
     * @param id The _unique_ provider ID of the resource to lookup.
     * @param options Optional settings to control the behavior of the CustomResource.
     */
    public static Index get(java.lang.String name, Output<java.lang.String> id, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        return new Index(name, id, options);
    }
}
