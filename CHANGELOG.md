## Unreleased

### Fixed

- Refreshing an `Index` resource will no longer fail if its stored credentials
  have expired. (https://github.com/pulumi/pulumi-docker-build/pull/194)

## 0.0.5 (2024-08-08)

### Fixed

- Fixed Go SDK publishing.

### Changed

- Upgraded docker from 27.0.3 to 27.1.0.

## 0.0.4 (2024-07-15)

### Changed

- Upgraded buildkit from 0.13.0 to 0.15.0.
- Upgraded buildx from 0.13.1. to 0.16.0.
- Upgraded docker from 26.0.0-rc1 to 27.0.3.
- Fixed an issue where warnings were not displayed correctly.

## 0.0.3 (2024-05-31)

### Fixed

- Fixed the default value for `ACTIONS_CACHE_URL` when using GitHub action caching. (https://github.com/pulumi/pulumi-docker-build/pull/80)
- Fixed Java SDK publishing. (https://github.com/pulumi/pulumi-docker-build/pull/89)
- Fixed a panic that could occur when `context` was omitted. (https://github.com/pulumi/pulumi-docker-build/pull/83)

### Changed

- The provider will now wait for new builders to fully boot.

## 0.0.2 (2024-04-25)

### Fixed

- Upgraded pulumi-go-provider to fix a panic during cancellation.

## 0.0.1 (2024-04-23)

Initial release.
