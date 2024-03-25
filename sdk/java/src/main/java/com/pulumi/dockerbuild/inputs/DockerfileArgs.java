// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.dockerbuild.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


public final class DockerfileArgs extends com.pulumi.resources.ResourceArgs {

    public static final DockerfileArgs Empty = new DockerfileArgs();

    /**
     * Raw Dockerfile contents.
     * 
     * Conflicts with `location`.
     * 
     * Equivalent to invoking Docker with `-f -`.
     * 
     */
    @Import(name="inline")
    private @Nullable Output<String> inline;

    /**
     * @return Raw Dockerfile contents.
     * 
     * Conflicts with `location`.
     * 
     * Equivalent to invoking Docker with `-f -`.
     * 
     */
    public Optional<Output<String>> inline() {
        return Optional.ofNullable(this.inline);
    }

    /**
     * Location of the Dockerfile to use.
     * 
     * Can be a relative or absolute path to a local file, or a remote URL.
     * 
     * Defaults to `${context.location}/Dockerfile` if context is on-disk.
     * 
     * Conflicts with `inline`.
     * 
     */
    @Import(name="location")
    private @Nullable Output<String> location;

    /**
     * @return Location of the Dockerfile to use.
     * 
     * Can be a relative or absolute path to a local file, or a remote URL.
     * 
     * Defaults to `${context.location}/Dockerfile` if context is on-disk.
     * 
     * Conflicts with `inline`.
     * 
     */
    public Optional<Output<String>> location() {
        return Optional.ofNullable(this.location);
    }

    private DockerfileArgs() {}

    private DockerfileArgs(DockerfileArgs $) {
        this.inline = $.inline;
        this.location = $.location;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(DockerfileArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private DockerfileArgs $;

        public Builder() {
            $ = new DockerfileArgs();
        }

        public Builder(DockerfileArgs defaults) {
            $ = new DockerfileArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param inline Raw Dockerfile contents.
         * 
         * Conflicts with `location`.
         * 
         * Equivalent to invoking Docker with `-f -`.
         * 
         * @return builder
         * 
         */
        public Builder inline(@Nullable Output<String> inline) {
            $.inline = inline;
            return this;
        }

        /**
         * @param inline Raw Dockerfile contents.
         * 
         * Conflicts with `location`.
         * 
         * Equivalent to invoking Docker with `-f -`.
         * 
         * @return builder
         * 
         */
        public Builder inline(String inline) {
            return inline(Output.of(inline));
        }

        /**
         * @param location Location of the Dockerfile to use.
         * 
         * Can be a relative or absolute path to a local file, or a remote URL.
         * 
         * Defaults to `${context.location}/Dockerfile` if context is on-disk.
         * 
         * Conflicts with `inline`.
         * 
         * @return builder
         * 
         */
        public Builder location(@Nullable Output<String> location) {
            $.location = location;
            return this;
        }

        /**
         * @param location Location of the Dockerfile to use.
         * 
         * Can be a relative or absolute path to a local file, or a remote URL.
         * 
         * Defaults to `${context.location}/Dockerfile` if context is on-disk.
         * 
         * Conflicts with `inline`.
         * 
         * @return builder
         * 
         */
        public Builder location(String location) {
            return location(Output.of(location));
        }

        public DockerfileArgs build() {
            return $;
        }
    }

}
