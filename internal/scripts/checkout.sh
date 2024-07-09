#!/bin/sh

set -e

GETH_DIR="$1"
VERSION_OR_COMMIT=$2

if [ -z "$GETH_DIR" ] || [ -z "$VERSION_OR_COMMIT" ]; then
  echo "Usage: $0 <directory> <version or commit>"
  exit 1
fi

echo "Checking out go-ethereum at version $VERSION_OR_COMMIT"
cd "$GETH_DIR" || { echo "Failed cd to $GETH_DIR"; exit 1; }

git fetch --quiet || { echo "Failed to fetch geth"; exit 1; }
git checkout --quiet "v$VERSION_OR_COMMIT" || { echo "Failed to checkout to $VERSION_OR_COMMIT."; echo "Does version exist?"; exit 1; }
