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

validate () {
    os_arch_pattern="^linux-(amd64|386|arm5|arm6|arm64|arm7)$"
    version_pattern="^[0-9]+\.[0-9]+\.[0-9]+$"
    if ! [[ $1 =~ $os_arch_pattern ]]; then
        error "<os-arch> must be a valid linux target architecture\nExample: linux-amd64"
    fi
    if ! [[ $2 =~ $version_pattern ]]; then
        error "<geth version> must be in format 'major.minor.patch'\nExample: 1.14.4"
    fi
}


get_arch_id () {
    arch=$(echo "$OS_ARCH" | cut -d'-' -f2)

    if [[ $arch == arm[0-9] ]]; then
        arm_version=$(echo "$arch" | sed 's/[^0-9]*//g')
        echo "ARM=$arm_version"
    else
        echo "$arch"
    fi
}

# Variable intialization and validation. Clean up.

OS_ARCH=$1
GETH_VERSION=$2
validate "$OS_ARCH" "$GETH_VERSION"
mkdir -p tmp
#rm -rf ./tmp/go-ethereum


# Fetch source code and check out at given version.
echo -e "\n\n[FETCHING GO ETHEREUM]\n"


cd tmp && echo -e "Downloading into directory $PWD\n" || error "Failed 'cd' to directory 'tmp'."
#git clone --branch master https://github.com/ethereum/go-ethereum.git || error "Failed cloning into go-ethereum"
cd go-ethereum && git fetch && git checkout --quiet "v$GETH_VERSION" || error "Failed checking out at version $GETH_VERSION\nDoes tag exist?"


# Construct URL for fetching binary.

echo -e "\n\n[CONSTRUCTING URL FOR BINARY DOWNLOAD]\n"
VERSION_COMMIT=$(git log -1 --format=%H)
SHORT_COMMIT=${VERSION_COMMIT:0:8}
TARGET_PKG="geth-$OS_ARCH-$GETH_VERSION-$SHORT_COMMIT"
URL="https://gethstore.blob.core.windows.net/builds/$TARGET_PKG.tar.gz"
echo -e "Version: ${GETH_VERSION}\nCommit: $VERSION_COMMIT ($SHORT_COMMIT)\nOS-ARCH: $OS_ARCH\n\n>>[URL] $URL\n"



echo -e "\n\n[RETRIEVING BUILD COMMANDS]\n"
ARCH_ID=$(get_arch_id)
# CONTINUE HERE @ REGEXP
grep -o 'go run build/ci.go install.*' .travis.yml | grep "$ARCH_ID" | sed 's/go run.*-dlgo/-dlgo/' | sed 's/- //'  # todo filter out only linux commands
