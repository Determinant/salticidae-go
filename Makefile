.PHONY: all clean

all: build/test_msgnet build/test_p2p_stress build/test_msgnet_tls build/bench_network


salticidae/libsalticidae.so:
	cd salticidae/; cmake -DCMAKE_BUILD_TYPE=Release -DBUILD_SHARED=ON  -DSALTICIDAE_DEBUG_LOG=OFF -DSALTICIDAE_CBINDINGS=ON -DBUILD_TEST=OFF ./
	make -C salticidae/ -j4

build:
	mkdir -p build

build/test_msgnet: salticidae/libsalticidae.so test_msgnet/main.go
	make -C salticidae/
	go build -o $@ github.com/Determinant/salticidae-go/test_msgnet
build/test_msgnet_tls: salticidae/libsalticidae.so test_msgnet_tls/main.go
	make -C salticidae/
	go build -o $@ github.com/Determinant/salticidae-go/test_msgnet_tls
build/test_p2p_stress: salticidae/libsalticidae.so test_p2p_stress/main.go
	make -C salticidae/
	go build -o $@ github.com/Determinant/salticidae-go/test_p2p_stress
build/bench_network: salticidae/libsalticidae.so bench_network/main.go
	make -C salticidae/
	go build -o $@ github.com/Determinant/salticidae-go/bench_network

clean:
	rm -rf build/
	cd salticidae/; make clean
	rm salticidae/CMakeCache.txt
