#!/bin/sh

DIR=$1
OUT=$2

if [ -z "$DIR"  ] || [ -z "$OUT" ]; then
  echo "Usage: $0 <binary dir> <absolute output filepath>"
  exit 1
fi

docker run --rm -t -w "$DIR" -v "$DIR:$DIR:ro" registry.salsa.debian.org/reproducible-builds/diffoscope \
    --no-progress --html - geth-reference geth-reproduce > "$OUT"