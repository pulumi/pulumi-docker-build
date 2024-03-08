import pulumi
import pulumi_docker_native as docker_native

my_random_resource = docker_native.Random("myRandomResource", length=24)
pulumi.export("output", {
    "value": my_random_resource.result,
})
