# k8s-serviceaccount-lib Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v2.0.1] - 2026-06-19
### Fixed
- Bump the Kubernetes CRD API version of `ServiceAccountRequest` and `ServiceAccountProducer` from `k8s.cloudogu.com/v1` to `k8s.cloudogu.com/v2`.
- Add the required `/v2` major-version suffix to the Go module path (`module github.com/cloudogu/k8s-serviceaccount-lib/v2`).
  
## [v2.0.0] - 2026-06-19

### Changed
- [#3] Improve status conditions for ServiceAccountRequests
- [#3] Simplify serviceaccount-request params to just a map of strings
  - When creating service accounts, the params are passed to the serviceaccount-producer via the serviceaccount operator

## [v1.0.0] - 2026-04-22
### Added
- [#1] Initial CRDs for ServiceAccountRequest and ServiceAccountProducer