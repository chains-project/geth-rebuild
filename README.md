# geth-rebuild

A rebuilder for [geth](https://github.com/ethereum/go-ethereum/).

Reproduce and verify source-to-binary semantics of a geth Linux binary.

## Build from source

`go build ./cmd/rebuild -o ./bin/rebuild`

## Usage

`cd ./bin`

`./rebuild <os-arch> <version>`

For example, `./rebuild linux-amd64 1.14.3`


## Cases of Unreproducibility 

When reproducing an artifact, cases of non-determinism need to be controlled.

In `./reports` four cases found for geth are shown:

- **buildid**: embedding of unreproducible build ids
- **date**: conditional embedding of release date
- **path**: embedding of absolute system paths
- **gcc**: differing gcc versions using identical build settings

> [!TIP]
> Reproduce these cases using `TODO PATH <case> <docker tag>`
> E.g. `TODO path my-path-tag`


## Build Inputs

When reproducing a geth binary, we need the correct **source code** and **build configurations** to reproduce the binary.

Given a certain...

- `GETH_VERSION`: E.g. 1.14.0 or 1.14.1-unstable
- `TARGET_ARCH`: target architecture, e.g. linux amd64
- `GETH_PKG`: relevant package (`geth` vs `geth-alltools`)

We need to fetch the following information:

- `GETH_COMMIT`: geth commit given version
  - **How:** `go version -m geth`
- `GO_VERSION`: Go compiler version
  - **How:** `go version -m geth`
- `BUILD_FLAGS`: additional go flags needed for build
  - **How:** get from travis.yml
- GCC version
  - **How:** `readelf -p .comment geth`


## Limitations

What can be bit-for-bit reproduced?

Linux`/cmd/geth` binaries for Linux
- amd64  ✅
- 386    ✅
- arm5   ✅
- arm6   ✅
- arm7   ✅
- arm64  ✅

Limitations/Not attempted ❌
- Linux geth-alltools releases
- OSX releases
- Windows releases
- Docker images
- ubuntu PPAs, homebrew etc.
