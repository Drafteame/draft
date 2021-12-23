package tpl

var (
	BuildSh = `#!/usr/bin/env bash

  BASE_PATH=./cmd
  OUT_PATH=./bin

  # shellcheck disable=SC2231
  for dir in $BASE_PATH/*/
  do
    dir=${dir%*/}
    handler=${dir##*/}

    export GO111MODULE=on
    env GOOS=linux go build -ldflags="-s -w" -o $OUT_PATH/"$handler" $BASE_PATH/"$handler"/main.go

    echo "Handler $handler compiled"
  done`
)
