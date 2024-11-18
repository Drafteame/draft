#!/usr/bin/env bash

if [ "$GIT_NV" = "1" ]; then
  echo "[go-mod-tidy] skipping go.mod files check"
  exit 0
fi

echo "[go-mod-tidy] checking go.mod files"

function get_base_dir() {
  local dir

  local file="$1"
  dir=$(dirname "$file")

  echo "$dir"
}

# Find and process go.mod files
go_mod_files=$(fd 'go.mod' --glob)
if [ -n "$go_mod_files" ]; then
  for file in $go_mod_files; do
    echo "[go-mod-tidy] tidying $file"

    cd "$(get_base_dir "$file")" && go mod tidy -v

    # shellcheck disable=SC2181
    if [ $? -ne 0 ]; then
      exit 2
    fi

    # shellcheck disable=SC2164
    cd - > /dev/null
  done
fi

git add --all
exit 0
