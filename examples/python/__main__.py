import pulumi
import pulumi_dockerbuild as dockerbuild

my_random_resource = dockerbuild.Random("myRandomResource", length=24)
pulumi.export("value", my_random_resource.result)
