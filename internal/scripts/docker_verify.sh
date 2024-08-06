#!/bin/sh

DOCKER_TAG=$1
BIN_DIR=$2
LOG_DIR=$3


if [ -z "$DOCKER_TAG"  ] || [ -z "$BIN_DIR" ] || [ -z "$LOG_DIR" ]; then
  echo "Usage: $0 <docker tag> <bin dir> <log dir>"
  exit 1
fi


# Run container in detached mode
# Runs comparison script on start (see dockerfile)
echo && echo "Running container $DOCKER_TAG in detached mode"

CONTAINER_ID=$(docker run -d "$DOCKER_TAG") || { echo "failed to run container $DOCKER_TAG"; exit 1; }
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

# Write rebuild result log
mkdir -p "$LOG_DIR"
LOG_TO="$LOG_DIR/$DOCKER_TAG.json"

RESULT=""
if [ "$SCRIPT_EXIT_STATUS" -eq 0 ]; then
    RESULT="match"
elif [ "$SCRIPT_EXIT_STATUS" -eq 1 ]; then
    RESULT="mismatch"
else
    echo "error: unexpected script exit code: $SCRIPT_EXIT_STATUS" 
    exit 1
fi

# Write
# TODO should include the SHA256 values here too...?
echo "{\"image\": \"$DOCKER_TAG\", \"status\": \"$RESULT\", \"cid\": \"$CONTAINER_ID\"}" \
    > "$LOG_TO" || { echo "failed to write log to $LOG_TO"; exit 1; }


# Copy produced binaries to local machine
COPY_TO="$BIN_DIR/$DOCKER_TAG"
mkdir -p "$COPY_TO"

echo "Copying produced binaries to $COPY_TO"
docker cp "$CONTAINER_ID:/bin/geth-reference" "$COPY_TO"  || { echo "failed to copy /bin/geth-reference to $COPY_TO"; exit 1; }
docker cp "$CONTAINER_ID:/bin/geth-reproduce" "$COPY_TO" ||  { echo "failed to copy /bin/geth-reference to $COPY_TO"; exit 1; }


# Stop container
echo && echo "Stopping container $CONTAINER_ID"
docker stop "$CONTAINER_ID" || { echo "failed to stop container $CONTAINER_ID"; exit 1; }
echo "Remove the container with docker rm $CONTAINER_ID" && echo


# bye bye log info
echo "You can run the verification again with 'docker run -it $DOCKER_TAG'" && echo 
