# NOTE: do NOT use the Makefile directly, please run `./build.sh` instead

.PHONY: all clean

all: build/test_msgnet build/test_p2p_stress build/test_msgnet_tls build/bench_network

build:
	mkdir -p build

build/test_msgnet: build/lib/libsalticidae.a test_msgnet/main.go
	go build -o $@ github.com/ava-labs/salticidae-go/test_msgnet
build/test_msgnet_tls: build/lib/libsalticidae.a test_msgnet_tls/main.go
	go build -o $@ github.com/ava-labs/salticidae-go/test_msgnet_tls
build/test_p2p_stress: build/lib/libsalticidae.a test_p2p_stress/main.go
	go build -o $@ github.com/ava-labs/salticidae-go/test_p2p_stress
build/bench_network: build/lib/libsalticidae.a bench_network/main.go
	go build -o $@ github.com/ava-labs/salticidae-go/bench_network

clean:
	rm -rf build/
	cd salticidae/; make clean
	rm salticidae/CMakeCache.txt
