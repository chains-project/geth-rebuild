#!/bin/sh
set -e

OS=$(uname)

if [ "$OS" = "Linux" ]; then
    CMD="dockerd &"
elif [ "$OS" = "Darwin" ]; then
    # docker desktop...
    CMD="open -a Docker"
else
    echo "Unsupported operating system: $OS"
    exit 1
fi


if ! docker info > /dev/null 2>&1; then
    echo "Docker is not running. Opening Docker..."
    eval "$CMD" || {
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
        sleep 4
    done
fi
echo "Docker is running."
