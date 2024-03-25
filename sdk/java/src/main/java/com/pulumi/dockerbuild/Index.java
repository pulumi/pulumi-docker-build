// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

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
 * An index (or manifest list) referencing one or more existing images.
 * 
 * Useful for crafting a multi-platform image from several
 * platform-specific images.
 * 
 * This creates an OCI image index or a Docker manifest list depending on
 * the media types of the source images.
 * 
 * ## Example Usage
 * ### Multi-platform registry caching
 * ```java
 * package generated_program;
 * 
 * import com.pulumi.Context;
 * import com.pulumi.Pulumi;
 * import com.pulumi.core.Output;
 * import com.pulumi.docker.buildx.Image;
 * import com.pulumi.docker.buildx.ImageArgs;
 * import com.pulumi.docker.buildx.inputs.CacheFromArgs;
 * import com.pulumi.docker.buildx.inputs.CacheFromRegistryArgs;
 * import com.pulumi.docker.buildx.inputs.CacheToArgs;
 * import com.pulumi.docker.buildx.inputs.CacheToRegistryArgs;
 * import com.pulumi.docker.buildx.inputs.BuildContextArgs;
 * import com.pulumi.docker.buildx.Index;
 * import com.pulumi.docker.buildx.IndexArgs;
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
 *         var amd64 = new Image(&#34;amd64&#34;, ImageArgs.builder()        
 *             .cacheFrom(CacheFromArgs.builder()
 *                 .registry(CacheFromRegistryArgs.builder()
 *                     .ref(&#34;docker.io/pulumi/pulumi:cache-amd64&#34;)
 *                     .build())
 *                 .build())
 *             .cacheTo(CacheToArgs.builder()
 *                 .registry(CacheToRegistryArgs.builder()
 *                     .mode(&#34;max&#34;)
 *                     .ref(&#34;docker.io/pulumi/pulumi:cache-amd64&#34;)
 *                     .build())
 *                 .build())
 *             .context(BuildContextArgs.builder()
 *                 .location(&#34;app&#34;)
 *                 .build())
 *             .platforms(&#34;linux/amd64&#34;)
 *             .tags(&#34;docker.io/pulumi/pulumi:3.107.0-amd64&#34;)
 *             .build());
 * 
 *         var arm64 = new Image(&#34;arm64&#34;, ImageArgs.builder()        
 *             .cacheFrom(CacheFromArgs.builder()
 *                 .registry(CacheFromRegistryArgs.builder()
 *                     .ref(&#34;docker.io/pulumi/pulumi:cache-arm64&#34;)
 *                     .build())
 *                 .build())
 *             .cacheTo(CacheToArgs.builder()
 *                 .registry(CacheToRegistryArgs.builder()
 *                     .mode(&#34;max&#34;)
 *                     .ref(&#34;docker.io/pulumi/pulumi:cache-arm64&#34;)
 *                     .build())
 *                 .build())
 *             .context(BuildContextArgs.builder()
 *                 .location(&#34;app&#34;)
 *                 .build())
 *             .platforms(&#34;linux/arm64&#34;)
 *             .tags(&#34;docker.io/pulumi/pulumi:3.107.0-arm64&#34;)
 *             .build());
 * 
 *         var index = new Index(&#34;index&#34;, IndexArgs.builder()        
 *             .sources(            
 *                 amd64.ref(),
 *                 arm64.ref())
 *             .tag(&#34;docker.io/pulumi/pulumi:3.107.0&#34;)
 *             .build());
 * 
 *         ctx.export(&#34;ref&#34;, index.ref());
 *     }
 * }
 * ```
 * 
 */
@ResourceType(type="dockerbuild:index:Index")
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
    public Index(String name) {
        this(name, IndexArgs.Empty);
    }
    /**
     *
     * @param name The _unique_ name of the resulting resource.
     * @param args The arguments to use to populate this resource's properties.
     */
    public Index(String name, IndexArgs args) {
        this(name, args, null);
    }
    /**
     *
     * @param name The _unique_ name of the resulting resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param options A bag of options that control this resource's behavior.
     */
    public Index(String name, IndexArgs args, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        super("dockerbuild:index:Index", name, args == null ? IndexArgs.Empty : args, makeResourceOptions(options, Codegen.empty()));
    }

    private Index(String name, Output<String> id, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        super("dockerbuild:index:Index", name, null, makeResourceOptions(options, id));
    }

    private static com.pulumi.resources.CustomResourceOptions makeResourceOptions(@Nullable com.pulumi.resources.CustomResourceOptions options, @Nullable Output<String> id) {
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
    public static Index get(String name, Output<String> id, @Nullable com.pulumi.resources.CustomResourceOptions options) {
        return new Index(name, id, options);
    }
}
