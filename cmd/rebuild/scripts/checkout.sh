#!/bin/sh

set -e

GETH_DIR="$1"
GETH_VERSION=$2

if [ -z "$GETH_DIR" ] || [ -z "$GETH_VERSION" ]; then
  echo "Usage: $0 <geth dir> <version>"
  exit 1
fi

echo "Checking out go-ethereum at version $GETH_VERSION"
cd "$GETH_DIR" || { echo "Failed cd to $GETH_DIR"; exit 1; }

git fetch --quiet
git checkout --quiet "v$GETH_VERSION" || { echo "Failed to checkout to version $GETH_VERSION."; echo "Does version exist?"; exit 1; }
