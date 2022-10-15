# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.3.1](https://github.com/Lord-Y/cypress-parallel/releases/tag/v0.3.1) - 2022-10-15

### Changed
- global:
  - Update default docker image version to 10.10.0-0.3.0
- ci:
  - Update ci requirements
- api:
  - Upgrade to golang version 1.19
  - Fix project update
- ui:
  - Update nodejs version to 16
  - Update packages

## [v0.2.0](https://github.com/Lord-Y/cypress-parallel/releases/tag/v0.2.0) - 2021-12-30

v0.2.0 release

### Added
- global:
  - Add git hooks

### Changed
- global:
  - Update default docker image version to 7.4.0-0.1.1
  - Avoid duplicate team or project creation
- api:
  - Enforce not null data
  - Refactoring postgresql queries
  - Update assets requirements
  - Update execution status
  - Improve search query
- ui:
  - Update packages
  - Update alert message when team or project exist
  - Fix team name in project

## [v0.2.0-beta2](https://github.com/Lord-Y/cypress-parallel/releases/tag/v0.2.0-beta2) - 2021-12-30

v0.2.0 release

### Added
- global:
  - Add git hooks

### Changed
- global:
  - Update default docker image version to 7.4.0-0.1.1
  - Avoid duplicate team or project creation
- api:
  - Enforce not null data
  - Refactoring postgresql queries
  - Update assets requirements
  - Update execution status
  - Improve search query
- ui:
  - Update packages
  - Update alert message when team or project exist
  - Fix team name in project

## [v0.2.0-beta1](https://github.com/Lord-Y/cypress-parallel/releases/tag/v0.2.0-beta1) - 2021-12-30

v0.2.0 release

### Added
- global:
  - Add git hooks

### Changed
- global:
  - Update default docker image version to 7.4.0-0.1.1
  - Avoid duplicate team or project creation
- api:
  - Enforce not null data
  - Refactoring postgresql queries
  - Update assets requirements
  - Update execution status
- ui:
  - Update packages
  - Update alert message when team or project exist
  - Fix team name in project

## [v0.1.0](https://github.com/Lord-Y/cypress-parallel/releases/tag/v0.1.0) - 2021-06-02

v0.1.0 release

### Added
- Add duplicate project feature
  - duplicate project
  - duplicate annotations
  - duplicate environments variables
- Add githook for pre-commit and pre-push

### Changed
- Fix css classes in executions
- Fix spec icon color
- Update default version of cypress docker image
- Update execution status

## [v0.1.0-beta1](https://github.com/Lord-Y/cypress-parallel/releases/tag/v0.1.0-beta1) - 2021-06-02

v0.1.0 release

### Added
- Add duplicate project feature
  - duplicate project
  - duplicate annotations
  - duplicate environments variables
- Add githook for pre-commit and pre-push

### Changed
- Fix css classes in executions
- Update default version of cypress docker image

## [v0.0.1](https://github.com/Lord-Y/cypress-parallel/releases/tag/v0.0.1) - 2021-05-28

Initial release version

### Added
- create, update and delete teams
- create, update and delete projects and also trigger new unit testing
- create, update and delete annotations that will be used by pods
- create, update and delete environments variables that will be used by pods
- View execution results of units testing

## [v0.0.1-beta1](https://github.com/Lord-Y/cypress-parallel/releases/tag/v0.0.1-beta1) - 2021-05-28

Initial beta version

### Added
- create, update and delete teams
- create, update and delete projects and also trigger new unit testing
- create, update and delete annotations that will be used by pods
- create, update and delete environments variables that will be used by pods
- View execution results of units testing
