// *** WARNING: this file was generated by pulumi-language-dotnet. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Immutable;

namespace Pulumi.DockerBuild
{
    public static class Config
    {
        [global::System.Diagnostics.CodeAnalysis.SuppressMessage("Microsoft.Design", "IDE1006", Justification = 
        "Double underscore prefix used to avoid conflicts with variable names.")]
        private sealed class __Value<T>
        {
            private readonly Func<T> _getter;
            private T _value = default!;
            private bool _set;

            public __Value(Func<T> getter)
            {
                _getter = getter;
            }

            public T Get() => _set ? _value : _getter();

            public void Set(T value)
            {
                _value = value;
                _set = true;
            }
        }

        private static readonly global::Pulumi.Config __config = new global::Pulumi.Config("docker-build");

        private static readonly __Value<string?> _host = new __Value<string?>(() => __config.Get("host") ?? Utilities.GetEnv("DOCKER_HOST") ?? "");
        /// <summary>
        /// The build daemon's address.
        /// </summary>
        public static string? Host
        {
            get => _host.Get();
            set => _host.Set(value);
        }

        private static readonly __Value<ImmutableArray<Types.Registry>> _registries = new __Value<ImmutableArray<Types.Registry>>(() => __config.GetObject<ImmutableArray<Types.Registry>>("registries"));
        public static ImmutableArray<Types.Registry> Registries
        {
            get => _registries.Get();
            set => _registries.Set(value);
        }

        public static class Types
        {

             public class Registry
             {
            /// <summary>
            /// The registry's address (e.g. "docker.io").
            /// </summary>
                public string Address { get; set; }
            /// <summary>
            /// Password or token for the registry.
            /// </summary>
                public string? Password { get; set; } = null!;
            /// <summary>
            /// Username for the registry.
            /// </summary>
                public string? Username { get; set; } = null!;
            }
        }
    }
}
