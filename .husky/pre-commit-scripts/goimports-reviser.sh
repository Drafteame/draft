#!/usr/bin/env bash

if [ "$GIT_NV" = "1" ]; then
  echo "[goimports-reviser] skipping goimports formatting"
  exit 0
fi

echo "[goimports-reviser] formatting code"

goimports-reviser -format ./...

# shellcheck disable=SC2181
if [ $? -ne 0 ]; then
  echo "[goimports-reviser] found issues formatting go files"
  exit 1
fi

git add --all
exit 0
