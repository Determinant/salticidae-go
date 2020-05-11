#!/bin/bash -e

export SALTICIDAE_VER="v0.3.1"
export SALTICIDAE_PATH=$GOPATH/pkg/mod/github.com/ava-labs/salticidae@$SALTICIDAE_VER
export SALTICIDAE_GO_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )"; cd .. && pwd ) # Directory above this script

if [[ "$OSTYPE" == "linux-gnu" ]]; then
    export CGO_CFLAGS="-I$SALTICIDAE_PATH/build/include/"
    export CGO_LDFLAGS="-L$SALTICIDAE_PATH/build/lib/ -lsalticidae -luv -lssl -lcrypto -lstdc++ -g"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    export CGO_CFLAGS="-I/usr/local/opt/openssl/include"
    export CGO_LDFLAGS="-L/usr/local/opt/openssl/lib/ -lsalticidae -luv -lssl -lcrypto"
else
    echo "Your operating system is not supported"
    exit 1
fi