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
    public sealed class CacheFromGitHubActions
    {
        /// <summary>
        /// The scope to use for cache keys. Defaults to `buildkit`.
        /// 
        /// This should be set if building and caching multiple images in one
        /// workflow, otherwise caches will overwrite each other.
        /// </summary>
        public readonly string? Scope;
        /// <summary>
        /// The GitHub Actions token to use. This is not a personal access tokens
        /// and is typically generated automatically as part of each job.
        /// 
        /// Defaults to `$ACTIONS_RUNTIME_TOKEN`, although a separate action like
        /// `crazy-max/ghaction-github-runtime` is recommended to expose this
        /// environment variable to your jobs.
        /// </summary>
        public readonly string? Token;
        /// <summary>
        /// The cache server URL to use for artifacts.
        /// 
        /// Defaults to `$ACTIONS_RUNTIME_URL`, although a separate action like
        /// `crazy-max/ghaction-github-runtime` is recommended to expose this
        /// environment variable to your jobs.
        /// </summary>
        public readonly string? Url;

        [OutputConstructor]
        private CacheFromGitHubActions(
            string? scope,

            string? token,

            string? url)
        {
            Scope = scope;
            Token = token;
            Url = url;
        }
    }
}
