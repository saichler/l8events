#!/usr/bin/env bash

set -e

# Download api.proto dependency (for L8MetaData in List types)
wget https://raw.githubusercontent.com/saichler/l8types/refs/heads/main/proto/api.proto

# Generate bindings for l8events proto files
docker run --user "$(id -u):$(id -g)" -e PROTO="l8events.proto l8events_categories.proto" --mount type=bind,source="$PWD",target=/home/proto/ -i saichler/protoc:latest

# Clean up downloaded proto
rm api.proto

# Move generated bindings to the types directory and clean up
rm -rf ../go/types
mkdir -p ../go/types
mv ./types/* ../go/types/.
rm -rf ./types

# Fix import paths for l8api types
cd ../go
find . -name "*.go" -type f -exec sed -i 's|"./types/l8api"|"github.com/saichler/l8types/go/types/l8api"|g' {} +

# Clean up
cd ../proto
rm -rf *.rs
