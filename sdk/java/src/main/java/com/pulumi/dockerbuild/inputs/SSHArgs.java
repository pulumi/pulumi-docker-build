// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.dockerbuild.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.String;
import java.util.List;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


public final class SSHArgs extends com.pulumi.resources.ResourceArgs {

    public static final SSHArgs Empty = new SSHArgs();

    /**
     * Useful for distinguishing different servers that are part of the same
     * build.
     * 
     * A value of `default` is appropriate if only dealing with a single host.
     * 
     */
    @Import(name="id", required=true)
    private Output<String> id;

    /**
     * @return Useful for distinguishing different servers that are part of the same
     * build.
     * 
     * A value of `default` is appropriate if only dealing with a single host.
     * 
     */
    public Output<String> id() {
        return this.id;
    }

    /**
     * SSH agent socket or private keys to expose to the build under the given
     * identifier.
     * 
     * Defaults to `[$SSH_AUTH_SOCK]`.
     * 
     * Note that your keys are **not** automatically added when using an
     * agent. Run `ssh-add -l` locally to confirm which public keys are
     * visible to the agent; these will be exposed to your build.
     * 
     */
    @Import(name="paths")
    private @Nullable Output<List<String>> paths;

    /**
     * @return SSH agent socket or private keys to expose to the build under the given
     * identifier.
     * 
     * Defaults to `[$SSH_AUTH_SOCK]`.
     * 
     * Note that your keys are **not** automatically added when using an
     * agent. Run `ssh-add -l` locally to confirm which public keys are
     * visible to the agent; these will be exposed to your build.
     * 
     */
    public Optional<Output<List<String>>> paths() {
        return Optional.ofNullable(this.paths);
    }

    private SSHArgs() {}

    private SSHArgs(SSHArgs $) {
        this.id = $.id;
        this.paths = $.paths;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(SSHArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private SSHArgs $;

        public Builder() {
            $ = new SSHArgs();
        }

        public Builder(SSHArgs defaults) {
            $ = new SSHArgs(Objects.requireNonNull(defaults));
        }

        /**
         * @param id Useful for distinguishing different servers that are part of the same
         * build.
         * 
         * A value of `default` is appropriate if only dealing with a single host.
         * 
         * @return builder
         * 
         */
        public Builder id(Output<String> id) {
            $.id = id;
            return this;
        }

        /**
         * @param id Useful for distinguishing different servers that are part of the same
         * build.
         * 
         * A value of `default` is appropriate if only dealing with a single host.
         * 
         * @return builder
         * 
         */
        public Builder id(String id) {
            return id(Output.of(id));
        }

        /**
         * @param paths SSH agent socket or private keys to expose to the build under the given
         * identifier.
         * 
         * Defaults to `[$SSH_AUTH_SOCK]`.
         * 
         * Note that your keys are **not** automatically added when using an
         * agent. Run `ssh-add -l` locally to confirm which public keys are
         * visible to the agent; these will be exposed to your build.
         * 
         * @return builder
         * 
         */
        public Builder paths(@Nullable Output<List<String>> paths) {
            $.paths = paths;
            return this;
        }

        /**
         * @param paths SSH agent socket or private keys to expose to the build under the given
         * identifier.
         * 
         * Defaults to `[$SSH_AUTH_SOCK]`.
         * 
         * Note that your keys are **not** automatically added when using an
         * agent. Run `ssh-add -l` locally to confirm which public keys are
         * visible to the agent; these will be exposed to your build.
         * 
         * @return builder
         * 
         */
        public Builder paths(List<String> paths) {
            return paths(Output.of(paths));
        }

        /**
         * @param paths SSH agent socket or private keys to expose to the build under the given
         * identifier.
         * 
         * Defaults to `[$SSH_AUTH_SOCK]`.
         * 
         * Note that your keys are **not** automatically added when using an
         * agent. Run `ssh-add -l` locally to confirm which public keys are
         * visible to the agent; these will be exposed to your build.
         * 
         * @return builder
         * 
         */
        public Builder paths(String... paths) {
            return paths(List.of(paths));
        }

        public SSHArgs build() {
            if ($.id == null) {
                throw new MissingRequiredPropertyException("SSHArgs", "id");
            }
            return $;
        }
    }

}
