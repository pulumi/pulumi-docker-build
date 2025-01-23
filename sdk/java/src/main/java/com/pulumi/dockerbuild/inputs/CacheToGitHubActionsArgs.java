// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.dockerbuild.inputs;

import com.pulumi.core.Output;
import com.pulumi.core.annotations.Import;
import com.pulumi.core.internal.Codegen;
import com.pulumi.dockerbuild.enums.CacheMode;
import java.lang.Boolean;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;


public final class CacheToGitHubActionsArgs extends com.pulumi.resources.ResourceArgs {

    public static final CacheToGitHubActionsArgs Empty = new CacheToGitHubActionsArgs();

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
     * The scope to use for cache keys. Defaults to `buildkit`.
     * 
     * This should be set if building and caching multiple images in one
     * workflow, otherwise caches will overwrite each other.
     * 
     */
    @Import(name="scope")
    private @Nullable Output<String> scope;

    /**
     * @return The scope to use for cache keys. Defaults to `buildkit`.
     * 
     * This should be set if building and caching multiple images in one
     * workflow, otherwise caches will overwrite each other.
     * 
     */
    public Optional<Output<String>> scope() {
        return Optional.ofNullable(this.scope);
    }

    /**
     * The GitHub Actions token to use. This is not a personal access tokens
     * and is typically generated automatically as part of each job.
     * 
     * Defaults to `$ACTIONS_RUNTIME_TOKEN`, although a separate action like
     * `crazy-max/ghaction-github-runtime` is recommended to expose this
     * environment variable to your jobs.
     * 
     */
    @Import(name="token")
    private @Nullable Output<String> token;

    /**
     * @return The GitHub Actions token to use. This is not a personal access tokens
     * and is typically generated automatically as part of each job.
     * 
     * Defaults to `$ACTIONS_RUNTIME_TOKEN`, although a separate action like
     * `crazy-max/ghaction-github-runtime` is recommended to expose this
     * environment variable to your jobs.
     * 
     */
    public Optional<Output<String>> token() {
        return Optional.ofNullable(this.token);
    }

    /**
     * The cache server URL to use for artifacts.
     * 
     * Defaults to `$ACTIONS_CACHE_URL`, although a separate action like
     * `crazy-max/ghaction-github-runtime` is recommended to expose this
     * environment variable to your jobs.
     * 
     */
    @Import(name="url")
    private @Nullable Output<String> url;

    /**
     * @return The cache server URL to use for artifacts.
     * 
     * Defaults to `$ACTIONS_CACHE_URL`, although a separate action like
     * `crazy-max/ghaction-github-runtime` is recommended to expose this
     * environment variable to your jobs.
     * 
     */
    public Optional<Output<String>> url() {
        return Optional.ofNullable(this.url);
    }

    private CacheToGitHubActionsArgs() {}

    private CacheToGitHubActionsArgs(CacheToGitHubActionsArgs $) {
        this.ignoreError = $.ignoreError;
        this.mode = $.mode;
        this.scope = $.scope;
        this.token = $.token;
        this.url = $.url;
    }

    public static Builder builder() {
        return new Builder();
    }
    public static Builder builder(CacheToGitHubActionsArgs defaults) {
        return new Builder(defaults);
    }

    public static final class Builder {
        private CacheToGitHubActionsArgs $;

        public Builder() {
            $ = new CacheToGitHubActionsArgs();
        }

        public Builder(CacheToGitHubActionsArgs defaults) {
            $ = new CacheToGitHubActionsArgs(Objects.requireNonNull(defaults));
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
         * @param scope The scope to use for cache keys. Defaults to `buildkit`.
         * 
         * This should be set if building and caching multiple images in one
         * workflow, otherwise caches will overwrite each other.
         * 
         * @return builder
         * 
         */
        public Builder scope(@Nullable Output<String> scope) {
            $.scope = scope;
            return this;
        }

        /**
         * @param scope The scope to use for cache keys. Defaults to `buildkit`.
         * 
         * This should be set if building and caching multiple images in one
         * workflow, otherwise caches will overwrite each other.
         * 
         * @return builder
         * 
         */
        public Builder scope(String scope) {
            return scope(Output.of(scope));
        }

        /**
         * @param token The GitHub Actions token to use. This is not a personal access tokens
         * and is typically generated automatically as part of each job.
         * 
         * Defaults to `$ACTIONS_RUNTIME_TOKEN`, although a separate action like
         * `crazy-max/ghaction-github-runtime` is recommended to expose this
         * environment variable to your jobs.
         * 
         * @return builder
         * 
         */
        public Builder token(@Nullable Output<String> token) {
            $.token = token;
            return this;
        }

        /**
         * @param token The GitHub Actions token to use. This is not a personal access tokens
         * and is typically generated automatically as part of each job.
         * 
         * Defaults to `$ACTIONS_RUNTIME_TOKEN`, although a separate action like
         * `crazy-max/ghaction-github-runtime` is recommended to expose this
         * environment variable to your jobs.
         * 
         * @return builder
         * 
         */
        public Builder token(String token) {
            return token(Output.of(token));
        }

        /**
         * @param url The cache server URL to use for artifacts.
         * 
         * Defaults to `$ACTIONS_CACHE_URL`, although a separate action like
         * `crazy-max/ghaction-github-runtime` is recommended to expose this
         * environment variable to your jobs.
         * 
         * @return builder
         * 
         */
        public Builder url(@Nullable Output<String> url) {
            $.url = url;
            return this;
        }

        /**
         * @param url The cache server URL to use for artifacts.
         * 
         * Defaults to `$ACTIONS_CACHE_URL`, although a separate action like
         * `crazy-max/ghaction-github-runtime` is recommended to expose this
         * environment variable to your jobs.
         * 
         * @return builder
         * 
         */
        public Builder url(String url) {
            return url(Output.of(url));
        }

        public CacheToGitHubActionsArgs build() {
            $.ignoreError = Codegen.booleanProp("ignoreError").output().arg($.ignoreError).def(false).getNullable();
            $.mode = Codegen.objectProp("mode", CacheMode.class).output().arg($.mode).def(CacheMode.Min).getNullable();
            $.scope = Codegen.stringProp("scope").output().arg($.scope).env("buildkit").def("").getNullable();
            $.token = Codegen.stringProp("token").secret().arg($.token).env("ACTIONS_RUNTIME_TOKEN").def("").getNullable();
            $.url = Codegen.stringProp("url").output().arg($.url).env("ACTIONS_CACHE_URL").def("").getNullable();
            return $;
        }
    }

}
