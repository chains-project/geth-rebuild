#!/bin/bash


if [ $# != 2 ]; then
    echo "Usage $0: <os-arch> <geth version>"
    echo "Example: $0 linux-amd64 1.14.3"
    exit 2
fi


# Function declarations.

error () {
    echo -e "Error: $1"
    exit 1
}

validate() {
    os_arch_pattern="^linux-(amd64|386|arm5|arm6|arm64|arm7)$"
    version_pattern="^[0-9]+\.[0-9]+\.[0-9]+$"
    if ! [[ $1 =~ $os_arch_pattern ]]; then
        error "<os-arch> must be a valid linux target architecture\nExample: linux-amd64"
    fi
    if ! [[ $2 =~ $version_pattern ]]; then
        error "<geth version> must be in format 'major.minor.patch'\nExample: 1.14.4"
    fi

}


# clean up.

mkdir -p tmp
rm -rf ./tmp/go-ethereum


# Variable intialization and validation.

OS_ARCH=$1
GETH_VERSION=$2
validate "$OS_ARCH" "$GETH_VERSION"


# Fetch source code and check out at given version.

cd tmp || error "Failed changing to directory 'tmp'."
git clone --branch master https://github.com/ethereum/go-ethereum.git || error "Failed cloning into go-ethereum"
cd go-ethereum && git fetch && git checkout --quiet "v$GETH_VERSION" || error "Failed checking out at version $GETH_VERSION\nDoes tag exist?"


# Construct URL for fetching binary.

VERSION_COMMIT=$(git log -1 --format=%H) && echo "Commit at version $GETH_VERSION: $VERSION_COMMIT"
SHORT_COMMIT="${VERSION_COMMIT:0:8}"
URL="geth-$OS_ARCH-$GETH_VERSION-$SHORT_COMMIT"
echo "URL is: $URL"


