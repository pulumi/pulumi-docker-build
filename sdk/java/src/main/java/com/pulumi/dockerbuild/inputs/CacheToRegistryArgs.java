// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.dockerbuild.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.core.internal.Codegen;
import com.pulumi.dockerbuild.enums.CacheMode;
import com.pulumi.dockerbuild.enums.CompressionType;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.Boolean;
import java.lang.Integer;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


public final class CacheToRegistryArgs extends com.pulumi.resources.ResourceArgs {

    public static final CacheToRegistryArgs Empty = new CacheToRegistryArgs();

    /**
     * The compression type to use.
     * 
     */
    @Import(name="compression")
    private @Nullable Output<CompressionType> compression;

    /**
     * @return The compression type to use.
     * 
     */
    public Optional<Output<CompressionType>> compression() {
        return Optional.ofNullable(this.compression);
    }

    /**
     * Compression level from 0 to 22.
     * 
     */
    @Import(name="compressionLevel")
    private @Nullable Output<Integer> compressionLevel;

    /**
     * @return Compression level from 0 to 22.
     * 
     */
    public Optional<Output<Integer>> compressionLevel() {
        return Optional.ofNullable(this.compressionLevel);
    }

    /**
     * Forcefully apply compression.
     * 
     */
    @Import(name="forceCompression")
    private @Nullable Output<Boolean> forceCompression;

    /**
     * @return Forcefully apply compression.
     * 
     */
    public Optional<Output<Boolean>> forceCompression() {
        return Optional.ofNullable(this.forceCompression);
    }

    /**
     * Ignore errors caused by failed cache exports.
     * 
     */
    @Import(name="ignoreError")
    private @Nullable Output<Boolean> ignoreError;

    /**
     * @return Ignore errors caused by failed cache exports.
     * 
     */
    public Optional<Output<Boolean>> ignoreError() {
        return Optional.ofNullable(this.ignoreError);
    }

    /**
     * Export cache manifest as an OCI-compatible image manifest instead of a
     * manifest list. Requires `ociMediaTypes` to also be `true`.
     * 
     * Some registries like AWS ECR will not work with caching if this is
     * `false`.
     * 
     * Defaults to `false` to match Docker&#39;s default behavior.
     * 
     */
    @Import(name="imageManifest")
    private @Nullable Output<Boolean> imageManifest;

    /**
     * @return Export cache manifest as an OCI-compatible image manifest instead of a
     * manifest list. Requires `ociMediaTypes` to also be `true`.
     * 
     * Some registries like AWS ECR will not work with caching if this is
     * `false`.
     * 
     * Defaults to `false` to match Docker&#39;s default behavior.
     * 
     */
    public Optional<Output<Boolean>> imageManifest() {
        return Optional.ofNullable(this.imageManifest);
    }

    /**
     * The cache mode to use. Defaults to `min`.
     * 
     */
    @Import(name="mode")
    private @Nullable Output<CacheMode> mode;

    /**
     * @return The cache mode to use. Defaults to `min`.
     * 
     */
    public Optional<Output<CacheMode>> mode() {
        return Optional.ofNullable(this.mode);
    }

    /**
     * Whether to use OCI media types in exported manifests. Defaults to
     * `true`.
     * 
     */
    @Import(name="ociMediaTypes")
    private @Nullable Output<Boolean> ociMediaTypes;

    /**
     * @return Whether to use OCI media types in exported manifests. Defaults to
     * `true`.
     * 
     */
    public Optional<Output<Boolean>> ociMediaTypes() {
        return Optional.ofNullable(this.ociMediaTypes);
    }

    /**
     * Fully qualified name of the cache image to import.
     * 
     */
    @Import(name="ref", required=true)
    private Output<String> ref;

    /**
     * @return Fully qualified name of the cache image to import.
     * 
     */
    public Output<String> ref() {
        return this.ref;
    }

    private CacheToRegistryArgs() {}

