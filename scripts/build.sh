#!/bin/bash -e

SRC_DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd ) # Directory above this script

source $SRC_DIR/env.sh

# Fetch dependencies (salticidae)
echo "Fetching dependencies..."
go mod download

if [ ! -d $SALTICIDAE_PATH ]; then
    echo "couldn't find salticidae version ${SALTICIDAE_VER} at ${SALTICIDAE_PATH}"
    echo "build failed"
    exit 1
fi

# Build salticidae
echo "Building salticidae..."
if [[ "$OSTYPE" == "linux-gnu" ]]; then
    chmod -R u+w $SALTICIDAE_PATH
    cd $SALTICIDAE_PATH
    cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX="$SALTICIDAE_PATH/build" .
    make -j4
    make install
    cd -
    mkdir -p $SALTICIDAE_GO_PATH/build
    ln -svf "$SALTICIDAE_PATH/build/lib/libsalticidae.a" "$SALTICIDAE_GO_PATH/build/libsalticidae.a"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    brew install Determinant/salticidae/salticidae
else
    echo "Your operating system is not supported."
    exit 1
fi
