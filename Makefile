.PHONY: all clean

all: build/test_msgnet

salticidae/libsalticidae.so:
	cd salticidae/; cmake -DCMAKE_BUILD_TYPE=Release -DBUILD_SHARED=ON  -DSALTICIDAE_DEBUG_LOG=OFF -DSALTICIDAE_CBINDINGS=ON -DBUILD_TEST=OFF ./
	make -C salticidae/

build:
	mkdir -p build

build/test_msgnet: salticidae/libsalticidae.so
	go build -o $@ salticidae-go/test_msgnet
	go build -o $@ salticidae-go/test_p2p_stress
