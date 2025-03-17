// *** WARNING: this file was generated by pulumi-language-java. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

package com.pulumi.dockerbuild.outputs;

import com.pulumi.core.annotations.CustomType;
import com.pulumi.exceptions.MissingRequiredPropertyException;
import java.lang.String;
import java.util.Objects;

@CustomType
public final class Context {
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

    private Context() {}
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

    public static Builder builder() {
        return new Builder();
    }

    public static Builder builder(Context defaults) {
        return new Builder(defaults);
    }
    @CustomType.Builder
    public static final class Builder {
        private String location;
        public Builder() {}
        public Builder(Context defaults) {
    	      Objects.requireNonNull(defaults);
    	      this.location = defaults.location;
        }

        @CustomType.Setter
        public Builder location(String location) {
            if (location == null) {
              throw new MissingRequiredPropertyException("Context", "location");
            }
            this.location = location;
            return this;
        }
        public Context build() {
            final var _resultValue = new Context();
            _resultValue.location = location;
            return _resultValue;
        }
    }
}
