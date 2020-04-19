#!/bin/bash -e
tmpdir=$(mktemp -d -t salticidae-go-XXXXXXXX)
version='v0.1.0'
cd "$tmpdir"
curl -s https://raw.githubusercontent.com/ava-labs/salticidae-go/"$version"/scripts/build.sh -o ./build.sh
curl -s https://raw.githubusercontent.com/ava-labs/salticidae-go/"$version"/scripts/env.sh -o ./env.sh
chmod +x ./build.sh
source ./env.sh
./build.sh
cd -
rm -rf "$tmpdir"
unset tmpdir