    private CacheToRegistryArgs(CacheToRegistryArgs $) {
        this.compression = $.compression;
        this.compressionLevel = $.compressionLevel;
        this.forceCompression = $.forceCompression;
        this.ignoreError = $.ignoreError;
        this.imageManifest = $.imageManifest;
        this.mode = $.mode;
        this.ociMediaTypes = $.ociMediaTypes;
        this.ref = $.ref;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(CacheToRegistryArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private CacheToRegistryArgs $;

        public Builder() {
            $ = new CacheToRegistryArgs();
        }

        public Builder(CacheToRegistryArgs defaults) {
            $ = new CacheToRegistryArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param compression The compression type to use.
         * 
         * @return builder
         * 
         */
        public Builder compression(@Nullable Output<CompressionType> compression) {
            $.compression = compression;
            return this;
        }

        /**
         * @param compression The compression type to use.
         * 
         * @return builder
         * 
         */
        public Builder compression(CompressionType compression) {
            return compression(Output.of(compression));
        }

        /**
         * @param compressionLevel Compression level from 0 to 22.
         * 
         * @return builder
         * 
         */
        public Builder compressionLevel(@Nullable Output<Integer> compressionLevel) {
            $.compressionLevel = compressionLevel;
            return this;
        }

        /**
         * @param compressionLevel Compression level from 0 to 22.
         * 
         * @return builder
         * 
         */
        public Builder compressionLevel(Integer compressionLevel) {
            return compressionLevel(Output.of(compressionLevel));
        }

        /**
         * @param forceCompression Forcefully apply compression.
         * 
         * @return builder
         * 
         */
        public Builder forceCompression(@Nullable Output<Boolean> forceCompression) {
            $.forceCompression = forceCompression;
            return this;
        }

        /**
         * @param forceCompression Forcefully apply compression.
         * 
         * @return builder
         * 
         */
        public Builder forceCompression(Boolean forceCompression) {
            return forceCompression(Output.of(forceCompression));
        }

        /**
         * @param ignoreError Ignore errors caused by failed cache exports.
         * 
         * @return builder
         * 
         */
        public Builder ignoreError(@Nullable Output<Boolean> ignoreError) {
            $.ignoreError = ignoreError;
            return this;
        }

        /**
         * @param ignoreError Ignore errors caused by failed cache exports.
         * 
         * @return builder
         * 
         */
        public Builder ignoreError(Boolean ignoreError) {
            return ignoreError(Output.of(ignoreError));
        }

        /**
         * @param imageManifest Export cache manifest as an OCI-compatible image manifest instead of a
         * manifest list. Requires `ociMediaTypes` to also be `true`.
         * 
         * Some registries like AWS ECR will not work with caching if this is
         * `false`.
         * 
         * Defaults to `false` to match Docker&#39;s default behavior.
         * 
         * @return builder
         * 
         */
        public Builder imageManifest(@Nullable Output<Boolean> imageManifest) {
            $.imageManifest = imageManifest;
            return this;
        }

        /**
         * @param imageManifest Export cache manifest as an OCI-compatible image manifest instead of a
         * manifest list. Requires `ociMediaTypes` to also be `true`.
         * 
         * Some registries like AWS ECR will not work with caching if this is
         * `false`.
         * 
         * Defaults to `false` to match Docker&#39;s default behavior.
         * 
         * @return builder
         * 
         */
        public Builder imageManifest(Boolean imageManifest) {
            return imageManifest(Output.of(imageManifest));
        }

        /**
         * @param mode The cache mode to use. Defaults to `min`.
         * 
         * @return builder
         * 
         */
        public Builder mode(@Nullable Output<CacheMode> mode) {
            $.mode = mode;
            return this;
        }

        /**
         * @param mode The cache mode to use. Defaults to `min`.
         * 
         * @return builder
         * 
         */
        public Builder mode(CacheMode mode) {
            return mode(Output.of(mode));
        }

        /**
         * @param ociMediaTypes Whether to use OCI media types in exported manifests. Defaults to
         * `true`.
         * 
         * @return builder
         * 
         */
        public Builder ociMediaTypes(@Nullable Output<Boolean> ociMediaTypes) {
            $.ociMediaTypes = ociMediaTypes;
            return this;
        }

        /**
         * @param ociMediaTypes Whether to use OCI media types in exported manifests. Defaults to
         * `true`.
         * 
         * @return builder
         * 
         */
        public Builder ociMediaTypes(Boolean ociMediaTypes) {
            return ociMediaTypes(Output.of(ociMediaTypes));
        }

        /**
         * @param ref Fully qualified name of the cache image to import.
         * 
         * @return builder
         * 
         */
        public Builder ref(Output<String> ref) {
            $.ref = ref;
            return this;
        }

        /**
         * @param ref Fully qualified name of the cache image to import.
         * 
         * @return builder
         * 
         */
        public Builder ref(String ref) {
            return ref(Output.of(ref));
        }

        public CacheToRegistryArgs build() {
            $.compression = Codegen.objectProp("compression", CompressionType.class).output().arg($.compression).def(CompressionType.Gzip).getNullable();
            $.compressionLevel = Codegen.integerProp("compressionLevel").output().arg($.compressionLevel).def(0).getNullable();
            $.forceCompression = Codegen.booleanProp("forceCompression").output().arg($.forceCompression).def(false).getNullable();
            $.ignoreError = Codegen.booleanProp("ignoreError").output().arg($.ignoreError).def(false).getNullable();
            $.imageManifest = Codegen.booleanProp("imageManifest").output().arg($.imageManifest).def(false).getNullable();
            $.mode = Codegen.objectProp("mode", CacheMode.class).output().arg($.mode).def(CacheMode.Min).getNullable();
            $.ociMediaTypes = Codegen.booleanProp("ociMediaTypes").output().arg($.ociMediaTypes).def(true).getNullable();
            if ($.ref == null) {
                throw new MissingRequiredPropertyException("CacheToRegistryArgs", "ref");
            }
            return $;
        }
    }

}
