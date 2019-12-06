#!/bin/bash
SRC_DIR="$(dirname "${BASH_SOURCE[0]}")"

source "$SRC_DIR/env.sh"
cd "$SALTICIDAE_PATH"
make clean
