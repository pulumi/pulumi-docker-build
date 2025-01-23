// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.dockerbuild.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.String;
import java.util.Objects;

@CustomType
public final class ExportTar {
    /**
     * @return Output path.
     * 
     */
    private String dest;

    private ExportTar() {}
    /**
     * @return Output path.
     * 
     */
    public String dest() {
        return this.dest;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(ExportTar defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private String dest;
        public Builder() {}
        public Builder(ExportTar defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.dest = defaults.dest;
        }

        @CustomType.Setter
        public Builder dest(String dest) {
            if (dest == null) {
              throw new MissingRequiredPropertyException("ExportTar", "dest");
            }
            this.dest = dest;
            return this;
        }
        public ExportTar build() {
            final var _resultValue = new ExportTar();
            _resultValue.dest = dest;
            return _resultValue;
        }
    }
}
