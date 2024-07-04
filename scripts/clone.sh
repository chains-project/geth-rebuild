#!/bin/sh

set -e

DIR=$1
GETH_DIR="$DIR/go-ethereum"
URL="https://github.com/ethereum/go-ethereum.git"
BRANCH="master"

if [ -z "$DIR" ]; then
  echo "Usage: $0 <target directory>"
  exit 1
fi



echo "Cloning go ethereum branch $BRANCH from $URL"
mkdir -p "$DIR"

if [ -d "$GETH_DIR" ]; then
    rm -rf "$GETH_DIR"
fi

git clone -v --branch $BRANCH $URL "$GETH_DIR" || { echo "Failed to clone Geth sources."; exit 1; } # TODO shallow copy. Decide proper --depth OR use --single-branch
