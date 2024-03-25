using System.Collections.Generic;
using System.Linq;
using Pulumi;
using Dockerbuild = Pulumi.Dockerbuild;

return await Deployment.RunAsync(() => 
{
    var config = new Config();
    var dockerHubPassword = config.Require("dockerHubPassword");
    var multiPlatform = new Dockerbuild.Image("multiPlatform", new()
    {
        Dockerfile = new Dockerbuild.Inputs.DockerfileArgs
        {
            Location = "./app/Dockerfile.multiPlatform",
        },
        Context = new Dockerbuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
        Platforms = new[]
        {
            Dockerbuild.Platform.Plan9_amd64,
            Dockerbuild.Platform.Plan9_386,
        },
    });

    var registryPush = new Dockerbuild.Image("registryPush", new()
    {
        Context = new Dockerbuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
        Tags = new[]
        {
            "docker.io/pulumibot/buildkit-e2e:example",
        },
        Exports = new[]
        {
            new Dockerbuild.Inputs.ExportArgs
            {
                Registry = new Dockerbuild.Inputs.ExportRegistryArgs
                {
                    OciMediaTypes = true,
                    Push = false,
                },
            },
        },
        Registries = new[]
        {
            new Dockerbuild.Inputs.RegistryArgs
            {
                Address = "docker.io",
                Username = "pulumibot",
                Password = dockerHubPassword,
            },
        },
    });

    var cached = new Dockerbuild.Image("cached", new()
    {
        Context = new Dockerbuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
        CacheTo = new[]
        {
            new Dockerbuild.Inputs.CacheToArgs
            {
                Local = new Dockerbuild.Inputs.CacheToLocalArgs
                {
                    Dest = "tmp/cache",
                    Mode = Dockerbuild.CacheMode.Max,
                },
            },
        },
        CacheFrom = new[]
        {
            new Dockerbuild.Inputs.CacheFromArgs
            {
                Local = new Dockerbuild.Inputs.CacheFromLocalArgs
                {
                    Src = "tmp/cache",
                },
            },
        },
    });

    var buildArgs = new Dockerbuild.Image("buildArgs", new()
    {
        Dockerfile = new Dockerbuild.Inputs.DockerfileArgs
        {
            Location = "./app/Dockerfile.buildArgs",
        },
        Context = new Dockerbuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
        BuildArgs = 
        {
            { "SET_ME_TO_TRUE", "true" },
        },
    });

    var extraHosts = new Dockerbuild.Image("extraHosts", new()
    {
        Dockerfile = new Dockerbuild.Inputs.DockerfileArgs
        {
            Location = "./app/Dockerfile.extraHosts",
        },
        Context = new Dockerbuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
        AddHosts = new[]
        {
            "metadata.google.internal:169.254.169.254",
        },
    });

    var sshMount = new Dockerbuild.Image("sshMount", new()
    {
        Dockerfile = new Dockerbuild.Inputs.DockerfileArgs
        {
            Location = "./app/Dockerfile.sshMount",
        },
        Context = new Dockerbuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
        Ssh = new[]
        {
            new Dockerbuild.Inputs.SSHArgs
            {
                Id = "default",
            },
        },
    });

    var secrets = new Dockerbuild.Image("secrets", new()
    {
        Dockerfile = new Dockerbuild.Inputs.DockerfileArgs
        {
            Location = "./app/Dockerfile.secrets",
        },
        Context = new Dockerbuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
        Secrets = 
        {
            { "password", "hunter2" },
        },
    });

    var labels = new Dockerbuild.Image("labels", new()
    {
        Context = new Dockerbuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
        Labels = 
        {
            { "description", "This image will get a descriptive label 👍" },
        },
    });

    var target = new Dockerbuild.Image("target", new()
    {
        Dockerfile = new Dockerbuild.Inputs.DockerfileArgs
        {
            Location = "./app/Dockerfile.target",
        },
        Context = new Dockerbuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
        Target = "build-me",
    });

    var namedContexts = new Dockerbuild.Image("namedContexts", new()
    {
        Dockerfile = new Dockerbuild.Inputs.DockerfileArgs
        {
            Location = "./app/Dockerfile.namedContexts",
        },
        Context = new Dockerbuild.Inputs.BuildContextArgs
        {
            Location = "./app",
            Named = 
            {
                { "golang:latest", new Dockerbuild.Inputs.ContextArgs
                {
                    Location = "docker-image://golang@sha256:b8e62cf593cdaff36efd90aa3a37de268e6781a2e68c6610940c48f7cdf36984",
                } },
            },
        },
    });

    var remoteContext = new Dockerbuild.Image("remoteContext", new()
    {
        Context = new Dockerbuild.Inputs.BuildContextArgs
        {
            Location = "https://raw.githubusercontent.com/pulumi/pulumi-docker/api-types/provider/testdata/Dockerfile",
        },
    });

    var remoteContextWithInline = new Dockerbuild.Image("remoteContextWithInline", new()
    {
        Dockerfile = new Dockerbuild.Inputs.DockerfileArgs
        {
            Inline = @"FROM busybox
COPY hello.c ./
",
        },
        Context = new Dockerbuild.Inputs.BuildContextArgs
        {
            Location = "https://github.com/docker-library/hello-world.git",
        },
    });

    var inline = new Dockerbuild.Image("inline", new()
    {
        Dockerfile = new Dockerbuild.Inputs.DockerfileArgs
        {
            Inline = @"FROM alpine
RUN echo ""This uses an inline Dockerfile! 👍""
",
        },
        Context = new Dockerbuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
    });

    var dockerLoad = new Dockerbuild.Image("dockerLoad", new()
    {
        Context = new Dockerbuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
        Exports = new[]
        {
            new Dockerbuild.Inputs.ExportArgs
            {
                Docker = new Dockerbuild.Inputs.ExportDockerArgs
                {
                    Tar = true,
                },
            },
        },
    });

    return new Dictionary<string, object?>
    {
        ["platforms"] = multiPlatform.Platforms,
    };
});

