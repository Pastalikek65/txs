# Changelog

All notable changes to `txs` will be documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-08-27
### Added
- 5 panels: github (Link pagination, PathEscape), jobs (push/retry, webhook), db (ListTables), sys (/proc), files (ls)
- Single binary 9.9M, `CGO_ENABLED=0`, `store` WAL 5 tables, `tui` tab router
- 12 tests <0.5s, `go vet` PASS
