# coding=utf-8
# *** WARNING: this file was generated by pulumi-language-python. ***
# *** Do not edit by hand unless you're certain you know what you are doing! ***

import copy
import warnings
import sys
import pulumi
import pulumi.runtime
from typing import Any, Mapping, Optional, Sequence, Union, overload
if sys.version_info >= (3, 11):
    from typing import NotRequired, TypedDict, TypeAlias
else:
    from typing_extensions import NotRequired, TypedDict, TypeAlias
from . import _utilities
from . import outputs
from ._inputs import *

__all__ = ['IndexArgs', 'Index']

@pulumi.input_type
class IndexArgs:
    def __init__(__self__, *,
                 sources: pulumi.Input[Sequence[pulumi.Input[str]]],
                 tag: pulumi.Input[str],
                 push: Optional[pulumi.Input[bool]] = None,
                 registry: Optional[pulumi.Input['RegistryArgs']] = None):
        """
        The set of arguments for constructing a Index resource.
        :param pulumi.Input[Sequence[pulumi.Input[str]]] sources: Existing images to include in the index.
        :param pulumi.Input[str] tag: The tag to apply to the index.
        :param pulumi.Input[bool] push: If true, push the index to the target registry.
               
               Defaults to `true`.
        :param pulumi.Input['RegistryArgs'] registry: Authentication for the registry where the tagged index will be pushed.
               
               Credentials can also be included with the provider's configuration.
        """
        pulumi.set(__self__, "sources", sources)
        pulumi.set(__self__, "tag", tag)
        if push is None:
            push = True
        if push is not None:
            pulumi.set(__self__, "push", push)
        if registry is not None:
            pulumi.set(__self__, "registry", registry)

    @property
    @pulumi.getter
    def sources(self) -> pulumi.Input[Sequence[pulumi.Input[str]]]:
        """
        Existing images to include in the index.
        """
        return pulumi.get(self, "sources")

    @sources.setter
    def sources(self, value: pulumi.Input[Sequence[pulumi.Input[str]]]):
        pulumi.set(self, "sources", value)

    @property
    @pulumi.getter
    def tag(self) -> pulumi.Input[str]:
        """
        The tag to apply to the index.
        """
        return pulumi.get(self, "tag")

    @tag.setter
    def tag(self, value: pulumi.Input[str]):
        pulumi.set(self, "tag", value)

    @property
    @pulumi.getter
    def push(self) -> Optional[pulumi.Input[bool]]:
        """
        If true, push the index to the target registry.

        Defaults to `true`.
        """
        return pulumi.get(self, "push")

    @push.setter
    def push(self, value: Optional[pulumi.Input[bool]]):
        pulumi.set(self, "push", value)

    @property
    @pulumi.getter
    def registry(self) -> Optional[pulumi.Input['RegistryArgs']]:
        """
        Authentication for the registry where the tagged index will be pushed.

        Credentials can also be included with the provider's configuration.
        """
        return pulumi.get(self, "registry")

    @registry.setter
    def registry(self, value: Optional[pulumi.Input['RegistryArgs']]):
        pulumi.set(self, "registry", value)


