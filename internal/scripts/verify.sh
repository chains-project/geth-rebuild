#!/bin/sh

DOCKER_TAG=$1
BIN_DIR=$2
LOG_DIR=$3


if [ -z "$DOCKER_TAG"  ] || [ -z "$BIN_DIR" ] || [ -z "$LOG_DIR" ]; then
  echo "Usage: $0 <docker tag> <bin dir> <log dir>"
  exit 1
fi


# Run container in detached mode
# Runs comparison script on start
echo
echo "Running container $DOCKER_TAG in detached mode"

CONTAINER_ID=$(docker run -d "$DOCKER_TAG") || { echo "failed to start container $DOCKER_TAG"; exit 1; }
docker logs -f "$CONTAINER_ID"  # Capture docker compare script output


# Docker script status
# 0: binaries match
# 1: binaries do not match
# 2: error: no binaries procuced

SCRIPT_EXIT_STATUS=$(docker wait "$CONTAINER_ID") 

# If no binaries produced, abort
if [ "$SCRIPT_EXIT_STATUS" -eq 2 ]; then
    echo "Aborting" 
    exit 2
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
mkdir -p "$BIN_DIR/$DOCKER_TAG"
echo "Copying produced binaries to $BIN_DIR/$DOCKER_TAG"

docker cp "$CONTAINER_ID:/bin/geth-reference" "$BIN_DIR/$DOCKER_TAG"  || { echo "failed to copy /bin/geth-reference to $BIN_DIR/$DOCKER_TAG"; exit 1; }
docker cp "$CONTAINER_ID:/bin/geth-reproduce" "$BIN_DIR/$DOCKER_TAG" ||   { echo "failed to copy /bin/geth-reference to $BIN_DIR/$DOCKER_TAG"; exit 1; }

# Stop container
echo
echo "Stopping container $CONTAINER_ID"
docker stop "$CONTAINER_ID" || { echo "error: container id not found"; exit 1; }
echo "Remove the container with docker rm $CONTAINER_ID"
echo

# Log info
echo "You can run the verification again with 'docker run -it $DOCKER_TAG'"
echo 


mkdir -p "$LOG_DIR"

if [ "$SCRIPT_EXIT_STATUS" -eq 0 ]; then
    echo "{\"image\": \"$DOCKER_TAG\", \"status\": \"match\", \"cid\": \"$CONTAINER_ID\"}" > "$LOG_DIR/$DOCKER_TAG.json"
    exit 0
fi

if [ "$SCRIPT_EXIT_STATUS" -eq 1 ]; then
    echo "{\"image\": \"$DOCKER_TAG\", \"status\": \"mismatch\", \"cid\": \"$CONTAINER_ID\"}" > "$LOG_DIR/$DOCKER_TAG.json"
    exit 0
fi