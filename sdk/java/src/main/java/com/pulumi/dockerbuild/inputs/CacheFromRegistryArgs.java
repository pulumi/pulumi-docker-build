// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.dockerbuild.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.String;
import java.util.Objects;


public final class CacheFromRegistryArgs extends com.pulumi.resources.ResourceArgs {

    public static final CacheFromRegistryArgs Empty = new CacheFromRegistryArgs();

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

    private CacheFromRegistryArgs() {}

    private CacheFromRegistryArgs(CacheFromRegistryArgs $) {
        this.ref = $.ref;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(CacheFromRegistryArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private CacheFromRegistryArgs $;

        public Builder() {
            $ = new CacheFromRegistryArgs();
        }

        public Builder(CacheFromRegistryArgs defaults) {
            $ = new CacheFromRegistryArgs(Objects.requireNonNull(defaults));
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

        public CacheFromRegistryArgs build() {
            if ($.ref == null) {
                throw new MissingRequiredPropertyException("CacheFromRegistryArgs", "ref");
            }
            return $;
        }
    }

}