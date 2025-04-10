// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.dockerbuild.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.String;
import java.util.Objects;

@CustomType
public final class CacheFromRegistry {
    /**
     * @return Fully qualified name of the cache image to import.
     * 
     */
    private String ref;

    private CacheFromRegistry() {}
    /**
     * @return Fully qualified name of the cache image to import.
     * 
     */
    public String ref() {
        return this.ref;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(CacheFromRegistry defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private String ref;
        public Builder() {}
        public Builder(CacheFromRegistry defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.ref = defaults.ref;
        }

        @CustomType.Setter
        public Builder ref(String ref) {
            if (ref == null) {
              throw new MissingRequiredPropertyException("CacheFromRegistry", "ref");
            }
            this.ref = ref;
            return this;
        }
        public CacheFromRegistry build() {
            final var _resultValue = new CacheFromRegistry();
            _resultValue.ref = ref;
            return _resultValue;
        }
    }
}
