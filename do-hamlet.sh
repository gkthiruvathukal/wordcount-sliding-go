#!/usr/bin/env bash

if [ ! -d build ]; then
   echo "No build outputs"
   exit 1
fi

if [ -f "build/sliding_wordcount" ]; then
   echo "Sequential"
   time cat data/hamlet-10-copies.txt | ./build/sliding_wordcount > /dev/null
   echo "Functional Style / Channels / Varying buffer size"
   for BUFFER_SIZE in 10 100 1000 10000 100000 1000000; do
      echo "buffer size", $BUFFER_SIZE
      time cat data/hamlet-10-copies.txt | ./build/sliding_wordcount -go-routines -channel-size $BUFFER_SIZE > /dev/null
   done
else
   echo "No executable"
   exit 2
fi
