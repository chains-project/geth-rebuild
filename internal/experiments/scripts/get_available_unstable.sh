#!/bin/sh

COMMIT_1=$1
COMMIT_2=$2
VERSION_OLDER=$3
VERSION_NEWER=$4
GETH_DIR=$5

if [ -z "$COMMIT_1"  ] || [ -z "$COMMIT_2" ] || [ -z "$GETH_DIR" ] || [ -z "$VERSION_OLDER" ] || [ -z "$VERSION_NEWER" ]; then
  echo "Usage: $0 <from commit> <to commit> <version older> <version newer> <geth dir>"
  exit 1
fi


#   ff6e43e8 6eb42a6b

cd "$GETH_DIR" || { echo "failed changing directory to $GETH_DIR"; exit 1; }
COMMITS=$(git log "$COMMIT_1..$COMMIT_2" --format=%H) || { echo "failed retrieving git logs 1"; exit 1; }


git checkout "v$VERSION_OLDER"

LAST_COMMIT_OLDER=$( git log --format=%H -n 2 | tail -n 1 ) || { echo "failed retrieving git logs 2"; exit 1; }

json_output="["


for commit in $COMMITS; do  
  is_ancestor=$(git merge-base --is-ancestor "$commit" "$LAST_COMMIT_OLDER" && echo "yes")

  # Compare commit with LAST_COMMIT_OLDER to determine version assignment
  if [ "$commit" = "$LAST_COMMIT_OLDER" ]; then
    version=$VERSION_OLDER
  elif [ "$is_ancestor" = "yes" ]; then
    version=$VERSION_OLDER
  else
    version=$VERSION_NEWER
  fi

  if $first_item; then
    json_output="$json_output\n    {\"commit\": \"$commit\", \"version\": \"$version\"}"
    first_item=false
  else
    json_output="$json_output,\n    {\"commit\": \"$commit\", \"version\": \"$version\"}"
  fi
done

json_output="$json_output\n]"

output_file=~/geth-rebuild/internal/experiments/data/available_unstable_commits.json
echo "$json_output" > "$output_file" || { echo "Failed writing to $output_file"; exit 1; }

echo "JSON data has been written to $output_file"