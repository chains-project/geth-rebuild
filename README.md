# geth-rebuild

A rebuilder for geth

# State

What has been bit-for-bit reproduced?

**Official geth binary releases...**

- `/cmd/geth` binaries for Linux
  - amd64 version 1.14.3 ✅

**In a locally reproduced Travis pipeline...**

- `/cmd/geth` binaries for Linux
  - amd64 ✅
  - 386 ✅
  - arm5 ✅
  - arm6 🟡 (assuming reproducible)
  - arm7 🟡
  - arm64 🟡


- Not attempted ❌
  - Linux geth-alltools releases
  - OSX releases (probably will not attempt)
  - Windows releases (probably will not attempt)
  - Docker images  (probably will not attempt)
  - ubuntu PPAs, homebrew etc.  (probably will not attempt)

# Rebuilding a geth binary

## Build Inputs

When reproducing a geth binary, we need the correct **source code** and **build configurations** to reproduce the binary.

Given a certain...

- `GETH_VERSION`: E.g. 1.14.0 or 1.14.1-unstable
- `TARGET_ARCH`: target architecture, e.g. linux amd64
- `GETH_PKG`: relevant package (`geth` vs `geth-alltools`) **(??)**

We need to fetch the following information:

- `GETH_COMMIT`: geth commit given version
  - **How:** `go version -m geth`
- `GO_VERSION`: Go compiler version
  - **How:** `go version -m geth`
- `BUILD_FLAGS`: additional go flags needed for build
  - **How:** get from travis.yml
- GCC version
  - **How:** `readelf -p .comment geth`
