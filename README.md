# geth-rebuild

A rebuilder for [geth](https://github.com/ethereum/go-ethereum/).

Reproduce and verify source-to-binary semantics of a geth Linux binary.

## Build from source

`go build . -o ./bin/rebuild`

## Usage

`cd ./bin`

`./rebuild <os> <arch> <version>`

For example, `./rebuild linux-amd64 1.14.3`


## Cases of Unreproducibility 

When reproducing an artifact, cases of non-determinism need to be controlled.

In `.non-determinism/reports` four cases found for geth are shown:

- **buildid**: embedding of unreproducible build ids
- **date**: conditional embedding of release date
- **path**: embedding of absolute system paths
- **gcc**: differing gcc versions using identical build settings

> [!TIP]
> Reproduce these cases using `TODO PATH <case> <docker tag>`
> E.g. `TODO path my-path-tag`


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
