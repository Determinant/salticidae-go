#!/bin/bash -e

prefix="$(pwd)/build"
SRC_DIR="$(dirname "${BASH_SOURCE[0]}")"

source "${SRC_DIR}/env.sh"

if [[ "$OSTYPE" == "linux-gnu" ]]; then
    go get -d "github.com/$SALTICIDAE_ORG/salticidae-go"
    cd "$SALTICIDAE_PATH"
    cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX="$SALTICIDAE_PATH/build" .
    make -j4
    make install
    cd -
elif [[ "$OSTYPE" == "darwin"* ]]; then
    brew install Determinant/salticidae/salticidae
else
    echo "Your system is not supported yet."
    exit 1
fi

rm -f "$prefix/libsalticidae.a"
ln -sv "$GOPATH/src/github.com/$SALTICIDAE_ORG/salticidae-go/salticidae/libsalticidae.a" "$prefix/libsalticidae.a"
