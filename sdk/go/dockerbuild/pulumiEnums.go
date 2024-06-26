// Code generated by pulumi-language-go DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package dockerbuild

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumix"
)

type CacheMode string

const (
	// Only layers that are exported into the resulting image are cached.
	CacheModeMin = CacheMode("min")
	// All layers are cached, even those of intermediate steps.
	CacheModeMax = CacheMode("max")
)

func (CacheMode) ElementType() reflect.Type {
	return reflect.TypeOf((*CacheMode)(nil)).Elem()
}

func (e CacheMode) ToCacheModeOutput() CacheModeOutput {
	return pulumi.ToOutput(e).(CacheModeOutput)
}

func (e CacheMode) ToCacheModeOutputWithContext(ctx context.Context) CacheModeOutput {
	return pulumi.ToOutputWithContext(ctx, e).(CacheModeOutput)
}

func (e CacheMode) ToCacheModePtrOutput() CacheModePtrOutput {
	return e.ToCacheModePtrOutputWithContext(context.Background())
}

func (e CacheMode) ToCacheModePtrOutputWithContext(ctx context.Context) CacheModePtrOutput {
	return CacheMode(e).ToCacheModeOutputWithContext(ctx).ToCacheModePtrOutputWithContext(ctx)
}

func (e CacheMode) ToStringOutput() pulumi.StringOutput {
	return pulumi.ToOutput(pulumi.String(e)).(pulumi.StringOutput)
}

func (e CacheMode) ToStringOutputWithContext(ctx context.Context) pulumi.StringOutput {
	return pulumi.ToOutputWithContext(ctx, pulumi.String(e)).(pulumi.StringOutput)
}

func (e CacheMode) ToStringPtrOutput() pulumi.StringPtrOutput {
	return pulumi.String(e).ToStringPtrOutputWithContext(context.Background())
}

func (e CacheMode) ToStringPtrOutputWithContext(ctx context.Context) pulumi.StringPtrOutput {
	return pulumi.String(e).ToStringOutputWithContext(ctx).ToStringPtrOutputWithContext(ctx)
}

type CacheModeOutput struct{ *pulumi.OutputState }

func (CacheModeOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*CacheMode)(nil)).Elem()
}

func (o CacheModeOutput) ToCacheModeOutput() CacheModeOutput {
	return o
}

func (o CacheModeOutput) ToCacheModeOutputWithContext(ctx context.Context) CacheModeOutput {
	return o
}

func (o CacheModeOutput) ToCacheModePtrOutput() CacheModePtrOutput {
	return o.ToCacheModePtrOutputWithContext(context.Background())
}

func (o CacheModeOutput) ToCacheModePtrOutputWithContext(ctx context.Context) CacheModePtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, v CacheMode) *CacheMode {
		return &v
	}).(CacheModePtrOutput)
}

func (o CacheModeOutput) ToOutput(ctx context.Context) pulumix.Output[CacheMode] {
	return pulumix.Output[CacheMode]{
		OutputState: o.OutputState,
	}
}

func (o CacheModeOutput) ToStringOutput() pulumi.StringOutput {
	return o.ToStringOutputWithContext(context.Background())
}

func (o CacheModeOutput) ToStringOutputWithContext(ctx context.Context) pulumi.StringOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, e CacheMode) string {
		return string(e)
	}).(pulumi.StringOutput)
}

func (o CacheModeOutput) ToStringPtrOutput() pulumi.StringPtrOutput {
	return o.ToStringPtrOutputWithContext(context.Background())
}

func (o CacheModeOutput) ToStringPtrOutputWithContext(ctx context.Context) pulumi.StringPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, e CacheMode) *string {
		v := string(e)
		return &v
	}).(pulumi.StringPtrOutput)
}

type CacheModePtrOutput struct{ *pulumi.OutputState }

func (CacheModePtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**CacheMode)(nil)).Elem()
}

func (o CacheModePtrOutput) ToCacheModePtrOutput() CacheModePtrOutput {
	return o
}

func (o CacheModePtrOutput) ToCacheModePtrOutputWithContext(ctx context.Context) CacheModePtrOutput {
	return o
}

