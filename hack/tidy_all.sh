#! /bin/bash
cd "$(git rev-parse --show-toplevel)"
find . -name go.mod |xargs dirname |while read dir ; do pushd $dir ; go mod tidy; popd; done
