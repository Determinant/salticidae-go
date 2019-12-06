export GOPATH="$(go env GOPATH)"
export SALTICIDAE_ORG="Determinant"

if [[ "$OSTYPE" == "linux-gnu" ]]; then
    export SALTICIDAE_PATH="$GOPATH/src/github.com/$SALTICIDAE_ORG/salticidae-go/salticidae"
    export CGO_CFLAGS="-I$SALTICIDAE_PATH/build/include/"
    export CGO_LDFLAGS="-L$SALTICIDAE_PATH/build/lib/ -lsalticidae -luv -lssl -lcrypto -lstdc++"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    export CGO_CFLAGS="-I/usr/local/opt/openssl/include"
    export CGO_LDFLAGS="-L/usr/local/opt/openssl/lib/ -lsalticidae -luv -lssl -lcrypto"
else
    echo "Your system is not supported yet."
    exit 1
fi
