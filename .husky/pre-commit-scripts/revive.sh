#!/usr/bin/env bash

if [ "$GIT_NV" = "1" ]; then
  echo "[revive] skipping go files check"
  exit 0
fi

echo "[revive] checking go files"

revive -config=revive.toml -formatter=friendly ./...

# shellcheck disable=SC2181
if [ $? -ne 0 ]; then
  echo "[revive] found issues in go files"
  exit 1
fi

exit 0