func (o CacheModePtrOutput) ToOutput(ctx context.Context) pulumix.Output[*CacheMode] {
	return pulumix.Output[*CacheMode]{
		OutputState: o.OutputState,
	}
}

func (o CacheModePtrOutput) Elem() CacheModeOutput {
	return o.ApplyT(func(v *CacheMode) CacheMode {
		if v != nil {
			return *v
		}
		var ret CacheMode
		return ret
	}).(CacheModeOutput)
}

func (o CacheModePtrOutput) ToStringPtrOutput() pulumi.StringPtrOutput {
	return o.ToStringPtrOutputWithContext(context.Background())
}

func (o CacheModePtrOutput) ToStringPtrOutputWithContext(ctx context.Context) pulumi.StringPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, e *CacheMode) *string {
		if e == nil {
			return nil
		}
		v := string(*e)
		return &v
	}).(pulumi.StringPtrOutput)
}

// CacheModeInput is an input type that accepts values of the CacheMode enum
// A concrete instance of `CacheModeInput` can be one of the following:
//
//	CacheModeMin
//	CacheModeMax
type CacheModeInput interface {
	pulumi.Input

	ToCacheModeOutput() CacheModeOutput
	ToCacheModeOutputWithContext(context.Context) CacheModeOutput
}

var cacheModePtrType = reflect.TypeOf((**CacheMode)(nil)).Elem()

type CacheModePtrInput interface {
	pulumi.Input

	ToCacheModePtrOutput() CacheModePtrOutput
	ToCacheModePtrOutputWithContext(context.Context) CacheModePtrOutput
}

type cacheModePtr string

func CacheModePtr(v string) CacheModePtrInput {
	return (*cacheModePtr)(&v)
}

func (*cacheModePtr) ElementType() reflect.Type {
	return cacheModePtrType
}

func (in *cacheModePtr) ToCacheModePtrOutput() CacheModePtrOutput {
	return pulumi.ToOutput(in).(CacheModePtrOutput)
}

func (in *cacheModePtr) ToCacheModePtrOutputWithContext(ctx context.Context) CacheModePtrOutput {
	return pulumi.ToOutputWithContext(ctx, in).(CacheModePtrOutput)
}

func (in *cacheModePtr) ToOutput(ctx context.Context) pulumix.Output[*CacheMode] {
	return pulumix.Output[*CacheMode]{
		OutputState: in.ToCacheModePtrOutputWithContext(ctx).OutputState,
	}
}

type CompressionType string

const (
	// Use `gzip` for compression.
	CompressionTypeGzip = CompressionType("gzip")
	// Use `estargz` for compression.
	CompressionTypeEstargz = CompressionType("estargz")
	// Use `zstd` for compression.
	CompressionTypeZstd = CompressionType("zstd")
)

func (CompressionType) ElementType() reflect.Type {
	return reflect.TypeOf((*CompressionType)(nil)).Elem()
}

func (e CompressionType) ToCompressionTypeOutput() CompressionTypeOutput {
	return pulumi.ToOutput(e).(CompressionTypeOutput)
}

func (e CompressionType) ToCompressionTypeOutputWithContext(ctx context.Context) CompressionTypeOutput {
	return pulumi.ToOutputWithContext(ctx, e).(CompressionTypeOutput)
}

func (e CompressionType) ToCompressionTypePtrOutput() CompressionTypePtrOutput {
	return e.ToCompressionTypePtrOutputWithContext(context.Background())
}

func (e CompressionType) ToCompressionTypePtrOutputWithContext(ctx context.Context) CompressionTypePtrOutput {
	return CompressionType(e).ToCompressionTypeOutputWithContext(ctx).ToCompressionTypePtrOutputWithContext(ctx)
}

func (e CompressionType) ToStringOutput() pulumi.StringOutput {
	return pulumi.ToOutput(pulumi.String(e)).(pulumi.StringOutput)
}

func (e CompressionType) ToStringOutputWithContext(ctx context.Context) pulumi.StringOutput {
	return pulumi.ToOutputWithContext(ctx, pulumi.String(e)).(pulumi.StringOutput)
}

