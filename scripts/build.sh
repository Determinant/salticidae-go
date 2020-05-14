#!/bin/bash -e

PREFIX="${PREFIX:-$(pwd)/build}"
SRC_DIR="$(dirname "${BASH_SOURCE[0]}")"

source "${SRC_DIR}/env.sh"

if [[ "$OSTYPE" == "linux-gnu" ]]; then
    go get -u -d "github.com/$SALTICIDAE_ORG/salticidae-go"
    cd "$SALTICIDAE_GO_PATH"
    git -c advice.detachedHead=false checkout "$SALTICIDAE_GO_VER"
    git submodule update --init --recursive
    cd "$SALTICIDAE_PATH"
    cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX="$SALTICIDAE_PATH/build" .
    make -j4
    make install
    cd -
    mkdir -p "$PREFIX"
    rm -f "$PREFIX/libsalticidae.a"
    ln -sv "$SALTICIDAE_PATH/build/lib/libsalticidae.a" "$PREFIX/libsalticidae.a"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    brew install Determinant/salticidae/salticidae
else
    echo "Your system is not supported yet."
    exit 1
fi
