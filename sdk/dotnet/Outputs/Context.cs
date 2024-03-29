// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Dockerbuild.Outputs
{

    [OutputType]
    public sealed class Context
    {
        /// <summary>
        /// Resources to use for build context.
        /// 
        /// The location can be:
        /// * A relative or absolute path to a local directory (`.`, `./app`,
        ///   `/app`, etc.).
        /// * A remote URL of a Git repository, tarball, or plain text file
        ///   (`https://github.com/user/myrepo.git`, `http://server/context.tar.gz`,
        ///   etc.).
        /// </summary>
        public readonly string Location;

        [OutputConstructor]
        private Context(string location)
        {
            Location = location;
        }
    }
}