func (e CompressionType) ToStringPtrOutput() pulumi.StringPtrOutput {
	return pulumi.String(e).ToStringPtrOutputWithContext(context.Background())
}

func (e CompressionType) ToStringPtrOutputWithContext(ctx context.Context) pulumi.StringPtrOutput {
	return pulumi.String(e).ToStringOutputWithContext(ctx).ToStringPtrOutputWithContext(ctx)
}

type CompressionTypeOutput struct{ *pulumi.OutputState }

func (CompressionTypeOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*CompressionType)(nil)).Elem()
}

func (o CompressionTypeOutput) ToCompressionTypeOutput() CompressionTypeOutput {
	return o
}

func (o CompressionTypeOutput) ToCompressionTypeOutputWithContext(ctx context.Context) CompressionTypeOutput {
	return o
}

func (o CompressionTypeOutput) ToCompressionTypePtrOutput() CompressionTypePtrOutput {
	return o.ToCompressionTypePtrOutputWithContext(context.Background())
}

func (o CompressionTypeOutput) ToCompressionTypePtrOutputWithContext(ctx context.Context) CompressionTypePtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, v CompressionType) *CompressionType {
		return &v
	}).(CompressionTypePtrOutput)
}

func (o CompressionTypeOutput) ToOutput(ctx context.Context) pulumix.Output[CompressionType] {
	return pulumix.Output[CompressionType]{
		OutputState: o.OutputState,
	}
}

func (o CompressionTypeOutput) ToStringOutput() pulumi.StringOutput {
	return o.ToStringOutputWithContext(context.Background())
}

func (o CompressionTypeOutput) ToStringOutputWithContext(ctx context.Context) pulumi.StringOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, e CompressionType) string {
		return string(e)
	}).(pulumi.StringOutput)
}

func (o CompressionTypeOutput) ToStringPtrOutput() pulumi.StringPtrOutput {
	return o.ToStringPtrOutputWithContext(context.Background())
}

func (o CompressionTypeOutput) ToStringPtrOutputWithContext(ctx context.Context) pulumi.StringPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, e CompressionType) *string {
		v := string(e)
		return &v
	}).(pulumi.StringPtrOutput)
}

type CompressionTypePtrOutput struct{ *pulumi.OutputState }

func (CompressionTypePtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**CompressionType)(nil)).Elem()
}

func (o CompressionTypePtrOutput) ToCompressionTypePtrOutput() CompressionTypePtrOutput {
	return o
}

func (o CompressionTypePtrOutput) ToCompressionTypePtrOutputWithContext(ctx context.Context) CompressionTypePtrOutput {
	return o
}

func (o CompressionTypePtrOutput) ToOutput(ctx context.Context) pulumix.Output[*CompressionType] {
	return pulumix.Output[*CompressionType]{
		OutputState: o.OutputState,
	}
}

func (o CompressionTypePtrOutput) Elem() CompressionTypeOutput {
	return o.ApplyT(func(v *CompressionType) CompressionType {
		if v != nil {
			return *v
		}
		var ret CompressionType
		return ret
	}).(CompressionTypeOutput)
}

func (o CompressionTypePtrOutput) ToStringPtrOutput() pulumi.StringPtrOutput {
	return o.ToStringPtrOutputWithContext(context.Background())
}

func (o CompressionTypePtrOutput) ToStringPtrOutputWithContext(ctx context.Context) pulumi.StringPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, e *CompressionType) *string {
		if e == nil {
			return nil
		}
		v := string(*e)
		return &v
	}).(pulumi.StringPtrOutput)
}

// CompressionTypeInput is an input type that accepts values of the CompressionType enum
// A concrete instance of `CompressionTypeInput` can be one of the following:
//
//	CompressionTypeGzip
//	CompressionTypeEstargz
//	CompressionTypeZstd
type CompressionTypeInput interface {
	pulumi.Input

	ToCompressionTypeOutput() CompressionTypeOutput
	ToCompressionTypeOutputWithContext(context.Context) CompressionTypeOutput
}

var compressionTypePtrType = reflect.TypeOf((**CompressionType)(nil)).Elem()

