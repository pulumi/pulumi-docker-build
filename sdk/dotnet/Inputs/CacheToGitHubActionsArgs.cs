// *** WARNING: this file was generated by pulumi. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Dockerbuild.Inputs
{

    public sealed class CacheToGitHubActionsArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// Ignore errors caused by failed cache exports.
        /// </summary>
        [Input("ignoreError")]
        public Input<bool>? IgnoreError { get; set; }

        /// <summary>
        /// The cache mode to use. Defaults to `min`.
        /// </summary>
        [Input("mode")]
        public Input<Pulumi.Dockerbuild.CacheMode>? Mode { get; set; }

        /// <summary>
        /// The scope to use for cache keys. Defaults to `buildkit`.
        /// 
        /// This should be set if building and caching multiple images in one
        /// workflow, otherwise caches will overwrite each other.
        /// </summary>
        [Input("scope")]
        public Input<string>? Scope { get; set; }

        [Input("token")]
        private Input<string>? _token;

        /// <summary>
        /// The GitHub Actions token to use. This is not a personal access tokens
        /// and is typically generated automatically as part of each job.
        /// 
        /// Defaults to `$ACTIONS_RUNTIME_TOKEN`, although a separate action like
        /// `crazy-max/ghaction-github-runtime` is recommended to expose this
        /// environment variable to your jobs.
        /// </summary>
        public Input<string>? Token
        {
            get => _token;
            set
            {
                var emptySecret = Output.CreateSecret(0);
                _token = Output.Tuple<Input<string>?, int>(value, emptySecret).Apply(t => t.Item1);
            }
        }

        /// <summary>
        /// The cache server URL to use for artifacts.
        /// 
        /// Defaults to `$ACTIONS_RUNTIME_URL`, although a separate action like
        /// `crazy-max/ghaction-github-runtime` is recommended to expose this
        /// environment variable to your jobs.
        /// </summary>
        [Input("url")]
        public Input<string>? Url { get; set; }

        public CacheToGitHubActionsArgs()
        {
            IgnoreError = false;
            Mode = Pulumi.Dockerbuild.CacheMode.Min;
            Scope = Utilities.GetEnv("buildkit") ?? "";
            Token = Utilities.GetEnv("ACTIONS_RUNTIME_TOKEN") ?? "";
            Url = Utilities.GetEnv("ACTIONS_RUNTIME_URL") ?? "";
        }
        public static new CacheToGitHubActionsArgs Empty => new CacheToGitHubActionsArgs();
    }
}