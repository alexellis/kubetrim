#!/bin/sh

for f in bin/kubetrim*; do shasum -a 256 $f > $f.sha256; done