type CompressionTypePtrInput interface {
	pulumi.Input

	ToCompressionTypePtrOutput() CompressionTypePtrOutput
	ToCompressionTypePtrOutputWithContext(context.Context) CompressionTypePtrOutput
}

type compressionTypePtr string

func CompressionTypePtr(v string) CompressionTypePtrInput {
	return (*compressionTypePtr)(&v)
}

func (*compressionTypePtr) ElementType() reflect.Type {
	return compressionTypePtrType
}

func (in *compressionTypePtr) ToCompressionTypePtrOutput() CompressionTypePtrOutput {
	return pulumi.ToOutput(in).(CompressionTypePtrOutput)
}

func (in *compressionTypePtr) ToCompressionTypePtrOutputWithContext(ctx context.Context) CompressionTypePtrOutput {
	return pulumi.ToOutputWithContext(ctx, in).(CompressionTypePtrOutput)
}

func (in *compressionTypePtr) ToOutput(ctx context.Context) pulumix.Output[*CompressionType] {
	return pulumix.Output[*CompressionType]{
		OutputState: in.ToCompressionTypePtrOutputWithContext(ctx).OutputState,
	}
}

type NetworkMode string

const (
	// The default sandbox network mode.
	NetworkModeDefault = NetworkMode("default")
	// Host network mode.
	NetworkModeHost = NetworkMode("host")
	// Disable network access.
	NetworkModeNone = NetworkMode("none")
)

func (NetworkMode) ElementType() reflect.Type {
	return reflect.TypeOf((*NetworkMode)(nil)).Elem()
}

func (e NetworkMode) ToNetworkModeOutput() NetworkModeOutput {
	return pulumi.ToOutput(e).(NetworkModeOutput)
}

func (e NetworkMode) ToNetworkModeOutputWithContext(ctx context.Context) NetworkModeOutput {
	return pulumi.ToOutputWithContext(ctx, e).(NetworkModeOutput)
}

func (e NetworkMode) ToNetworkModePtrOutput() NetworkModePtrOutput {
	return e.ToNetworkModePtrOutputWithContext(context.Background())
}

func (e NetworkMode) ToNetworkModePtrOutputWithContext(ctx context.Context) NetworkModePtrOutput {
	return NetworkMode(e).ToNetworkModeOutputWithContext(ctx).ToNetworkModePtrOutputWithContext(ctx)
}

func (e NetworkMode) ToStringOutput() pulumi.StringOutput {
	return pulumi.ToOutput(pulumi.String(e)).(pulumi.StringOutput)
}

func (e NetworkMode) ToStringOutputWithContext(ctx context.Context) pulumi.StringOutput {
	return pulumi.ToOutputWithContext(ctx, pulumi.String(e)).(pulumi.StringOutput)
}

func (e NetworkMode) ToStringPtrOutput() pulumi.StringPtrOutput {
	return pulumi.String(e).ToStringPtrOutputWithContext(context.Background())
}

func (e NetworkMode) ToStringPtrOutputWithContext(ctx context.Context) pulumi.StringPtrOutput {
	return pulumi.String(e).ToStringOutputWithContext(ctx).ToStringPtrOutputWithContext(ctx)
}

type NetworkModeOutput struct{ *pulumi.OutputState }

func (NetworkModeOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*NetworkMode)(nil)).Elem()
}

func (o NetworkModeOutput) ToNetworkModeOutput() NetworkModeOutput {
	return o
}

func (o NetworkModeOutput) ToNetworkModeOutputWithContext(ctx context.Context) NetworkModeOutput {
	return o
}

func (o NetworkModeOutput) ToNetworkModePtrOutput() NetworkModePtrOutput {
	return o.ToNetworkModePtrOutputWithContext(context.Background())
}

func (o NetworkModeOutput) ToNetworkModePtrOutputWithContext(ctx context.Context) NetworkModePtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, v NetworkMode) *NetworkMode {
		return &v
	}).(NetworkModePtrOutput)
}

func (o NetworkModeOutput) ToOutput(ctx context.Context) pulumix.Output[NetworkMode] {
	return pulumix.Output[NetworkMode]{
		OutputState: o.OutputState,
	}
}

