// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Dockerbuild.Inputs
{

    public sealed class CacheToArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// Push cache to Azure's blob storage service.
        /// </summary>
        [Input("azblob")]
        public Input<Inputs.CacheToAzureBlobArgs>? Azblob { get; set; }

        /// <summary>
        /// When `true` this entry will be excluded. Defaults to `false`.
        /// </summary>
        [Input("disabled")]
        public Input<bool>? Disabled { get; set; }

        /// <summary>
        /// Recommended for use with GitHub Actions workflows.
        /// 
        /// An action like `crazy-max/ghaction-github-runtime` is recommended to
        /// expose appropriate credentials to your GitHub workflow.
        /// </summary>
        [Input("gha")]
        public Input<Inputs.CacheToGitHubActionsArgs>? Gha { get; set; }

        /// <summary>
        /// The inline cache storage backend is the simplest implementation to get
        /// started with, but it does not handle multi-stage builds. Consider the
        /// `registry` cache backend instead.
        /// </summary>
        [Input("inline")]
        public Input<Inputs.CacheToInlineArgs>? Inline { get; set; }

        /// <summary>
        /// A simple backend which caches imagines on your local filesystem.
        /// </summary>
        [Input("local")]
        public Input<Inputs.CacheToLocalArgs>? Local { get; set; }

        /// <summary>
        /// A raw string as you would provide it to the Docker CLI (e.g.,
        /// `type=inline`)
        /// </summary>
        [Input("raw")]
        public Input<string>? Raw { get; set; }

        /// <summary>
        /// Push caches to remote registries. Incompatible with the `docker` build
        /// driver.
        /// </summary>
        [Input("registry")]
        public Input<Inputs.CacheToRegistryArgs>? Registry { get; set; }

        /// <summary>
        /// Push cache to AWS S3 or S3-compatible services such as MinIO.
        /// </summary>
        [Input("s3")]
        public Input<Inputs.CacheToS3Args>? S3 { get; set; }

        public CacheToArgs()
        {
        }
        public static new CacheToArgs Empty => new CacheToArgs();
    }
}
