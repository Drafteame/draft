#!/usr/bin/env bash

if [ "$GIT_NV" = "1" ]; then
  echo "[force-mod-version] skipping go.mod version check"
  exit 0
fi

echo "[force-mod-version] checking go.mod files"

set -e -o pipefail

version="1.22"
if [ $# -gt 0 ]; then
  version="$1"
fi


function sed_replace() {
  local os
  os=$(uname)

  if [ "$os" = "Darwin" ]; then
    sed -i '' "$1" "$2"
  else
    sed -i "$1" "$2"
  fi
}

function update_go_version() {
  local file="$1"

  # Check if file exists
  if [ ! -f "$file" ]; then
    echo "[force-mod-version] skipping $file (file not found)"
    return 0
  fi

  if grep -q "^toolchain go" "$file"; then
    echo "[force-mod-version] removing gotoolchain version from $file"
    sed_replace '/^toolchain go/d' "$file"
  fi

  if current_version=$(grep "^go [0-9]" "$file" | awk '{print $2}'); then
    if [ "$current_version" != "$version" ]; then
      echo "[force-mod-version] updating go version from $current_version to $version in $file"
      pattern="s|^go .*|go $version|g"
      sed_replace "$pattern" "$file"
    fi
  fi
}

# Update go.work if it exists
if [ -f "./go.work" ]; then
  update_go_version "./go.work"
fi

# Find and process go.mod files
go_mod_files=$(fd 'go.mod' --glob 2>/dev/null)
if [ -n "$go_mod_files" ]; then
  while IFS= read -r file; do
    update_go_version "$file"
  done <<< "$go_mod_files"
fi

git add --all
exit 0
