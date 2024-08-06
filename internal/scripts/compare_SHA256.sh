#!/bin/sh

set -e

BIN_REF=$1
BIN_REP=$2

if [ -z "$BIN_REF" ] || [ -z "$BIN_REP" ]; then
  echo "Usage: $0 <artifact 1> <artifact 2>"
  exit 1
fi

echo && echo "Comparing SHA256 of binary artefacts..."

OS=$(uname)
if [ "$OS" = "Linux" ]; then
    REF_SHA=$(sha256sum "$BIN_REF" | awk '{print $1}')
    REP_SHA=$(sha256sum "$BIN_REP" | awk '{print $1}')
elif [ "$OS" = "Darwin" ]; then
    REF_SHA=$(shasum -a 256 "$BIN_REF" | awk '{print $1}')
    REP_SHA=$(shasum -a 256 "$BIN_REP" | awk '{print $1}')
else 
    echo "OS $OS not supported" && exit 1
fi

echo "Reference build:      SHA256    $REF_SHA"
echo "Reproducing build:    SHA256    $REP_SHA" && echo


# TODO tests for exit status
if [ -z "$REF_SHA" ] || [ -z "$REP_SHA" ]; then
    echo && echo "error producing one or both artifacts" && echo && echo
    exit 2
elif [ "$REF_SHA" = "$REP_SHA" ]; then
    echo && echo ">>>[REPRODUCTION SUCCESSFUL] Binaries match." && echo && echo
    exit 0
else
    echo && echo ">>>[REPRODUCTION FAILED] Binaries do not match." && echo && echo
    exit 1
fi
