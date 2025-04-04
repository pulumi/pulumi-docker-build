// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.dockerbuild.enums;

import com.pulumi.core.annotations.EnumType;
import java.lang.String;
import java.util.Objects;
import java.util.StringJoiner;

    @EnumType
    public enum CompressionType {
        /**
         * Use `gzip` for compression.
         * 
         */
        Gzip("gzip"),
        /**
         * Use `estargz` for compression.
         * 
         */
        Estargz("estargz"),
        /**
         * Use `zstd` for compression.
         * 
         */
        Zstd("zstd");

        private final String value;

        CompressionType(String value) {
            this.value = Objects.requireNonNull(value);
        }

        @EnumType.Converter
        public String getValue() {
            return this.value;
        }

        @Override
        public java.lang.String toString() {
            return new StringJoiner(", ", "CompressionType[", "]")
                .add("value='" + this.value + "'")
                .toString();
        }
    }
