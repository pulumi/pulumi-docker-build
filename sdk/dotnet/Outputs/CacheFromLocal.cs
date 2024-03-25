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
    public sealed class CacheFromLocal
    {
        /// <summary>
        /// Digest of manifest to import.
        /// </summary>
        public readonly string? Digest;
        /// <summary>
        /// Path of the local directory where cache gets imported from.
        /// </summary>
        public readonly string Src;

        [OutputConstructor]
        private CacheFromLocal(
            string? digest,

            string src)
        {
            Digest = digest;
            Src = src;
        }
    }
}
