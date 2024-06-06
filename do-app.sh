#!/usr/bin/env bash

for app in app/*; do
   go build -o build $app/$(basename $app).go
   go run $app/$(basename $app).go -h
   pushd $app
   go test -v
   popd
done
