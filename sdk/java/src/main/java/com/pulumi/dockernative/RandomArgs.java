// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.dockernative;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.Integer;
import java.util.Objects;


public final class RandomArgs extends com.pulumi.resources.ResourceArgs {

    public static final RandomArgs Empty = new RandomArgs();

    @Import(name="length", required=true)
    private Output<Integer> length;

    public Output<Integer> length() {
        return this.length;
    }

    private RandomArgs() {}

    private RandomArgs(RandomArgs $) {
        this.length = $.length;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(RandomArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private RandomArgs $;

        public Builder() {
            $ = new RandomArgs();
        }

        public Builder(RandomArgs defaults) {
            $ = new RandomArgs(Objects.requireNonNull(defaults));
        }

        public Builder length(Output<Integer> length) {
            $.length = length;
            return this;
        }

        public Builder length(Integer length) {
            return length(Output.of(length));
        }

        public RandomArgs build() {
            if ($.length == null) {
                throw new MissingRequiredPropertyException("RandomArgs", "length");
            }
            return $;
        }
    }

}
