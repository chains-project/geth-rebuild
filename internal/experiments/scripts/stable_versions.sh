#!/bin/bash

# Check if at least one version number is provided
if [ $# -eq 0 ]; then
  echo "No version numbers provided."
  exit 1
fi

# Initialize the output JSON file


output_file=$(PWD)"/stable_version_commits.json"
echo "[" > $output_file
cd ~/go-ethereum

# Iterate over each version number provided as an argument
for version in "$@"; do
  # Checkout the version
  git checkout "v$version" &> /dev/null
  
  # Get the latest commit hash for the version
  commit_hash=$(git log --format=%H -n 1)

  # Append the version and commit hash to the JSON file
  echo "  {" >> $output_file
  echo "    \"version\": \"$version\"," >> $output_file
  echo "    \"commit\": \"$commit_hash\"" >> $output_file
  echo "  }," >> $output_file
done

# Remove the trailing comma from the last JSON object and close the array
echo "]" >> $output_file

# Output success message
echo "Version and commit information has been saved to $output_file."