class Index(pulumi.CustomResource):
    @overload
    def __init__(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 push: Optional[pulumi.Input[bool]] = None,
                 registry: Optional[pulumi.Input[Union['RegistryArgs', 'RegistryArgsDict']]] = None,
                 sources: Optional[pulumi.Input[Sequence[pulumi.Input[str]]]] = None,
                 tag: Optional[pulumi.Input[str]] = None,
                 __props__=None):
        """
        A wrapper around `docker buildx imagetools create` to create an index
        (or manifest list) referencing one or more existing images.

        In most cases you do not need an `Index` to build a multi-platform
        image -- specifying multiple platforms on the `Image` will handle this
        for you automatically.

        However, as of April 2024, building multi-platform images _with
        caching_ will only export a cache for one platform at a time (see [this
        discussion](https://github.com/docker/buildx/discussions/1382) for more
        details).

        Therefore this resource can be helpful if you are building
        multi-platform images with caching: each platform can be built and
        cached separately, and an `Index` can join them all together. An
        example of this is shown below.

        This resource creates an OCI image index or a Docker manifest list
        depending on the media types of the source images.

        ## Example Usage
        ### Multi-platform registry caching
        ```python
        import pulumi
        import pulumi_docker_build as docker_build

        amd64 = docker_build.Image("amd64",
            cache_from=[{
                "registry": {
                    "ref": "docker.io/pulumi/pulumi:cache-amd64",
                },
            }],
            cache_to=[{
                "registry": {
                    "mode": docker_build.CacheMode.MAX,
                    "ref": "docker.io/pulumi/pulumi:cache-amd64",
                },
            }],
            context={
                "location": "app",
            },
            platforms=[docker_build.Platform.LINUX_AMD64],
            tags=["docker.io/pulumi/pulumi:3.107.0-amd64"])
        arm64 = docker_build.Image("arm64",
            cache_from=[{
                "registry": {
                    "ref": "docker.io/pulumi/pulumi:cache-arm64",
                },
            }],
            cache_to=[{
                "registry": {
                    "mode": docker_build.CacheMode.MAX,
                    "ref": "docker.io/pulumi/pulumi:cache-arm64",
                },
            }],
            context={
                "location": "app",
            },
            platforms=[docker_build.Platform.LINUX_ARM64],
            tags=["docker.io/pulumi/pulumi:3.107.0-arm64"])
        index = docker_build.Index("index",
            sources=[
                amd64.ref,
                arm64.ref,
            ],
            tag="docker.io/pulumi/pulumi:3.107.0")
        pulumi.export("ref", index.ref)
        ```

        :param str resource_name: The name of the resource.
        :param pulumi.ResourceOptions opts: Options for the resource.
        :param pulumi.Input[bool] push: If true, push the index to the target registry.
               
               Defaults to `true`.
        :param pulumi.Input[Union['RegistryArgs', 'RegistryArgsDict']] registry: Authentication for the registry where the tagged index will be pushed.
               
               Credentials can also be included with the provider's configuration.
        :param pulumi.Input[Sequence[pulumi.Input[str]]] sources: Existing images to include in the index.
        :param pulumi.Input[str] tag: The tag to apply to the index.
        """
        ...
    @overload
    def __init__(__self__,
                 resource_name: str,
                 args: IndexArgs,
                 opts: Optional[pulumi.ResourceOptions] = None):
        """
        A wrapper around `docker buildx imagetools create` to create an index
        (or manifest list) referencing one or more existing images.

        In most cases you do not need an `Index` to build a multi-platform
        image -- specifying multiple platforms on the `Image` will handle this
        for you automatically.

        However, as of April 2024, building multi-platform images _with
        caching_ will only export a cache for one platform at a time (see [this
        discussion](https://github.com/docker/buildx/discussions/1382) for more
        details).

        Therefore this resource can be helpful if you are building
        multi-platform images with caching: each platform can be built and
        cached separately, and an `Index` can join them all together. An
        example of this is shown below.

        This resource creates an OCI image index or a Docker manifest list
        depending on the media types of the source images.

        ## Example Usage
        ### Multi-platform registry caching
        ```python
        import pulumi
        import pulumi_docker_build as docker_build

        amd64 = docker_build.Image("amd64",
            cache_from=[{
                "registry": {
                    "ref": "docker.io/pulumi/pulumi:cache-amd64",
                },
            }],
            cache_to=[{
                "registry": {
                    "mode": docker_build.CacheMode.MAX,
                    "ref": "docker.io/pulumi/pulumi:cache-amd64",
                },
            }],
            context={
                "location": "app",
            },
            platforms=[docker_build.Platform.LINUX_AMD64],
            tags=["docker.io/pulumi/pulumi:3.107.0-amd64"])
        arm64 = docker_build.Image("arm64",
            cache_from=[{
                "registry": {
                    "ref": "docker.io/pulumi/pulumi:cache-arm64",
                },
            }],
            cache_to=[{
                "registry": {
                    "mode": docker_build.CacheMode.MAX,
                    "ref": "docker.io/pulumi/pulumi:cache-arm64",
                },
            }],
            context={
                "location": "app",
            },
            platforms=[docker_build.Platform.LINUX_ARM64],
            tags=["docker.io/pulumi/pulumi:3.107.0-arm64"])
        index = docker_build.Index("index",
            sources=[
                amd64.ref,
                arm64.ref,
            ],
            tag="docker.io/pulumi/pulumi:3.107.0")
        pulumi.export("ref", index.ref)
        ```

        :param str resource_name: The name of the resource.
        :param IndexArgs args: The arguments to use to populate this resource's properties.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        ...
    def __init__(__self__, resource_name: str, *args, **kwargs):
        resource_args, opts = _utilities.get_resource_args_opts(IndexArgs, pulumi.ResourceOptions, *args, **kwargs)
        if resource_args is not None:
            __self__._internal_init(resource_name, opts, **resource_args.__dict__)
        else:
            __self__._internal_init(resource_name, *args, **kwargs)

    def _internal_init(__self__,
                 resource_name: str,
                 opts: Optional[pulumi.ResourceOptions] = None,
                 push: Optional[pulumi.Input[bool]] = None,
                 registry: Optional[pulumi.Input[Union['RegistryArgs', 'RegistryArgsDict']]] = None,
                 sources: Optional[pulumi.Input[Sequence[pulumi.Input[str]]]] = None,
                 tag: Optional[pulumi.Input[str]] = None,
                 __props__=None):
        opts = pulumi.ResourceOptions.merge(_utilities.get_resource_opts_defaults(), opts)
        if not isinstance(opts, pulumi.ResourceOptions):
            raise TypeError('Expected resource options to be a ResourceOptions instance')
        if opts.id is None:
            if __props__ is not None:
                raise TypeError('__props__ is only valid when passed in combination with a valid opts.id to get an existing resource')
            __props__ = IndexArgs.__new__(IndexArgs)

            if push is None:
                push = True
            __props__.__dict__["push"] = push
            __props__.__dict__["registry"] = registry
            if sources is None and not opts.urn:
                raise TypeError("Missing required property 'sources'")
            __props__.__dict__["sources"] = sources
            if tag is None and not opts.urn:
                raise TypeError("Missing required property 'tag'")
            __props__.__dict__["tag"] = tag
            __props__.__dict__["ref"] = None
        super(Index, __self__).__init__(
            'docker-build:index:Index',
            resource_name,
            __props__,
            opts)

    @staticmethod
    def get(resource_name: str,
            id: pulumi.Input[str],
            opts: Optional[pulumi.ResourceOptions] = None) -> 'Index':
        """
        Get an existing Index resource's state with the given name, id, and optional extra
        properties used to qualify the lookup.

        :param str resource_name: The unique name of the resulting resource.
        :param pulumi.Input[str] id: The unique provider ID of the resource to lookup.
        :param pulumi.ResourceOptions opts: Options for the resource.
        """
        opts = pulumi.ResourceOptions.merge(opts, pulumi.ResourceOptions(id=id))

        __props__ = IndexArgs.__new__(IndexArgs)

        __props__.__dict__["push"] = None
        __props__.__dict__["ref"] = None
        __props__.__dict__["registry"] = None
        __props__.__dict__["sources"] = None
        __props__.__dict__["tag"] = None
        return Index(resource_name, opts=opts, __props__=__props__)

    @property
    @pulumi.getter
    def push(self) -> pulumi.Output[Optional[bool]]:
        """
        If true, push the index to the target registry.

        Defaults to `true`.
        """
        return pulumi.get(self, "push")

    @property
    @pulumi.getter
    def ref(self) -> pulumi.Output[str]:
        """
        The pushed tag with digest.

        Identical to the tag if the index was not pushed.
        """
        return pulumi.get(self, "ref")

    @property
    @pulumi.getter
    def registry(self) -> pulumi.Output[Optional['outputs.Registry']]:
        """
        Authentication for the registry where the tagged index will be pushed.

        Credentials can also be included with the provider's configuration.
        """
        return pulumi.get(self, "registry")

    @property
    @pulumi.getter
    def sources(self) -> pulumi.Output[Sequence[str]]:
        """
        Existing images to include in the index.
        """
        return pulumi.get(self, "sources")

    @property
    @pulumi.getter
    def tag(self) -> pulumi.Output[str]:
        """
        The tag to apply to the index.
        """
        return pulumi.get(self, "tag")

