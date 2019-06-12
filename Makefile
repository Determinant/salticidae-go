.PHONY: all clean

all: build/test_msgnet build/test_p2p_stress

salticidae/libsalticidae.so:
	cd salticidae/; cmake -DCMAKE_BUILD_TYPE=Release -DBUILD_SHARED=ON  -DSALTICIDAE_DEBUG_LOG=OFF -DSALTICIDAE_CBINDINGS=ON -DBUILD_TEST=OFF ./
	make -C salticidae/

build:
	mkdir -p build

build/test_msgnet: salticidae/libsalticidae.so test_msgnet/main.go
	go build -o $@ github.com/Determinant/salticidae-go/test_msgnet
build/test_p2p_stress: salticidae/libsalticidae.so test_p2p_stress/main.go
	go build -o $@ github.com/Determinant/salticidae-go/test_p2p_stress

clean:
	rm -r build/
	cd salticidae/; make clean
	rm salticidae/CMakeCache.txt
