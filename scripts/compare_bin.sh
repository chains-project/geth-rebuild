#!/bin/sh

set -e

BIN_REF=$1
BIN_REP=$1

if [ -z "$BIN_REF"  ] || [ -z "$BIN_REP" ]; then
  echo "Usage: $0 <bin> <bin>"
  exit 1
fi


OS=$(uname)
if [ "$OS" = "Linux" ]; then
    md5_reference=$(md5sum "$BIN_REF" | awk '{print $1}')
    md5_reproduce=$(md5sum "$BIN_REP" | awk '{print $1}')
else
    if [ "$OS" = "Darwin" ]; then
    md5_reference=$(md5 "$BIN_REF" | awk '{print $NF}')
    md5_reproduce=$(md5 "$BIN_REP" | awk '{print $NF}')
    else 
        echo "OS $OS not supported." && exit 1
    fi
fi

echo "First build has hash $md5_reference"
echo "Second build has hash $md5_reproduce"

if [ "$md5_reproduce" != "$md5_reference" ]; then
    echo "Binaries mismatch. Running diffoscope..."
    # diff the binaries -- this takes a while.
    cd "$BIN_DIR" || exit 1
    docker run --rm -t -w "$(pwd)" -v "$(pwd)":"$(pwd)":rw registry.salsa.debian.org/reproducible-builds/diffoscope --progress geth-reference geth-reproduce
else
    if [ "$md5_reproduce" = "" ]; then
        { echo "Error: no binary produced."; exit 1; }
    else
        echo "Binaries match."
    fi
fi
