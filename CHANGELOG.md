# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).



## [Unreleased]

### Changed

- Upgrade to Go 1.17

## [0.5.0] - 2021-01-04

### Added

- Add Logger.WithIncreasedCallerDepth to support wrapping in other interfaces.

### Fixed

- Fix caller for Logger.Debugf and Logger.Errorf.
- Fix caller in ActivationLogger.



## [0.4.0] - 2020-12-01

### Added

- Add Logger.Debugf and Logger.Errorf.



## [0.3.4] - 2020-11-05

### Fixed

- Fix `isVerbosityAllowed` default case (log verbosity undefined).

## [0.3.3] - 2020-09-15

### Fixed

- Fix indirect dependency vulnerability detected by CI.



## [0.3.2] - 2020-09-15

### Fixed

- Fix order of log level activation.



## [0.3.1] 2020-03-20

### Fixed

- Fix LogCtx panic when no LoggerMeta is given.



## [0.3.0] 2020-03-17

### Changed

- Remove error from the spec.



## [0.2.0] 2020-03-03

### Changed

- Switch to Go modules.



## [0.1.0] 2020-02-13

### Added

- First release.



[Unreleased]: https://github.com/giantswarm/micrologger/compare/v0.5.0...HEAD
[0.5.0]: https://github.com/giantswarm/micrologger/compare/v0.4.0...v0.5.0
[0.4.0]: https://github.com/giantswarm/micrologger/compare/v0.3.4...v0.4.0
[0.3.4]: https://github.com/giantswarm/micrologger/compare/v0.3.3...v0.3.4
[0.3.3]: https://github.com/giantswarm/micrologger/compare/v0.3.2...v0.3.3
[0.3.2]: https://github.com/giantswarm/micrologger/compare/v0.3.1...v0.3.2
[0.3.1]: https://github.com/giantswarm/micrologger/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/giantswarm/micrologger/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/giantswarm/micrologger/compare/v0.1.0...v0.2.0

[0.1.0]: https://github.com/giantswarm/micrologger/releases/tag/v0.1.0
