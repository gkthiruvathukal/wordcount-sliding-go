#!/usr/bin/env bash

for lib in libs/*; do
  pushd $lib
  go test -v
  popd
done
