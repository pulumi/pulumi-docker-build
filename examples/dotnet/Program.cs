using System.Collections.Generic;
using System.Linq;
using Pulumi;
using DockerBuild = Pulumi.DockerBuild;

return await Deployment.RunAsync(() => 
{
    var config = new Config();
    var dockerHubPassword = config.Require("dockerHubPassword");
    var multiPlatform = new DockerBuild.Image("multiPlatform", new()
    {
        Dockerfile = new DockerBuild.Inputs.DockerfileArgs
        {
            Location = "./app/Dockerfile.multiPlatform",
        },
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
        Platforms = new[]
        {
            DockerBuild.Platform.Plan9_amd64,
            DockerBuild.Platform.Plan9_386,
        },
    });

    var registryPush = new DockerBuild.Image("registryPush", new()
    {
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
        Tags = new[]
        {
            "docker.io/pulumibot/buildkit-e2e:example",
        },
        Exports = new[]
        {
            new DockerBuild.Inputs.ExportArgs
            {
                Registry = new DockerBuild.Inputs.ExportRegistryArgs
                {
                    OciMediaTypes = true,
                    Push = false,
                },
            },
        },
        Registries = new[]
        {
            new DockerBuild.Inputs.RegistryArgs
            {
                Address = "docker.io",
                Username = "pulumibot",
                Password = dockerHubPassword,
            },
        },
    });

    var cached = new DockerBuild.Image("cached", new()
    {
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
        CacheTo = new[]
        {
            new DockerBuild.Inputs.CacheToArgs
            {
                Local = new DockerBuild.Inputs.CacheToLocalArgs
                {
                    Dest = "tmp/cache",
                    Mode = DockerBuild.CacheMode.Max,
                },
            },
        },
        CacheFrom = new[]
        {
            new DockerBuild.Inputs.CacheFromArgs
            {
                Local = new DockerBuild.Inputs.CacheFromLocalArgs
                {
                    Src = "tmp/cache",
                },
            },
        },
    });

    var buildArgs = new DockerBuild.Image("buildArgs", new()
    {
        Dockerfile = new DockerBuild.Inputs.DockerfileArgs
        {
            Location = "./app/Dockerfile.buildArgs",
        },
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
        BuildArgs = 
        {
            { "SET_ME_TO_TRUE", "true" },
        },
    });

    var extraHosts = new DockerBuild.Image("extraHosts", new()
    {
        Dockerfile = new DockerBuild.Inputs.DockerfileArgs
        {
            Location = "./app/Dockerfile.extraHosts",
        },
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
        AddHosts = new[]
        {
            "metadata.google.internal:169.254.169.254",
        },
    });

    var sshMount = new DockerBuild.Image("sshMount", new()
    {
        Dockerfile = new DockerBuild.Inputs.DockerfileArgs
        {
            Location = "./app/Dockerfile.sshMount",
        },
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
        Ssh = new[]
        {
            new DockerBuild.Inputs.SSHArgs
            {
                Id = "default",
            },
        },
    });

    var secrets = new DockerBuild.Image("secrets", new()
    {
        Dockerfile = new DockerBuild.Inputs.DockerfileArgs
        {
            Location = "./app/Dockerfile.secrets",
        },
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
        Secrets = 
        {
            { "password", "hunter2" },
        },
    });

    var labels = new DockerBuild.Image("labels", new()
    {
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
        Labels = 
        {
            { "description", "This image will get a descriptive label 👍" },
        },
    });

    var target = new DockerBuild.Image("target", new()
    {
        Dockerfile = new DockerBuild.Inputs.DockerfileArgs
        {
            Location = "./app/Dockerfile.target",
        },
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
        Target = "build-me",
    });

    var namedContexts = new DockerBuild.Image("namedContexts", new()
    {
        Dockerfile = new DockerBuild.Inputs.DockerfileArgs
        {
            Location = "./app/Dockerfile.namedContexts",
        },
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "./app",
            Named = 
            {
                { "golang:latest", new DockerBuild.Inputs.ContextArgs
                {
                    Location = "docker-image://golang@sha256:b8e62cf593cdaff36efd90aa3a37de268e6781a2e68c6610940c48f7cdf36984",
                } },
            },
        },
    });

    var remoteContext = new DockerBuild.Image("remoteContext", new()
    {
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "https://raw.githubusercontent.com/pulumi/pulumi-docker/api-types/provider/testdata/Dockerfile",
        },
    });

    var remoteContextWithInline = new DockerBuild.Image("remoteContextWithInline", new()
    {
        Dockerfile = new DockerBuild.Inputs.DockerfileArgs
        {
            Inline = @"FROM busybox
COPY hello.c ./
",
        },
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "https://github.com/docker-library/hello-world.git",
        },
    });

    var inline = new DockerBuild.Image("inline", new()
    {
        Dockerfile = new DockerBuild.Inputs.DockerfileArgs
        {
            Inline = @"FROM alpine
RUN echo ""This uses an inline Dockerfile! 👍""
",
        },
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
    });

    var dockerLoad = new DockerBuild.Image("dockerLoad", new()
    {
        Context = new DockerBuild.Inputs.BuildContextArgs
        {
            Location = "./app",
        },
        Exports = new[]
        {
            new DockerBuild.Inputs.ExportArgs
            {
                Docker = new DockerBuild.Inputs.ExportDockerArgs
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

