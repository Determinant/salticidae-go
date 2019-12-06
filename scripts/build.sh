#!/bin/bash -e

prefix="${1:-$(pwd)}/build"
SRC_DIR="$(dirname "${BASH_SOURCE[0]}")"

source "${SRC_DIR}/env.sh"

if [[ "$OSTYPE" == "linux-gnu" ]]; then
    go get -u -d "github.com/$SALTICIDAE_ORG/salticidae-go"
    cd "$SALTICIDAE_GO_PATH"
    git submodule update --init --recursive
    cd "$SALTICIDAE_PATH"
    cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX="$SALTICIDAE_PATH/build" .
    make -j4
    make install
    cd -
    rm -f "$prefix/libsalticidae.a"
    ln -sv "$SALTICIDAE_PATH/libsalticidae.a" "$prefix/libsalticidae.a"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    brew install Determinant/salticidae/salticidae
else
    echo "Your system is not supported yet."
    exit 1
fi
