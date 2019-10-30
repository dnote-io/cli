# CHANAGELOG

All notable changes to the projects under this repository will be documented in this file.

* [Server](#server)
* [CLI](#cli)
* [Browser Extensions](#browser-extensions)

## Server

The following log documentes the history of the server project.

### [Unreleased]

#### Upgrade Guide

* Please define a new environment variable `WebURL` whose value is the URL to your Dnote server, without the trailing slash. (e.g. `https://my-server.com`) (Please see #290)

#### Fixed

- Allow to customize the app URL in the emails (#290)

### 0.2.0 - 2019-10-28

#### Added

- Specify spaced repetition rule (#280)

#### Changed

- Treat a linebreak as a new line in the preview (#261)
- Allow to have multiple editor states for adding and editing notes (#260)

#### Fixed

- Fix jumping focus on editor (#265)

### 0.1.1 - 2019-09-30

#### Fixed

- Fix asset loading (#257)


### 0.1.0 - 2019-09-30

#### Added

- Full-text search (#254)
- Password recovery (#254)
- Embedded notes in the digest emails (#254)

#### Removed

- **Breaking Change**: End-to-end encryption was removed. Existing users need to go to `/classic` and follow the automated migration steps. (#254)
- **Breaking Change**: `v1` and `v2` API endpoints were removed, and `v3` API was added as a replacement.

#### Migration guide

- In your application, navigate to `/classic` and follow the automated migration steps.


## CLI

The following log documentes the history of the CLI project

### 0.10.0 - 2019-09-30

#### Removed

- **Breaking Change**: End-to-end encryption was removed. Previous versions will no longer be able to interact with the web API, because `v1` and `v2` endpoints were replaced by a new `v3` endpoint to remove encryption.

#### Migration guide

- If you are using Dnote Pro, change the value of `apiEndpoint` in `~/.dnote/dnoterc` to `https://api.getdnote.com`.

## Browser Extensions

The following log documentes the history of the browser extensions project

### [Unreleased]

N/A

### 2.0.0 - 2019-10-29

- Allow to customize API and web URLs (#285)

### 1.1.1 - 2019-10-02

- Fix failing requests (#263)

### 1.1.0 - 2019-09-30

#### Removed

- **Breaking Change**: End-to-end encryption was removed. Previous versions will no longer be able to interact with the web API, because `v1` and `v2` endpoints were replaced by a new `v3` endpoint to remove encryption.
