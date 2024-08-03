#!/bin/sh


REF=$1
REP=$2

readelf -p .rodata "$REF" | grep /home/travis > ref-path-1.txt
readelf -p .rodata "$REP" | grep /root/go/pkg > ref-path-2.txt

colordiff ref-path-1.txt ref-path-2.txt

