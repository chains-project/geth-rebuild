#!/bin/sh

set -e

DOCKER_TAG=$1
OUTPUT_DIR=$2

if [ -z "$DOCKER_TAG"  ] || [ -z "$OUTPUT_DIR" ]; then
  echo "Usage: $0 <docker tag> <output dir>"
  exit 1
fi

mkdir -p "$OUTPUT_DIR"

# run container in detached mode
echo "Running container $DOCKER_TAG in detached mode."
CONTAINER_ID=$(docker run -d "$DOCKER_TAG" /bin/sh) # cannot use --rm here: loses cid

# copy binaries and stop container
echo "Copying binaries..."
docker cp -q "$CONTAINER_ID":/bin/geth-reference "$OUTPUT_DIR"
docker cp -q "$CONTAINER_ID":/bin/geth-reproduce "$OUTPUT_DIR"
docker stop "$CONTAINER_ID"


echo "Stopped container $CONTAINER_ID"
echo "You can run it again with docker run $CONTAINER_ID /bin/sh"
echo "Remove the conotainer with docker rm $CONTAINER_ID"
echo ""
echo "Binaries copied to $OUTPUT_DIR"