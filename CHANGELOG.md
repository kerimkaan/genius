# Changelog

All notable changes to this project will be documented in this file.

## v0.0.10 - 2026-02-11

### Changed
- **Breaking**: Replaced `Run` with `RunE` on all commands for proper error propagation
- **Breaking**: Removed `miekg/dns` dependency; DNS servers are now parsed directly from `/etc/resolv.conf`
- Network interfaces: show all active non-loopback interfaces instead of hardcoded `en0`/`ens160`
- NTP parser: rewritten line-by-line; supports `/etc/ntp.conf`, `/etc/chrony/chrony.conf`, `/etc/chrony.conf`
- NTP server: now configurable via `--ntp-server` flag (default: `pool.ntp.org`)
- Windows compatibility check moved to root command `PersistentPreRunE` (applies to all subcommands)
- Output formatting: consistent aligned labels with `fmt.Printf`
- Go naming: `VERSION` constant replaced with `Version` variable (supports ldflags override)
- `FolderSize` type moved to `types/` package
- `helpers/packages.go`: external commands now use `context.WithTimeout` (5s)
- `helpers/file.go`: goroutine pool limited to 10 concurrent walkers

### Added
- CI pipeline (`.github/workflows/ci.yml`): build, vet, test, golangci-lint
- `--count`/`-n` flag for `largest-folders` command
- `lf` alias for `largest-folders` command
- `IsLinux()` helper function
- Comprehensive test suite: 16 tests across `cmd/` and `helpers/` packages

### Fixed
- CI: `build-release.yml` darwin/amd64 build had wrong `build_tags: darwin-arm64`
- Unsafe `addr.(*net.IPNet)` type assertion that could panic
- NTP config parser incorrectly removing everything after first `#` character
- Error messages now follow Go conventions (lowercase, no trailing punctuation)

### Removed
- `miekg/dns` dependency and its transitive dependencies
- Unused `GetLargestFolders` (non-concurrent version)
- Unused `toggle` flag on root command
- Redundant `MAJOR`, `MINOR`, `PATCH` constants
- Scaffold boilerplate comments from cobra

## v0.0.9 - 2025-07-01

### Changed
- Updated dns to v1.1.66
- Updated go to 1.24.4

## v0.0.8 - 2025-05-01

### Added
- `largest-folders` command: Lists the top 5 largest folders in the user's home directory, using concurrency for speed.
- Async folder size calculation for better performance and fewer permission errors.
- Helper function to get the user's home directory.
