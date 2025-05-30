// *** WARNING: this file was generated by pulumi-language-dotnet. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.DockerBuild.Inputs
{

    public sealed class CacheToAzureBlobArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// Base URL of the storage account.
        /// </summary>
        [Input("accountUrl")]
        public Input<string>? AccountUrl { get; set; }

        /// <summary>
        /// Ignore errors caused by failed cache exports.
        /// </summary>
        [Input("ignoreError")]
        public Input<bool>? IgnoreError { get; set; }

        /// <summary>
        /// The cache mode to use. Defaults to `min`.
        /// </summary>
        [Input("mode")]
        public Input<Pulumi.DockerBuild.CacheMode>? Mode { get; set; }

        /// <summary>
        /// The name of the cache image.
        /// </summary>
        [Input("name", required: true)]
        public Input<string> Name { get; set; } = null!;

        [Input("secretAccessKey")]
        private Input<string>? _secretAccessKey;

        /// <summary>
        /// Blob storage account key.
        /// </summary>
        public Input<string>? SecretAccessKey
        {
            get => _secretAccessKey;
            set
            {
                var emptySecret = Output.CreateSecret(0);
                _secretAccessKey = Output.Tuple<Input<string>?, int>(value, emptySecret).Apply(t => t.Item1);
            }
        }

        public CacheToAzureBlobArgs()
        {
            IgnoreError = false;
            Mode = Pulumi.DockerBuild.CacheMode.Min;
        }
        public static new CacheToAzureBlobArgs Empty => new CacheToAzureBlobArgs();
    }
}
