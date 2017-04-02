#!/bin/bash

set -ve

for d in $(find . -mindepth 1 -type d); do
    pushd $d &> /dev/null

    protoc -I/usr/local/include -I. \
        -I$GOPATH/src \
        -I${GOPATH}/src/github.com/rnd/kudu/protos \
        --go_out=plugins=grpc:$GOPATH/src \
        *.proto

    popd &> /dev/null
done