#!/bin/bash

set -e

# 33 unstable
NUM_COMMITS=$1

if [ -z "$NUM_COMMITS" ]; then
  echo "Usage: $0 <retrieve # commits>"
  exit 1
fi


cd ~/geth-rebuild/tmp/go-ethereum

git fetch
git checkout master
git pull

COMMITS=$(git log --format="%H" -n "$NUM_COMMITS")
OUT=~/geth-rebuild/internal/experiments/data/20_latest_commits.json

if [ -e $OUT ]; then
  rm $OUT
fi

VERSION="v1.14.8"

json_output="{\"commits\": ["

first=true
for COMMIT in $COMMITS; do
    if [ "$first" = true ]; then
        first=false
    else
        json_output+=","
    fi

  json_output+="{\"commit\":\"$COMMIT\",\"version\":\"$VERSION\"}"
done

json_output+="]}"
echo "$json_output" >> $OUT
