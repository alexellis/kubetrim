#!/bin/sh

cd bin

for f in kubetrim*; do shasum -a 256 $f > ../uploads/$f.sha256; done

