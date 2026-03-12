#!/usr/bin/env bash

set -e

# Generate bindings for l8events proto
docker run --user "$(id -u):$(id -g)" -e PROTO="l8events.proto" --mount type=bind,source="$PWD",target=/home/proto/ -i saichler/protoc:latest

# Move generated bindings to the types directory and clean up
rm -rf ../go/types
mkdir -p ../go/types
mv ./types/* ../go/types/.
rm -rf ./types

# Clean up
rm -rf *.rs
