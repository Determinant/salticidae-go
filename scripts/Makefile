# NOTE: do NOT use the Makefile directly, please run `./build.sh` instead

.PHONY: all clean

all: build/test_msgnet build/test_p2p_stress build/test_msgnet_tls build/bench_network

build/libsalticidae.a: build
	rm -f $@
	ln -sv "${GOPATH}/src/github.com/Determinant/salticidae-go/salticidae/libsalticidae.a" $@

build:
	mkdir -p build

build/test_msgnet: build/libsalticidae.a test_msgnet/main.go
	go build -o $@ github.com/Determinant/salticidae-go/test_msgnet
build/test_msgnet_tls: build/libsalticidae.a test_msgnet_tls/main.go
	go build -o $@ github.com/Determinant/salticidae-go/test_msgnet_tls
build/test_p2p_stress: build/libsalticidae.a test_p2p_stress/main.go
	go build -o $@ github.com/Determinant/salticidae-go/test_p2p_stress
build/bench_network: build/libsalticidae.a bench_network/main.go
	go build -o $@ github.com/Determinant/salticidae-go/bench_network

clean:
	rm -rf build/
	cd salticidae/; make clean
	rm salticidae/CMakeCache.txt
