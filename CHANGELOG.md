# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.4] - 2025-07-19

- **Paused Session Working Time**: The `flow status` command now shows how much time you've actually worked when a session is paused, excluding pause time for accurate productivity tracking.

## [1.1.3] - 2025-07-12

### Added

- **Session Targets**: Set a duration goal for your work session with `flow start --target 2h`. The `status` command will now display your progress and remaining time.
- **Daily Goals**: Set and track a daily focus goal with the new `flow goal` command. Use `flow goal --set 4h` to define your target and `flow goal` to view your progress.
- **Recent Sessions Summary**: Get a quick summary of today's completed sessions with the new `flow recent` command.
- **Productivity Insights**: Analyze your work patterns with the new `flow insights` command, which shows your busiest day, average session length, and more.
- **System Doctor**: Diagnose and troubleshoot your setup with the new `flow doctor` command to check for common configuration and data issues.

### Changed

- The `flow status` command now provides more detailed output for active sessions that have a target duration.
- The `flow pause` command output has been updated for better clarity.

### Fixed

- Corrected a test in the end-to-end suite that was failing due to updated command output, making the test suite more robust.

## [1.1.2] - 2025-07-10

### Added

- **Watcher**: A new `flow watch` command that runs as a long-running process to provide gentle, timely reminders to start, pause, resume, or end a focus session. This is an opt-in feature designed to help users who forget to interact with the timer.
- **Watcher Configuration**: The watcher's reminder timings can be customized via a new `~/.config/flow/config.yml` file. See the [Customization Guide](docs/CUSTOMIZATION.md) for details.

### Changed

- **Dependency**: Migrated from the archived `gopkg.in/yaml.v3` to the actively maintained `github.com/goccy/go-yaml` for improved security and reliability.

### Fixed

- Addressed a potential supply chain vulnerability by replacing an archived dependency.

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
