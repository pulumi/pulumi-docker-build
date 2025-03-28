// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.dockerbuild;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.core.internal.Codegen;
import com.pulumi.dockerbuild.inputs.RegistryArgs;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


public final class ProviderArgs extends com.pulumi.resources.ResourceArgs {

    public static final ProviderArgs Empty = new ProviderArgs();

    /**
     * The build daemon&#39;s address.
     * 
     */
    @Import(name="host")
    private @Nullable Output<String> host;

    /**
     * @return The build daemon&#39;s address.
     * 
     */
    public Optional<Output<String>> host() {
        return Optional.ofNullable(this.host);
    }

    @Import(name="registries", json=true)
    private @Nullable Output<List<RegistryArgs>> registries;

    public Optional<Output<List<RegistryArgs>>> registries() {
        return Optional.ofNullable(this.registries);
    }

    private ProviderArgs() {}

    private ProviderArgs(ProviderArgs $) {
        this.host = $.host;
        this.registries = $.registries;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(ProviderArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private ProviderArgs $;

        public Builder() {
            $ = new ProviderArgs();
        }

        public Builder(ProviderArgs defaults) {
            $ = new ProviderArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param host The build daemon&#39;s address.
         * 
         * @return builder
         * 
         */
        public Builder host(@Nullable Output<String> host) {
            $.host = host;
            return this;
        }

        /**
         * @param host The build daemon&#39;s address.
         * 
         * @return builder
         * 
         */
        public Builder host(String host) {
            return host(Output.of(host));
        }

        public Builder registries(@Nullable Output<List<RegistryArgs>> registries) {
            $.registries = registries;
            return this;
        }

        public Builder registries(List<RegistryArgs> registries) {
            return registries(Output.of(registries));
        }

        public Builder registries(RegistryArgs... registries) {
            return registries(List.of(registries));
        }

        public ProviderArgs build() {
            $.host = Codegen.stringProp("host").output().arg($.host).env("DOCKER_HOST").def("").getNullable();
            return $;
        }
    }

}
