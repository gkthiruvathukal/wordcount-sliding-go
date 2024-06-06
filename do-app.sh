#!/usr/bin/env bash

for app in app/*; do
   go run $app/$(basename $app).go -h
   pushd $app
   go test -v
   popd
done
