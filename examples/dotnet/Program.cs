using System.Collections.Generic;
using System.Linq;
using Pulumi;
using Dockerbuild = Pulumi.Dockerbuild;

return await Deployment.RunAsync(() => 
{
    var myRandomResource = new Dockerbuild.Random("myRandomResource", new()
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

