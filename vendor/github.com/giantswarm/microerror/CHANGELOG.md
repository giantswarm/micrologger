# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Add JSON function.

### Changed

- Use built-in errors package instead of juju/errgo.
- Print error stacks in JSON format instead of custom errgo format.

### Removed

- Drop Stack function in favour of JSON function.
- Drop Newf function.
- Drop Error.GoString method.
- Drop Error.String method.

## [1.0.0] - 2019-09-20

[Unreleased]: https://github.com/giantswarm/architect-orb/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/giantswarm/architect-orb/releases/tag/v1.0.0
