using System.Collections.Generic;
using System.Linq;
using Pulumi;
using DockerNative = Pulumi.DockerNative;

return await Deployment.RunAsync(() => 
{
    var myRandomResource = new DockerNative.Random("myRandomResource", new()
    {
        Length = 24,
    });

    return new Dictionary<string, object?>
    {
        ["output"] = 
        {
            { "value", myRandomResource.Result },
        },
    };
});

