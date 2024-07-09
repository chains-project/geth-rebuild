#!/bin/sh
set -e

if ! docker info > /dev/null 2>&1; then
    echo "Docker is not running. Opening Docker..."
    open -a Docker || {
        echo "Failed to start Docker."
        exit 1
    }

    start=$(date +%s)
    timeout=$((start + 75))

    while ! docker info > /dev/null 2>&1; do
        echo "Waiting for Docker to start..."
        now=$(date +%s)
        if [ $now -gt $timeout ]; then
            echo "Timed out waiting for Docker to start."
            exit 1
        fi
        sleep 5
    done
fi
echo "Docker is running."
