# Change Log
All notable changes to this project will be documented in this file.


## [2.0.1] - 2025-04-23

### Added

- Support for URLs in emails in addition to QR codes
- Support for custom Entra configuration (client id, scope, tenant)
- Support for a custom User Agent on backend requests to Microsoft

### Fixed

- Bug in the send email loop that, instead of sending each email to one recipient at a time, would send each email to *all* recipients


## [2.0.0] - 2025-04-22

### Added

- Web front end to interact with SquarePhish
- Database to store configuration, credentials, and metrics

### Changed

- Full codebase rewrite to Golang