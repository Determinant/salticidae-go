#!/bin/bash -xe
tmpdir=$(mktemp -d -t salticidae-go-XXXXXXXX)
cd "$tmpdir"
curl -s https://raw.githubusercontent.com/Determinant/salticidae-go/master/scripts/build.sh -o ./build.sh
curl -s https://raw.githubusercontent.com/Determinant/salticidae-go/master/scripts/env.sh -o ./env.sh
chmod +x ./build.sh
source ./env.sh
./build.sh "$SALTICIDAE_GO_PATH"
rm -rf "$tmpdir"