func (o NetworkModeOutput) ToStringOutput() pulumi.StringOutput {
	return o.ToStringOutputWithContext(context.Background())
}

func (o NetworkModeOutput) ToStringOutputWithContext(ctx context.Context) pulumi.StringOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, e NetworkMode) string {
		return string(e)
	}).(pulumi.StringOutput)
}

func (o NetworkModeOutput) ToStringPtrOutput() pulumi.StringPtrOutput {
	return o.ToStringPtrOutputWithContext(context.Background())
}

func (o NetworkModeOutput) ToStringPtrOutputWithContext(ctx context.Context) pulumi.StringPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, e NetworkMode) *string {
		v := string(e)
		return &v
	}).(pulumi.StringPtrOutput)
}

type NetworkModePtrOutput struct{ *pulumi.OutputState }

func (NetworkModePtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**NetworkMode)(nil)).Elem()
}

func (o NetworkModePtrOutput) ToNetworkModePtrOutput() NetworkModePtrOutput {
	return o
}

func (o NetworkModePtrOutput) ToNetworkModePtrOutputWithContext(ctx context.Context) NetworkModePtrOutput {
	return o
}

func (o NetworkModePtrOutput) ToOutput(ctx context.Context) pulumix.Output[*NetworkMode] {
	return pulumix.Output[*NetworkMode]{
		OutputState: o.OutputState,
	}
}

func (o NetworkModePtrOutput) Elem() NetworkModeOutput {
	return o.ApplyT(func(v *NetworkMode) NetworkMode {
		if v != nil {
			return *v
		}
		var ret NetworkMode
		return ret
	}).(NetworkModeOutput)
}

func (o NetworkModePtrOutput) ToStringPtrOutput() pulumi.StringPtrOutput {
	return o.ToStringPtrOutputWithContext(context.Background())
}

func (o NetworkModePtrOutput) ToStringPtrOutputWithContext(ctx context.Context) pulumi.StringPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, e *NetworkMode) *string {
		if e == nil {
			return nil
		}
		v := string(*e)
		return &v
	}).(pulumi.StringPtrOutput)
}

// NetworkModeInput is an input type that accepts values of the NetworkMode enum
// A concrete instance of `NetworkModeInput` can be one of the following:
//
//	NetworkModeDefault
//	NetworkModeHost
//	NetworkModeNone
type NetworkModeInput interface {
	pulumi.Input

	ToNetworkModeOutput() NetworkModeOutput
	ToNetworkModeOutputWithContext(context.Context) NetworkModeOutput
}

var networkModePtrType = reflect.TypeOf((**NetworkMode)(nil)).Elem()

type NetworkModePtrInput interface {
	pulumi.Input

	ToNetworkModePtrOutput() NetworkModePtrOutput
	ToNetworkModePtrOutputWithContext(context.Context) NetworkModePtrOutput
}

type networkModePtr string

func NetworkModePtr(v string) NetworkModePtrInput {
	return (*networkModePtr)(&v)
}

func (*networkModePtr) ElementType() reflect.Type {
	return networkModePtrType
}

func (in *networkModePtr) ToNetworkModePtrOutput() NetworkModePtrOutput {
	return pulumi.ToOutput(in).(NetworkModePtrOutput)
}

func (in *networkModePtr) ToNetworkModePtrOutputWithContext(ctx context.Context) NetworkModePtrOutput {
	return pulumi.ToOutputWithContext(ctx, in).(NetworkModePtrOutput)
}

func (in *networkModePtr) ToOutput(ctx context.Context) pulumix.Output[*NetworkMode] {
	return pulumix.Output[*NetworkMode]{
		OutputState: in.ToNetworkModePtrOutputWithContext(ctx).OutputState,
	}
}

type Platform string

