#!/usr/bin/env bash

rm -rf build
mkdir -p build
for app in app/*; do
   go build -o build/ $app/$(basename $app).go
   go run $app/$(basename $app).go -h
   pushd $app
   go test -v
   popd
done
