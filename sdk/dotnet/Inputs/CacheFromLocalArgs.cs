// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Dockerbuild.Inputs
{

    public sealed class CacheFromLocalArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// Digest of manifest to import.
        /// </summary>
        [Input("digest")]
        public Input<string>? Digest { get; set; }

        /// <summary>
        /// Path of the local directory where cache gets imported from.
        /// </summary>
        [Input("src", required: true)]
        public Input<string> Src { get; set; } = null!;

        public CacheFromLocalArgs()
        {
        }
        public static new CacheFromLocalArgs Empty => new CacheFromLocalArgs();
    }
}
