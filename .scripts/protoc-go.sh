#!/bin/bash

# This script installs protoc-gen-go for the given version.
# See: https://github.com/golang/protobuf/releases

VERSION=$1
if [ -z "$1" ] ; then
  VERSION=v1.2.0
fi

export GO111MODULE=off

echo "Locking protobuf to $VERSION..."
go get -u github.com/golang/protobuf/protoc-gen-go
git -C "$(go env GOPATH)"/src/github.com/golang/protobuf checkout "$VERSION"

echo "Installing to GOBIN..."
go install github.com/golang/protobuf/protoc-gen-go
echo "Done"
