#!/bin/bash

go mod download

APP="spring-boot-scanner"
for GOOS in linux darwin windows; do
  for GOARCH in arm64 amd64; do
    EXT=""
    if [[ "$GOOS" == "windows" ]]; then
      EXT=".exe"
    fi
    GOARCH=$GOARCH GOOS=$GOOS go build -o "build/${APP}-${GOOS}-${GOARCH}${EXT}"
  done
done
