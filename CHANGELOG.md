# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2024-01-29

### Fixed
- Added automatic fixes for common mistakes in changelog.md file on parsing
- Added required braces on version number (in output only!)
- Fixed style of init version of the CHANGELOG.md file (added new line after the title)
- Fixed `-file` parameter (it wasn't used)
- Fixed tag alias in the docker registry (v1.0.0 instead of v1-0-0)
- Fixed build and publish stages in CI pipeline
- Fixed typos in the README.md
- Fixed tests after using refactoring tool
- Fixed using env variables in github actions

### Added
- Created a simple skeleton for cli command
- Added markdown changelog format parser
- Supported `-command=diff` command
- Supported `-command=bump` command
- Supported `-command=init` command
- Add command `-command=latest_version` for getting latest described version (exclude unreleased, returns just version number)
- Allow to use `latest` keyword in `direction` and `diff` commands
- Add bool param `fail-on-empty`: if it's true and no changes is found then exit code will be not 0 (not ok)
- Add cli command for checking deployment direction (UPGRADE, ROLLBACK, REDEPLOY)
- Added gitlab CI pipeline
- Add release notes notification to the Slack
- Add tests for changes in the CHANGELOG.md on MR
- Support read changelog from STDIN

### Changed
- If `from` and `to` are equal in `diff` command then changes in this version will be output
- Use shared tasks in CI pipeline (partially)
- If no changes in the diff nothing will be output
