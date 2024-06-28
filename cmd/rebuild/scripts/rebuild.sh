#!/bin/sh
set -e

# Required environment variables
: "${UBUNTU_VERSION:?}"
: "${GO_VERSION:?}"
: "${GETH_VERSION:?}"
: "${GETH_COMMIT:?}"
: "${OS_ARCH:?}"
: "${REFERENCE_URL:?}"
: "${PACKAGES:?}"
: "${C_COMPILER:?}"
: "${TAG:?}"
: "${BUILD_CMD:?}"
: "${ELF_TARGET:?}"
: "${DOCKER_FILE:?}"

docker build -t "$TAG" \
  --build-arg UBUNTU_VERSION="${UBUNTU_VERSION}" \
  --build-arg GO_VERSION="${GO_VERSION}" \
  --build-arg GETH_VERSION="${GETH_VERSION}" \
  --build-arg GETH_COMMIT="${GETH_COMMIT}" \
  --build-arg OS_ARCH="${OS_ARCH}" \
  --build-arg REFERENCE_URL="${REFERENCE_URL}" \
  --build-arg PACKAGES="${PACKAGES}" \
  --build-arg C_COMPILER="${C_COMPILER}" \
  --build-arg BUILD_CMD="${BUILD_CMD}" \
  --build-arg ELF_TARGET="${ELF_TARGET}" \
  - < "$DOCKER_FILE"