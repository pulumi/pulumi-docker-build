// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.dockerbuild.enums;

import com.pulumi.core.annotations.EnumType;
import java.lang.String;
import java.util.Objects;
import java.util.StringJoiner;

    @EnumType
    public enum CacheMode {
        /**
         * Only layers that are exported into the resulting image are cached.
         * 
         */
        Min("min"),
        /**
         * All layers are cached, even those of intermediate steps.
         * 
         */
        Max("max");

        private final String value;

        CacheMode(String value) {
            this.value = Objects.requireNonNull(value);
        }

        @EnumType.Converter
        public String getValue() {
            return this.value;
        }

        @Override
        public java.lang.String toString() {
            return new StringJoiner(", ", "CacheMode[", "]")
                .add("value='" + this.value + "'")
                .toString();
        }
    }
