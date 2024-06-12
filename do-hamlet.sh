#!/usr/bin/env bash

rm -f *.txt*
if [ ! -d build ]; then
   echo "No build outputs"
   exit 1
fi

if [ -f "build/sliding_wordcount" ]; then
   echo "Sequential"
   cat data/hamlet-10-copies.txt | time -o no-go-routines.txt ./build/sliding_wordcount > /dev/null

   echo "Functional Style / Channels / Varying buffer size"
   for BUFFER_SIZE in 1 10 100 1000 10000 100000 1000000 10000000 ; do
      echo "buffer size", $BUFFER_SIZE
      cat data/hamlet-10-copies.txt | time -o go-routines-$BUFFER_SIZE.txt ./build/sliding_wordcount -go-routines -channel-size $BUFFER_SIZE > /dev/null
   done
   echo ":"
   cat no-go-routines.txt
   for BUFFER_SIZE in 1 10 100 1000 10000 100000 1000000 10000000 ; do
      echo $BUFFER_SIZE":"; cat go-routines-$BUFFER_SIZE.txt
   done
else
   echo "No executable"
   exit 2
fi
