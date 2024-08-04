#!/bin/sh

set -e

BIN_REF=$1
BIN_REP=$2

if [ -z "$BIN_REF"  ] || [ -z "$BIN_REP" ]; then
  echo "Usage: $0 <artifact 1> <artifact 2>"
  exit 1
fi

echo && echo "Comparing SHA256 of binary artefacts"

OS=$(uname)
if [ "$OS" = "Linux" ]; then
    md5_reference=$(sha256sum "$BIN_REF" | awk '{print $1}')
    md5_reproduce=$(sha256sum "$BIN_REP" | awk '{print $1}')
else
    if [ "$OS" = "Darwin" ]; then
    md5_reference=$(shasum -a 256 "$BIN_REF" | awk '{print $1}')
    md5_reproduce=$(shasum -a 256 "$BIN_REP" | awk '{print $1}')
    else 
        echo "OS $OS not supported" && exit 1
    fi
fi

echo "Reference build:      sha256    $md5_reference"
echo "Reproducing build:    sha256    $md5_reproduce" && echo

# TODO tests for exit status

if [ "$md5_reproduce" != "$md5_reference" ]; then
    echo "binaries do not match"
    exit 1
else
    if [ "$md5_reproduce" = "" ]; then
        echo "error: no binary produced"
        exit 2
    else
        echo "binaries match"
        exit 0
    fi
fi
