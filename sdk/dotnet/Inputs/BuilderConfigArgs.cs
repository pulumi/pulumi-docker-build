// *** WARNING: this file was generated by pulumi-language-dotnet. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.DockerBuild.Inputs
{

    public sealed class BuilderConfigArgs : global::Pulumi.ResourceArgs
    {
        /// <summary>
        /// Name of an existing buildx builder to use.
        /// 
        /// Only `docker-container`, `kubernetes`, or `remote` drivers are
        /// supported. The legacy `docker` driver is not supported.
        /// 
        /// Equivalent to Docker's `--builder` flag.
        /// </summary>
        [Input("name")]
        public Input<string>? Name { get; set; }

        public BuilderConfigArgs()
        {
        }
        public static new BuilderConfigArgs Empty => new BuilderConfigArgs();
    }
}
