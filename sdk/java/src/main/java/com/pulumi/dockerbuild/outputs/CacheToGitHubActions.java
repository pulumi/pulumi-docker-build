// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.dockerbuild.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.dockerbuild.enums.CacheMode;
import java.lang.Boolean;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class CacheToGitHubActions {
    /**
     * @return Ignore errors caused by failed cache exports.
     * 
     */
    private @Nullable Boolean ignoreError;
    /**
     * @return The cache mode to use. Defaults to `min`.
     * 
     */
    private @Nullable CacheMode mode;
    /**
     * @return The scope to use for cache keys. Defaults to `buildkit`.
     * 
     * This should be set if building and caching multiple images in one
     * workflow, otherwise caches will overwrite each other.
     * 
     */
    private @Nullable String scope;
    /**
     * @return The GitHub Actions token to use. This is not a personal access tokens
     * and is typically generated automatically as part of each job.
     * 
     * Defaults to `$ACTIONS_RUNTIME_TOKEN`, although a separate action like
     * `crazy-max/ghaction-github-runtime` is recommended to expose this
     * environment variable to your jobs.
     * 
     */
    private @Nullable String token;
    /**
     * @return The cache server URL to use for artifacts.
     * 
     * Defaults to `$ACTIONS_CACHE_URL`, although a separate action like
     * `crazy-max/ghaction-github-runtime` is recommended to expose this
     * environment variable to your jobs.
     * 
     */
    private @Nullable String url;

    private CacheToGitHubActions() {}
    /**
     * @return Ignore errors caused by failed cache exports.
     * 
     */
    public Optional<Boolean> ignoreError() {
        return Optional.ofNullable(this.ignoreError);
    }
    /**
     * @return The cache mode to use. Defaults to `min`.
     * 
     */
    public Optional<CacheMode> mode() {
        return Optional.ofNullable(this.mode);
    }
    /**
     * @return The scope to use for cache keys. Defaults to `buildkit`.
     * 
     * This should be set if building and caching multiple images in one
     * workflow, otherwise caches will overwrite each other.
     * 
     */
    public Optional<String> scope() {
        return Optional.ofNullable(this.scope);
    }
    /**
     * @return The GitHub Actions token to use. This is not a personal access tokens
     * and is typically generated automatically as part of each job.
     * 
     * Defaults to `$ACTIONS_RUNTIME_TOKEN`, although a separate action like
     * `crazy-max/ghaction-github-runtime` is recommended to expose this
     * environment variable to your jobs.
     * 
     */
    public Optional<String> token() {
        return Optional.ofNullable(this.token);
    }
    /**
     * @return The cache server URL to use for artifacts.
     * 
     * Defaults to `$ACTIONS_CACHE_URL`, although a separate action like
     * `crazy-max/ghaction-github-runtime` is recommended to expose this
     * environment variable to your jobs.
     * 
     */
    public Optional<String> url() {
        return Optional.ofNullable(this.url);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(CacheToGitHubActions defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable Boolean ignoreError;
        private @Nullable CacheMode mode;
        private @Nullable String scope;
        private @Nullable String token;
        private @Nullable String url;
        public Builder() {}
        public Builder(CacheToGitHubActions defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.ignoreError = defaults.ignoreError;
    	      this.mode = defaults.mode;
    	      this.scope = defaults.scope;
    	      this.token = defaults.token;
    	      this.url = defaults.url;
        }

        @CustomType.Setter
        public Builder ignoreError(@Nullable Boolean ignoreError) {

            this.ignoreError = ignoreError;
            return this;
        }
        @CustomType.Setter
        public Builder mode(@Nullable CacheMode mode) {

            this.mode = mode;
            return this;
        }
        @CustomType.Setter
        public Builder scope(@Nullable String scope) {

            this.scope = scope;
            return this;
        }
        @CustomType.Setter
        public Builder token(@Nullable String token) {

            this.token = token;
            return this;
        }
        @CustomType.Setter
        public Builder url(@Nullable String url) {

            this.url = url;
            return this;
        }
        public CacheToGitHubActions build() {
            final var _resultValue = new CacheToGitHubActions();
            _resultValue.ignoreError = ignoreError;
            _resultValue.mode = mode;
            _resultValue.scope = scope;
            _resultValue.token = token;
            _resultValue.url = url;
            return _resultValue;
        }
    }
}
