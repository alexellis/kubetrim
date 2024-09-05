#!/bin/sh

cd bin

for f in kubetrim*; do tar -cvzf ../uploads/$f.tgz $f; done