// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.dockerbuild.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.dockerbuild.outputs.Context;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.String;
import java.util.Map;
import java.util.Objects;
import javax.annotation.Nullable;

@CustomType
public final class BuildContext {
    /**
     * @return Resources to use for build context.
     * 
     * The location can be:
     * * A relative or absolute path to a local directory (`.`, `./app`,
     *   `/app`, etc.).
     * * A remote URL of a Git repository, tarball, or plain text file
     *   (`https://github.com/user/myrepo.git`, `http://server/context.tar.gz`,
     *   etc.).
     * 
     */
    private String location;
    /**
     * @return Additional build contexts to use.
     * 
     * These contexts are accessed with `FROM name` or `--from=name`
     * statements when using Dockerfile 1.4+ syntax.
     * 
     * Values can be local paths, HTTP URLs, or  `docker-image://` images.
     * 
     */
    private @Nullable Map<String,Context> named;

    private BuildContext() {}
    /**
     * @return Resources to use for build context.
     * 
     * The location can be:
     * * A relative or absolute path to a local directory (`.`, `./app`,
     *   `/app`, etc.).
     * * A remote URL of a Git repository, tarball, or plain text file
     *   (`https://github.com/user/myrepo.git`, `http://server/context.tar.gz`,
     *   etc.).
     * 
     */
    public String location() {
        return this.location;
    }
    /**
     * @return Additional build contexts to use.
     * 
     * These contexts are accessed with `FROM name` or `--from=name`
     * statements when using Dockerfile 1.4+ syntax.
     * 
     * Values can be local paths, HTTP URLs, or  `docker-image://` images.
     * 
     */
    public Map<String,Context> named() {
        return this.named == null ? Map.of() : this.named;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(BuildContext defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private String location;
        private @Nullable Map<String,Context> named;
        public Builder() {}
        public Builder(BuildContext defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.location = defaults.location;
    	      this.named = defaults.named;
        }

        @CustomType.Setter
        public Builder location(String location) {
            if (location == null) {
              throw new MissingRequiredPropertyException("BuildContext", "location");
            }
            this.location = location;
            return this;
        }
        @CustomType.Setter
        public Builder named(@Nullable Map<String,Context> named) {

            this.named = named;
            return this;
        }
        public BuildContext build() {
            final var _resultValue = new BuildContext();
            _resultValue.location = location;
            _resultValue.named = named;
            return _resultValue;
        }
    }
}
