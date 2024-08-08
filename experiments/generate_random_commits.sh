#!/bin/bash

set -e

NUM_COMMITS=$1

if [ -z "$NUM_COMMITS" ]; then
  echo "Usage: $0 <since commit> <to commit> <retrieve # commits>"
  exit 1
fi

cd ../tmp/go-ethereum

SINCE_COMMIT="71aa15c98f88ee03097e5b30ccbb564734180ca3"
TO_COMMIT="aa55f5ea200dfd07618fdf658d9d2741c3b376a8"

COMMITS=$(git log --format="%H" "$SINCE_COMMIT..$TO_COMMIT")
MAX=$(echo "$COMMITS" | wc -l | xargs)

# we use array and loops to get unique values
declare -a random_indices
while [ ${#random_indices[@]} -lt "$NUM_COMMITS" ]; do
    random_index=$(jot -r 1 0 "$MAX")

    # check if index exists
    if ! [[ "${random_indices[*]}" =~ $random_index ]]; then
        random_indices+=("$random_index")
    fi
done


OUT="../../experiments/random_commitws.txt"

if [ -e $OUT ]; then
  rm $OUT
fi

for index in "${random_indices[@]}"; do
    commit=$(echo "$COMMITS" | sed "${index}q;d")
    version=$(git describe --tags --abbrev=0 "$commit")
    echo "$commit $version" >> $OUT
done