#!/bin/sh

set -e

DIR=$1
GETH_DIR="$DIR/go-ethereum"
VERSION=$2

if [ -z "$DIR" ] || [ -z "$VERSION" ]; then
  echo "Usage: $0 <git-dir> <version>"
  exit 1
fi

echo "[CLONING GO ETHEREUM SOURCES]"
mkdir -p "$DIR"
rm -r -f "$GETH_DIR"
git clone --branch master https://github.com/ethereum/go-ethereum.git "$GETH_DIR" # TODO shallow copy. Decide proper --depth OR use --single-branch
cd "$GETH_DIR" || { echo "Failed cd to $GETH_DIR"; exit 1; }

git fetch --quiet
git checkout --quiet "v$VERSION" || { echo "Failed to checkout to version $VERSION."; echo "Does version exist?"; exit 1; }
