#!/usr/bin/env bash

if [ ! -d build ]; then
   echo "No build outputs"
   exit 1
fi

if [ -f "build/sliding_wordcount" ]; then
   cat data/hamlet-10-copies.txt | ./build/sliding_wordcount
else
   echo "No executable"
   exit 2
fi
