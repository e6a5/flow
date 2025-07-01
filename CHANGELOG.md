# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of Flow
- Focus timer with countdown display
- Task tagging with `--tag` flag
- Session logging to `~/.flowlog`
- Unix philosophy compliance - no built-in integrations
- Shell composition examples for Zenta integration
- Version information with `--version` flag
- Cross-platform build support
- CI/CD with GitHub Actions
- One-liner installation script (`install.sh`)
- Automated platform detection and binary installation
- Automatic creation of installation directory if it doesn't exist

### Changed
- Upgraded to Go 1.23.0
- Removed direct Zenta integration in favor of Unix composition

## Philosophy

Flow follows semantic versioning and Unix philosophy. Breaking changes will only be introduced in major versions, and we strive to maintain backward compatibility.

## [0.1.0] - 2025-07-01

### Added
- **Session Logging**: Automatic tracking of completed sessions with the `flow log` command.
- **Partitioned Log Files**: Logs are now stored in `YYYYMM_sessions.jsonl` files for improved performance and scalability.
- **Log Filtering**: Filter logs by `--today`, `--week`, `--month`, or a specific month (`YYYY-MM`).
- **Log Statistics**: View summary statistics with `flow log --stats`.
- **Shell Completions**: Added completions for the `log` command and its flags.

### Changed
- **Messaging**: Updated the `flow start` message to be more accurate ("Deep work session initiated").
- **Code Quality**: Fixed all `errcheck` linting issues for improved reliability.
- **Build Process**: Added a `lint` command to the `Makefile`.

### Fixed
- **Error Handling**: Improved error handling for file operations.

## [0.0.2] - 2024-06-15

### Added
- Initial release of Flow
- Session management: `start`, `pause`, `resume`, `end`
- XDG-compliant storage at `~/.local/share/flow/session`
- Automation hooks for workflow integration
- Single-session enforcement for focus 