#!/bin/bash

# Path to the JSON file
json_file=~/geth-rebuild/internal/experiments/data/unstable_versions.json
output_file=~/geth-rebuild/internal/experiments/data/ok_downloads.json

# Initialize the output JSON file
echo "[" > $output_file

# Read the JSON file and extract the version and commit information
commits=$(jq -c '.commits[]' "$json_file")

# Iterate over each commit object
first_entry=true

for commit in $commits; do
  # Extract the version and full commit hash
  version=$(echo "$commit" | jq -r '.version')
  full_commit=$(echo "$commit" | jq -r '.commit')

  # Get the first 8 characters of the commit hash (short commit)
  short_commit=${full_commit:0:8}

  # Construct the download URL
  url="https://gethstore.blob.core.windows.net/builds/geth-linux-amd64-$version-unstable-$short_commit.tar.gz"

  # Download the artifact using curl
  echo "Downloading $url..."
  curl -O "$url"

  # Check if the download was successful
  if [ $? -eq 0 ]; then
    echo "Downloaded artifact for $version ($short_commit) successfully."
    
    # Append the version and commit hash to the output JSON file
    if [ "$first_entry" = true ]; then
      first_entry=false
    else
      echo "," >> $output_file
    fi
    
    echo "  {" >> $output_file
    echo "    \"version\": \"$version\"," >> $output_file
    echo "    \"commit\": \"$full_commit\"" >> $output_file
    echo "  }" >> $output_file
  else
    echo "Failed to download artifact for $version ($short_commit)."
  fi
done

# Close the JSON array
echo "]" >> $output_file

# Output success message
echo "Successfully downloaded artifacts have been logged to $output_file."
