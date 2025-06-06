// *** WARNING: this file was generated by pulumi-language-dotnet. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.DockerBuild.Outputs
{

    [OutputType]
    public sealed class Dockerfile
    {
        /// <summary>
        /// Raw Dockerfile contents.
        /// 
        /// Conflicts with `location`.
        /// 
        /// Equivalent to invoking Docker with `-f -`.
        /// </summary>
        public readonly string? Inline;
        /// <summary>
        /// Location of the Dockerfile to use.
        /// 
        /// Can be a relative or absolute path to a local file, or a remote URL.
        /// 
        /// Defaults to `${context.location}/Dockerfile` if context is on-disk.
        /// 
        /// Conflicts with `inline`.
        /// </summary>
        public readonly string? Location;

        [OutputConstructor]
        private Dockerfile(
            string? inline,

            string? location)
        {
            Inline = inline;
            Location = location;
        }
    }
}
