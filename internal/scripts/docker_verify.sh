#!/bin/sh

DOCKER_TAG=$1
STORE_BINS=$2
LOG_FILE=$3


if [ -z "$DOCKER_TAG"  ] || [ -z "$STORE_BINS" ] || [ -z "$LOG_FILE" ]; then
  echo "Usage: $0 <docker tag> <unique bin dir> <json log file>"
  exit 1
fi

if [ ! -f "$LOG_FILE" ]; then
    echo "error: no log file found at $LOG_FILE"
    exit 1
fi

# Run container in detached mode
# Runs comparison script on start (see dockerfile)
echo && echo "Running container $DOCKER_TAG in detached mode"

CONTAINER_ID=$(docker run -d "$DOCKER_TAG") || { echo "failed to run container $DOCKER_TAG"; exit 1; }
docker logs -f "$CONTAINER_ID" || { echo "error retrieving docker logs"; exit 1; } # Capture docker compare script output


# Docker script status
# 0: binaries match
# 1: binaries do not match
# 2: error: no binaries procuced

SCRIPT_EXIT_STATUS=$(docker wait "$CONTAINER_ID") 


# Determine result status of rebuild
STATUS=""
if [ "$SCRIPT_EXIT_STATUS" -eq 0 ]; then
    STATUS="match"
elif [ "$SCRIPT_EXIT_STATUS" -eq 1 ]; then
    STATUS="mismatch"
elif [ "$SCRIPT_EXIT_STATUS" -eq 2 ]; then
    STATUS="error"
else
    STATUS="error"
    echo "error: unexpected script exit code: $SCRIPT_EXIT_STATUS" 
fi

# Write result to json
jq --arg key "STATUS" --arg value $STATUS \
'.[$key] = $value' "$LOG_FILE" > tmp.$$.json && mv tmp.$$.json "$LOG_FILE" \
|| { echo "failed to write log to $LOG_FILE"; exit 1; }

# Write contianer id to json
jq --arg key "CONTAINER_ID" --arg value "$CONTAINER_ID" \
'.[$key] = $value' "$LOG_FILE" > tmp.$$.json && mv tmp.$$.json "$LOG_FILE" \
|| { echo "failed to write log to $LOG_FILE"; exit 1; }

# If no binaries produced, abort
if [ "$STATUS" = "error" ]; then
    echo "Aborting"
    exit 2
fi

# Copy produced binaries to local machine
mkdir -p "$STORE_BINS" # TODO MOve

echo "Copying produced binaries to $STORE_BINS"
docker cp "$CONTAINER_ID:/bin/geth-reference" "$STORE_BINS"  || { echo "failed to copy /bin/geth-reference to $STORE_BINS"; exit 1; }
docker cp "$CONTAINER_ID:/bin/geth-reproduce" "$STORE_BINS" ||  { echo "failed to copy /bin/geth-reference to $STORE_BINS"; exit 1; }


# Stop container
echo && echo "Stopping container $CONTAINER_ID"
docker stop "$CONTAINER_ID" || { echo "failed to stop container $CONTAINER_ID"; exit 1; }
echo "Remove the container with docker rm $CONTAINER_ID" && echo


# bye bye log info
echo "You can run the verification again with 'docker run -it $DOCKER_TAG'" && echo 
