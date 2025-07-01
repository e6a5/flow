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