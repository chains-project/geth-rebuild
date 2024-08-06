# geth-rebuild

A rebuilder for [geth](https://github.com/ethereum/go-ethereum/).

Reproduce and verify source-to-binary semantics of a geth binary artifact.

## Usage

`go build ./cmd/gethrebuild -o ./gethrebuild`

`./gethrebuild <os> <arch> <version>`

For example, `gethrebuild linux-amd64 1.14.3`

See command documentation for optional arguments `gethrebuild --help`

> [!NOTE] Must be run inside project directory

## Cases of Unreproducibility

When reproducing an artifact, cases of non-determinism need to be controlled.

In `.non-determinism/reports` four cases found for geth are shown:

- **buildid**: embedding of unreproducible build ids
- **date**: conditional embedding of release date
- **path**: embedding of absolute system paths
- **gcc**: differing gcc versions using identical build settings

> [!TIP]
> Reproduce these cases by running #TODO

## Limitations

Supported os/arch pairs:

- Linux

  - amd64 ✅
  - 386   ✅
  - arm5  ✅
  - arm6  ✅
  - arm7  ✅
  - arm64 ✅

- Not supported currently
  - Darwin  ❌
  - Windows ❌
