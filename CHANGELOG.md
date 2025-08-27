## Unreleased

## 0.0.13 (2025-08-27)

### Changed

- Docker Build Cloud and `exec` errors are more helpful. (https://github.com/pulumi/pulumi-docker-build/issues/549)

### Fixed

- The provider is no longer replaced on version changes. (https://github.com/pulumi/pulumi-docker-build/issues/581)

## 0.0.12 (2025-05-16)

### Changed

- Upgraded pulumi-go-provider to v1.0.0-rc2.

### Fixed

- Builds now respect cancellation. (https://github.com/pulumi/pulumi-docker-build/issues/533, https://github.com/pulumi/pulumi-docker-build/pull/522)

## 0.0.11 (2025-04-11)

### Changed

- Upgraded buildx from 0.18.0 to 0.20.1 to remain compatible with upcoming
  changes to GitHub Actions. (https://github.com/pulumi/pulumi-docker-build/pull/519)

### Fixed

- Upgrading docker-build no longer causes resource replacements. (<https://github.com/pulumi/pulumi-docker-build/issues/404>)
- Fixed a panic that could occur in `exec` mode. (https://github.com/pulumi/pulumi-docker-build/issues/482)
- The default GitHub Actions cache scope is now correctly set as `buildkit`. (https://github.com/pulumi/pulumi-docker-build/issues/496)

## 0.0.10 (2025-01-27)

### Changed

- Windows binaries are now signed. (https://github.com/pulumi/pulumi-docker-build/pull/429)

## 0.0.9 (2025-01-16)

### Changed

- Upgraded pulumi-go-provider to v0.24.1. (https://github.com/pulumi/pulumi-docker-build/pull/413)

### Fixed

- `ACTIONS_RUNTIME_TOKEN` is now correctly marked as a secret. (https://github.com/pulumi/pulumi-docker-build/issues/403)

## 0.0.8 (2024-12-10)

### Added

- Multiple exports are now allowed if the build daemon is detected to have
  version 0.13 of Buildkit or newer.
  (https://github.com/pulumi/pulumi-docker-build/issues/21)

### Changed

- Upgraded buildx from 0.16.0 to 0.18.0.

### Fixed

- Custom `# syntax=` directives no longer cause validation errors.
  (https://github.com/pulumi/pulumi-docker-build/issues/300)

## 0.0.7 (2024-10-16)

### Fixed

- Fixed an issue where registry authentication couldn't be specified on the
  provider. (<https://github.com/pulumi/pulumi-docker-build/issues/262>)

## 0.0.6 (2024-08-13)

### Fixed

- Refreshing an `Index` resource will no longer fail if its stored credentials
  have expired. (<https://github.com/pulumi/pulumi-docker-build/pull/194>)

### Changed

- Local and tar exporters will now trigger an update if an export doesn't exist
  at the expected path. (<https://github.com/pulumi/pulumi-docker-build/pull/195>)

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

- Fixed the default value for `ACTIONS_CACHE_URL` when using GitHub action caching. (<https://github.com/pulumi/pulumi-docker-build/pull/80>)
- Fixed Java SDK publishing. (<https://github.com/pulumi/pulumi-docker-build/pull/89>)
- Fixed a panic that could occur when `context` was omitted. (<https://github.com/pulumi/pulumi-docker-build/pull/83>)

### Changed

- The provider will now wait for new builders to fully boot.

## 0.0.2 (2024-04-25)

### Fixed

- Upgraded pulumi-go-provider to fix a panic during cancellation.

## 0.0.1 (2024-04-23)

Initial release.
