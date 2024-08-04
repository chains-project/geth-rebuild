#!/bin/sh

DOCKER_TAG=$1
BIN_DIR=$2


if [ -z "$DOCKER_TAG"  ] || [ -z "$BIN_DIR" ]; then
  echo "Usage: $0 <docker tag> <bin dir>"
  exit 1
fi


# Run container in detached mode
# Runs comparison script on start
echo
echo "Running container $DOCKER_TAG in detached mode"

CONTAINER_ID=$(docker run -d "$DOCKER_TAG") || { echo "failed to start container $DOCKER_TAG"; exit 1; }
docker logs -f "$CONTAINER_ID"  # Capture docker compare script output


# Docker script status
# 42: binaries match
# 43: binaries do not match
# 44: error: no binaries procuced

SCRIPT_EXIT_STATUS=$(docker wait "$CONTAINER_ID") 

# If no binaries produced, abort
if [ "$SCRIPT_EXIT_STATUS" -eq 2 ]; then
    echo "Aborting" 
    exit 44
fi


# Set up clean binary dir
mkdir -p "$BIN_DIR"

REF_BIN="$BIN_DIR/geth-reference"
REP_BIN="$BIN_DIR/geth-reproduce"

if [ -f "$REF_BIN" ]; then
    rm "$REF_BIN" || { echo "failed to rm $$REF_BIN"; exit 1; }
fi

if [ -f "$REP_BIN" ]; then
    rm "$REP_BIN" || { echo "failed to rm $$REP_BIN"; exit 1; }
fi


# Copy produced binaries to local machine
echo "Copying produced binaries to $BIN_DIR..."

docker cp "$CONTAINER_ID:/bin/geth-reference" "$BIN_DIR"  || { echo "failed to copy /bin/geth-reference to $BIN_DIR"; exit 1; }
docker cp "$CONTAINER_ID:/bin/geth-reproduce" "$BIN_DIR" ||   { echo "failed to copy /bin/geth-reference to $BIN_DIR"; exit 1; }

# Stop container
echo
echo "Stopping container $CONTAINER_ID"
docker stop "$CONTAINER_ID" || { echo "error: container id not found"; exit 1; }

# Log info
echo "You can run it again with docker run -it $DOCKER_TAG"
echo "Remove the container with docker rm $CONTAINER_ID"


if [ "$SCRIPT_EXIT_STATUS" -eq 0 ]; then
    exit 0
fi

if [ "$SCRIPT_EXIT_STATUS" -eq 1 ]; then
    exit 1
fi