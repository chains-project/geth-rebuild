#!/bin/sh

if [ $# -ne 2 ]; then
    echo "Usage: $0 <non-determinism case (date/gcc/path)> <docker tag>"
    exit 1
fi


ND_CASE=$1
TAG=$2
DOCKER_PATH=./docker/non-determinism/Dockerfile-$ND_CASE
OUTPUT_DIR=./bin
REPORT_DIR=./reports
mkdir -p $OUTPUT_DIR
rm $OUTPUT_DIR/geth-reference $OUTPUT_DIR/geth-reproduce

# build image
echo "Starting docker build..."
docker build -t "$TAG" - < "$DOCKER_PATH"

if [ $? != 0 ]; then
    echo "Error: Docker build failed." && exit 1
fi

# start container
echo "Build finished. Running container in detached mode."
CONTAINER_ID=$(docker run -d "$TAG" /bin/sh) # cannot use --rm here: loses cid

# copy binaries and stop container
echo "Copying binaries..."
docker cp -q "$CONTAINER_ID":/bin/geth-reference "$OUTPUT_DIR"
docker cp -q "$CONTAINER_ID":/bin/geth-reproduce "$OUTPUT_DIR"

# check binary md5s and diff if neq
md5_reference=$(md5 "$OUTPUT_DIR"/geth-reference | awk '{print $NF}')
md5_reproduce=$(md5 "$OUTPUT_DIR"/geth-reproduce | awk '{print $NF}')
echo "First build has hash $md5_reference"
echo "Second build has hash $md5_reproduce"


if [ "$md5_reproduce" != "$md5_reference" ]; then
    echo "Binaries mismatch."
    echo "Writing report to $REPORT_DIR/non-determinism-$ND_CASE.md"
    docker cp -q "$CONTAINER_ID":/non-determinism.md "$REPORT_DIR/non-determinism-$ND_CASE.md"
    echo "You can run diffoscope with 'cd ./bin && docker run --rm -t -w '\$(pwd)' -v '\$(pwd)':'\$(pwd)':rw registry.salsa.debian.org/reproducible-builds/diffoscope --progress geth-reference geth-reproduce'"
else
    if [ "$md5_reproduce" = "" ]; then
        { echo "Error: no binary produced."; exit 1; }
    else
        echo "Binaries match."
    fi
fi

docker stop "$CONTAINER_ID"
docker rm "$CONTAINER_ID"