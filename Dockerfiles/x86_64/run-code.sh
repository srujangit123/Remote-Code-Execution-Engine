#!/bin/bash

LANGUAGE=$1
FILE=$2
INPUT=$3

case $LANGUAGE in
  "cpp")
    g++ -o /tmp/code $FILE && /tmp/code < $INPUT
    ;;
  "golang")
    go run $FILE < $INPUT
    ;;
  *)
    echo "Unsupported language"
    exit 1
    ;;
esac