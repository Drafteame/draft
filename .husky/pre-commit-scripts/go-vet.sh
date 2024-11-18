#!/usr/bin/env bash

if [ "$GIT_NV" = "1" ]; then
  echo "[go-vet] skipping go files check"
  exit 0
fi

echo "[go-vet] checking go files"

go vet ./...

# shellcheck disable=SC2181
if [ $? -ne 0 ]; then
  echo "[go-vet] found issues in go files"
  exit 1
fi

exit 0
