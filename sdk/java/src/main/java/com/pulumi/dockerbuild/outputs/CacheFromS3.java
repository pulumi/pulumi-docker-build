// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.dockerbuild.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.Boolean;
import java.lang.String;
import java.util.Objects;
import java.util.Optional;
import javax.annotation.Nullable;

@CustomType
public final class CacheFromS3 {
    /**
     * @return Defaults to `$AWS_ACCESS_KEY_ID`.
     * 
     */
    private @Nullable String accessKeyId;
    /**
     * @return Prefix to prepend to blob filenames.
     * 
     */
    private @Nullable String blobsPrefix;
    /**
     * @return Name of the S3 bucket.
     * 
     */
    private String bucket;
    /**
     * @return Endpoint of the S3 bucket.
     * 
     */
    private @Nullable String endpointUrl;
    /**
     * @return Prefix to prepend on manifest filenames.
     * 
     */
    private @Nullable String manifestsPrefix;
    /**
     * @return Name of the cache image.
     * 
     */
    private @Nullable String name;
    /**
     * @return The geographic location of the bucket. Defaults to `$AWS_REGION`.
     * 
     */
    private String region;
    /**
     * @return Defaults to `$AWS_SECRET_ACCESS_KEY`.
     * 
     */
    private @Nullable String secretAccessKey;
    /**
     * @return Defaults to `$AWS_SESSION_TOKEN`.
     * 
     */
    private @Nullable String sessionToken;
    /**
     * @return Uses `bucket` in the URL instead of hostname when `true`.
     * 
     */
    private @Nullable Boolean usePathStyle;

    private CacheFromS3() {}
    /**
     * @return Defaults to `$AWS_ACCESS_KEY_ID`.
     * 
     */
    public Optional<String> accessKeyId() {
        return Optional.ofNullable(this.accessKeyId);
    }
    /**
     * @return Prefix to prepend to blob filenames.
     * 
     */
    public Optional<String> blobsPrefix() {
        return Optional.ofNullable(this.blobsPrefix);
    }
    /**
     * @return Name of the S3 bucket.
     * 
     */
    public String bucket() {
        return this.bucket;
    }
    /**
     * @return Endpoint of the S3 bucket.
     * 
     */
    public Optional<String> endpointUrl() {
        return Optional.ofNullable(this.endpointUrl);
    }
    /**
     * @return Prefix to prepend on manifest filenames.
     * 
     */
    public Optional<String> manifestsPrefix() {
        return Optional.ofNullable(this.manifestsPrefix);
    }
    /**
     * @return Name of the cache image.
     * 
     */
    public Optional<String> name() {
        return Optional.ofNullable(this.name);
    }
    /**
     * @return The geographic location of the bucket. Defaults to `$AWS_REGION`.
     * 
     */
    public String region() {
        return this.region;
    }
    /**
     * @return Defaults to `$AWS_SECRET_ACCESS_KEY`.
     * 
     */
    public Optional<String> secretAccessKey() {
        return Optional.ofNullable(this.secretAccessKey);
    }
    /**
     * @return Defaults to `$AWS_SESSION_TOKEN`.
     * 
     */
    public Optional<String> sessionToken() {
        return Optional.ofNullable(this.sessionToken);
    }
    /**
     * @return Uses `bucket` in the URL instead of hostname when `true`.
     * 
     */
    public Optional<Boolean> usePathStyle() {
        return Optional.ofNullable(this.usePathStyle);
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(CacheFromS3 defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private @Nullable String accessKeyId;
        private @Nullable String blobsPrefix;
        private String bucket;
        private @Nullable String endpointUrl;
        private @Nullable String manifestsPrefix;
        private @Nullable String name;
        private String region;
        private @Nullable String secretAccessKey;
        private @Nullable String sessionToken;
        private @Nullable Boolean usePathStyle;
        public Builder() {}
        public Builder(CacheFromS3 defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.accessKeyId = defaults.accessKeyId;
    	      this.blobsPrefix = defaults.blobsPrefix;
    	      this.bucket = defaults.bucket;
    	      this.endpointUrl = defaults.endpointUrl;
    	      this.manifestsPrefix = defaults.manifestsPrefix;
    	      this.name = defaults.name;
    	      this.region = defaults.region;
    	      this.secretAccessKey = defaults.secretAccessKey;
    	      this.sessionToken = defaults.sessionToken;
    	      this.usePathStyle = defaults.usePathStyle;
        }

        @CustomType.Setter
        public Builder accessKeyId(@Nullable String accessKeyId) {

            this.accessKeyId = accessKeyId;
            return this;
        }
        @CustomType.Setter
        public Builder blobsPrefix(@Nullable String blobsPrefix) {

            this.blobsPrefix = blobsPrefix;
            return this;
        }
        @CustomType.Setter
        public Builder bucket(String bucket) {
            if (bucket == null) {
              throw new MissingRequiredPropertyException("CacheFromS3", "bucket");
            }
            this.bucket = bucket;
            return this;
        }
        @CustomType.Setter
        public Builder endpointUrl(@Nullable String endpointUrl) {

            this.endpointUrl = endpointUrl;
            return this;
        }
        @CustomType.Setter
        public Builder manifestsPrefix(@Nullable String manifestsPrefix) {

            this.manifestsPrefix = manifestsPrefix;
            return this;
        }
        @CustomType.Setter
        public Builder name(@Nullable String name) {

            this.name = name;
            return this;
        }
        @CustomType.Setter
        public Builder region(String region) {
            if (region == null) {
              throw new MissingRequiredPropertyException("CacheFromS3", "region");
            }
            this.region = region;
            return this;
        }
        @CustomType.Setter
        public Builder secretAccessKey(@Nullable String secretAccessKey) {

            this.secretAccessKey = secretAccessKey;
            return this;
        }
        @CustomType.Setter
        public Builder sessionToken(@Nullable String sessionToken) {

            this.sessionToken = sessionToken;
            return this;
        }
        @CustomType.Setter
        public Builder usePathStyle(@Nullable Boolean usePathStyle) {

            this.usePathStyle = usePathStyle;
            return this;
        }
        public CacheFromS3 build() {
            final var _resultValue = new CacheFromS3();
            _resultValue.accessKeyId = accessKeyId;
            _resultValue.blobsPrefix = blobsPrefix;
            _resultValue.bucket = bucket;
            _resultValue.endpointUrl = endpointUrl;
            _resultValue.manifestsPrefix = manifestsPrefix;
            _resultValue.name = name;
            _resultValue.region = region;
            _resultValue.secretAccessKey = secretAccessKey;
            _resultValue.sessionToken = sessionToken;
            _resultValue.usePathStyle = usePathStyle;
            return _resultValue;
        }
    }
}