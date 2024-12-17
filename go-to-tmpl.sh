#!/bin/bash

# Check if the user provided a directory as an argument
if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <directory>"
  exit 1
fi

# Get the directory from the first argument
directory=$1

# Check if the directory exists
if [ ! -d "$directory" ]; then
  echo "Error: Directory $directory does not exist."
  exit 1
fi

# Recursively find all .go files in the directory and rename them to .tmpl
fd ".*\.go$" "$directory" --type f | while read -r file; do
  new_name="${file}.tmpl"
  mv "$file" "$new_name"
  echo "Renamed: $file -> $new_name"
done

echo "Conversion complete."