# geth-rebuild

A rebuilder for geth

# State

What has been bit-for-bit reproduced?

## Locally reproduced Travis pipeline

- Linux binary releases

  - `/cmd/geth`
    - amd64 ✅
    - 386 ✅
    - arm5 ✅
    - arm6 🟡 (assuming reproducible)
    - arm7 🟡
    - arm64 🟡
  - `/cmd/*` releases ❌

- OSX releases ❌
- Windows releases ❌ (probably will not attempt)
- Docker images ❌
- ubuntu PPAs, homebrew et

# Rebuilding a geth binary

## Build Inputs

When reproducing a geth binary, we need the correct **source code** and **build configurations** to reproduce the binary.

- **Q:** can we just take in a geth binary and get all information? _Probably..._ except ubuntu dist? Check if can reproduce regardless of underlying ubuntu dist.

Given a certain...

- `GETH_VERSION`: E.g. 1.14.0 or 1.14.1-unstable
- Target OS (TODO: currently only Linux)
- `TARGET_ARCH`: target architecture, e.g. amd-65

We need to fetch the following information:

- `GETH_PKG`: relevant package (`geth` vs `geth-alltools`)
  - **How:** TODO
- `GETH_COMMIT`: geth commit given version
  - **How:** TODO
- `GO_FLAGS`: additional go flags needed for build
  - **How:** TODO, found in travis.yml
- GCC version
  - With `readelf -p .comment geth`