const (
	Platform_Darwin_386      = Platform("darwin/386")
	Platform_Darwin_amd64    = Platform("darwin/amd64")
	Platform_Darwin_arm      = Platform("darwin/arm")
	Platform_Darwin_arm64    = Platform("darwin/arm64")
	Platform_Dragonfly_amd64 = Platform("dragonfly/amd64")
	Platform_Freebsd_386     = Platform("freebsd/386")
	Platform_Freebsd_amd64   = Platform("freebsd/amd64")
	Platform_Freebsd_arm     = Platform("freebsd/arm")
	Platform_Linux_386       = Platform("linux/386")
	Platform_Linux_amd64     = Platform("linux/amd64")
	Platform_Linux_arm       = Platform("linux/arm")
	Platform_Linux_arm64     = Platform("linux/arm64")
	Platform_Linux_mips64    = Platform("linux/mips64")
	Platform_Linux_mips64le  = Platform("linux/mips64le")
	Platform_Linux_ppc64le   = Platform("linux/ppc64le")
	Platform_Linux_riscv64   = Platform("linux/riscv64")
	Platform_Linux_s390x     = Platform("linux/s390x")
	Platform_Netbsd_386      = Platform("netbsd/386")
	Platform_Netbsd_amd64    = Platform("netbsd/amd64")
	Platform_Netbsd_arm      = Platform("netbsd/arm")
	Platform_Openbsd_386     = Platform("openbsd/386")
	Platform_Openbsd_amd64   = Platform("openbsd/amd64")
	Platform_Openbsd_arm     = Platform("openbsd/arm")
	Platform_Plan9_386       = Platform("plan9/386")
	Platform_Plan9_amd64     = Platform("plan9/amd64")
	Platform_Solaris_amd64   = Platform("solaris/amd64")
	Platform_Windows_386     = Platform("windows/386")
	Platform_Windows_amd64   = Platform("windows/amd64")
)

func (Platform) ElementType() reflect.Type {
	return reflect.TypeOf((*Platform)(nil)).Elem()
}

func (e Platform) ToPlatformOutput() PlatformOutput {
	return pulumi.ToOutput(e).(PlatformOutput)
}

func (e Platform) ToPlatformOutputWithContext(ctx context.Context) PlatformOutput {
	return pulumi.ToOutputWithContext(ctx, e).(PlatformOutput)
}

func (e Platform) ToPlatformPtrOutput() PlatformPtrOutput {
	return e.ToPlatformPtrOutputWithContext(context.Background())
}

func (e Platform) ToPlatformPtrOutputWithContext(ctx context.Context) PlatformPtrOutput {
	return Platform(e).ToPlatformOutputWithContext(ctx).ToPlatformPtrOutputWithContext(ctx)
}

func (e Platform) ToStringOutput() pulumi.StringOutput {
	return pulumi.ToOutput(pulumi.String(e)).(pulumi.StringOutput)
}

func (e Platform) ToStringOutputWithContext(ctx context.Context) pulumi.StringOutput {
	return pulumi.ToOutputWithContext(ctx, pulumi.String(e)).(pulumi.StringOutput)
}

func (e Platform) ToStringPtrOutput() pulumi.StringPtrOutput {
	return pulumi.String(e).ToStringPtrOutputWithContext(context.Background())
}

func (e Platform) ToStringPtrOutputWithContext(ctx context.Context) pulumi.StringPtrOutput {
	return pulumi.String(e).ToStringOutputWithContext(ctx).ToStringPtrOutputWithContext(ctx)
}

type PlatformOutput struct{ *pulumi.OutputState }

func (PlatformOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*Platform)(nil)).Elem()
}

func (o PlatformOutput) ToPlatformOutput() PlatformOutput {
	return o
}

func (o PlatformOutput) ToPlatformOutputWithContext(ctx context.Context) PlatformOutput {
	return o
}

func (o PlatformOutput) ToPlatformPtrOutput() PlatformPtrOutput {
	return o.ToPlatformPtrOutputWithContext(context.Background())
}

func (o PlatformOutput) ToPlatformPtrOutputWithContext(ctx context.Context) PlatformPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, v Platform) *Platform {
		return &v
	}).(PlatformPtrOutput)
}

func (o PlatformOutput) ToOutput(ctx context.Context) pulumix.Output[Platform] {
	return pulumix.Output[Platform]{
		OutputState: o.OutputState,
	}
}

