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
    public sealed class ExportTar
    {
        /// <summary>
        /// Output path.
        /// </summary>
        public readonly string Dest;

        [OutputConstructor]
        private ExportTar(string dest)
        {
            Dest = dest;
        }
    }
}
