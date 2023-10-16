#!/bin/bash

for GOOS in linux darwin windows; do
  for GOARCH in arm64 amd64; do
    EXT=""
    if [[ "$GOOS" == "windows" ]]; then
      EXT=".exe"
    fi
    echo "Building binary for Architecture=${GOARCH} and OS=${GOOS} ..."
    GOARCH=$GOARCH GOOS=$GOOS go build -o "build/spring-boot-scanner-${GOOS}-${GOARCH}${EXT}"
  done
done