func (o PlatformOutput) ToStringOutput() pulumi.StringOutput {
	return o.ToStringOutputWithContext(context.Background())
}

func (o PlatformOutput) ToStringOutputWithContext(ctx context.Context) pulumi.StringOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, e Platform) string {
		return string(e)
	}).(pulumi.StringOutput)
}

func (o PlatformOutput) ToStringPtrOutput() pulumi.StringPtrOutput {
	return o.ToStringPtrOutputWithContext(context.Background())
}

func (o PlatformOutput) ToStringPtrOutputWithContext(ctx context.Context) pulumi.StringPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, e Platform) *string {
		v := string(e)
		return &v
	}).(pulumi.StringPtrOutput)
}

type PlatformPtrOutput struct{ *pulumi.OutputState }

func (PlatformPtrOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**Platform)(nil)).Elem()
}

func (o PlatformPtrOutput) ToPlatformPtrOutput() PlatformPtrOutput {
	return o
}

func (o PlatformPtrOutput) ToPlatformPtrOutputWithContext(ctx context.Context) PlatformPtrOutput {
	return o
}

func (o PlatformPtrOutput) ToOutput(ctx context.Context) pulumix.Output[*Platform] {
	return pulumix.Output[*Platform]{
		OutputState: o.OutputState,
	}
}

func (o PlatformPtrOutput) Elem() PlatformOutput {
	return o.ApplyT(func(v *Platform) Platform {
		if v != nil {
			return *v
		}
		var ret Platform
		return ret
	}).(PlatformOutput)
}

func (o PlatformPtrOutput) ToStringPtrOutput() pulumi.StringPtrOutput {
	return o.ToStringPtrOutputWithContext(context.Background())
}

func (o PlatformPtrOutput) ToStringPtrOutputWithContext(ctx context.Context) pulumi.StringPtrOutput {
	return o.ApplyTWithContext(ctx, func(_ context.Context, e *Platform) *string {
		if e == nil {
			return nil
		}
		v := string(*e)
		return &v
	}).(pulumi.StringPtrOutput)
}

// PlatformInput is an input type that accepts values of the Platform enum
// A concrete instance of `PlatformInput` can be one of the following:
//
//	Platform_Darwin_386
//	Platform_Darwin_amd64
//	Platform_Darwin_arm
//	Platform_Darwin_arm64
//	Platform_Dragonfly_amd64
//	Platform_Freebsd_386
//	Platform_Freebsd_amd64
//	Platform_Freebsd_arm
//	Platform_Linux_386
//	Platform_Linux_amd64
//	Platform_Linux_arm
//	Platform_Linux_arm64
//	Platform_Linux_mips64
//	Platform_Linux_mips64le
//	Platform_Linux_ppc64le
//	Platform_Linux_riscv64
//	Platform_Linux_s390x
//	Platform_Netbsd_386
//	Platform_Netbsd_amd64
//	Platform_Netbsd_arm
//	Platform_Openbsd_386
//	Platform_Openbsd_amd64
//	Platform_Openbsd_arm
//	Platform_Plan9_386
//	Platform_Plan9_amd64
//	Platform_Solaris_amd64
//	Platform_Windows_386
//	Platform_Windows_amd64
type PlatformInput interface {
	pulumi.Input

	ToPlatformOutput() PlatformOutput
	ToPlatformOutputWithContext(context.Context) PlatformOutput
}

var platformPtrType = reflect.TypeOf((**Platform)(nil)).Elem()

type PlatformPtrInput interface {
	pulumi.Input

	ToPlatformPtrOutput() PlatformPtrOutput
	ToPlatformPtrOutputWithContext(context.Context) PlatformPtrOutput
}

type platformPtr string

func PlatformPtr(v string) PlatformPtrInput {
	return (*platformPtr)(&v)
}

func (*platformPtr) ElementType() reflect.Type {
	return platformPtrType
}

func (in *platformPtr) ToPlatformPtrOutput() PlatformPtrOutput {
	return pulumi.ToOutput(in).(PlatformPtrOutput)
}

