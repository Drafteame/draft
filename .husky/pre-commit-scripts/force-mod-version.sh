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

update_go_version "./go.work"

# Find and process go.mod files
go_mod_files=$(fd 'go.mod' --glob)
if [ -n "$go_mod_files" ]; then
  for file in $go_mod_files; do
    update_go_version "$file"
  done
fi

git add --all
exit 0