func (in *platformPtr) ToPlatformPtrOutputWithContext(ctx context.Context) PlatformPtrOutput {
	return pulumi.ToOutputWithContext(ctx, in).(PlatformPtrOutput)
}

func (in *platformPtr) ToOutput(ctx context.Context) pulumix.Output[*Platform] {
	return pulumix.Output[*Platform]{
		OutputState: in.ToPlatformPtrOutputWithContext(ctx).OutputState,
	}
}

// PlatformArrayInput is an input type that accepts PlatformArray and PlatformArrayOutput values.
// You can construct a concrete instance of `PlatformArrayInput` via:
//
//	PlatformArray{ PlatformArgs{...} }
type PlatformArrayInput interface {
	pulumi.Input

	ToPlatformArrayOutput() PlatformArrayOutput
	ToPlatformArrayOutputWithContext(context.Context) PlatformArrayOutput
}

type PlatformArray []Platform

func (PlatformArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]Platform)(nil)).Elem()
}

func (i PlatformArray) ToPlatformArrayOutput() PlatformArrayOutput {
	return i.ToPlatformArrayOutputWithContext(context.Background())
}

func (i PlatformArray) ToPlatformArrayOutputWithContext(ctx context.Context) PlatformArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(PlatformArrayOutput)
}

func (i PlatformArray) ToOutput(ctx context.Context) pulumix.Output[[]Platform] {
	return pulumix.Output[[]Platform]{
		OutputState: i.ToPlatformArrayOutputWithContext(ctx).OutputState,
	}
}

type PlatformArrayOutput struct{ *pulumi.OutputState }

func (PlatformArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]Platform)(nil)).Elem()
}

func (o PlatformArrayOutput) ToPlatformArrayOutput() PlatformArrayOutput {
	return o
}

func (o PlatformArrayOutput) ToPlatformArrayOutputWithContext(ctx context.Context) PlatformArrayOutput {
	return o
}

func (o PlatformArrayOutput) ToOutput(ctx context.Context) pulumix.Output[[]Platform] {
	return pulumix.Output[[]Platform]{
		OutputState: o.OutputState,
	}
}

func (o PlatformArrayOutput) Index(i pulumi.IntInput) PlatformOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) Platform {
		return vs[0].([]Platform)[vs[1].(int)]
	}).(PlatformOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*CacheModeInput)(nil)).Elem(), CacheMode("min"))
	pulumi.RegisterInputType(reflect.TypeOf((*CacheModePtrInput)(nil)).Elem(), CacheMode("min"))
	pulumi.RegisterInputType(reflect.TypeOf((*CompressionTypeInput)(nil)).Elem(), CompressionType("gzip"))
	pulumi.RegisterInputType(reflect.TypeOf((*CompressionTypePtrInput)(nil)).Elem(), CompressionType("gzip"))
	pulumi.RegisterInputType(reflect.TypeOf((*NetworkModeInput)(nil)).Elem(), NetworkMode("default"))
	pulumi.RegisterInputType(reflect.TypeOf((*NetworkModePtrInput)(nil)).Elem(), NetworkMode("default"))
	pulumi.RegisterInputType(reflect.TypeOf((*PlatformInput)(nil)).Elem(), Platform("darwin/386"))
	pulumi.RegisterInputType(reflect.TypeOf((*PlatformPtrInput)(nil)).Elem(), Platform("darwin/386"))
	pulumi.RegisterInputType(reflect.TypeOf((*PlatformArrayInput)(nil)).Elem(), PlatformArray{})
	pulumi.RegisterOutputType(CacheModeOutput{})
	pulumi.RegisterOutputType(CacheModePtrOutput{})
	pulumi.RegisterOutputType(CompressionTypeOutput{})
	pulumi.RegisterOutputType(CompressionTypePtrOutput{})
	pulumi.RegisterOutputType(NetworkModeOutput{})
	pulumi.RegisterOutputType(NetworkModePtrOutput{})
	pulumi.RegisterOutputType(PlatformOutput{})
	pulumi.RegisterOutputType(PlatformPtrOutput{})
	pulumi.RegisterOutputType(PlatformArrayOutput{})
}
