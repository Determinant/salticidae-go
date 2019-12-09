.PHONY: all clean cdep examples

all: cdep examples

examples: build/test_msgnet build/test_p2p_stress build/test_msgnet_tls build/bench_network

cdep: build/libsalticidae.a

build/libsalticidae.a:
	scripts/build.sh
build/test_msgnet: build/libsalticidae.a test_msgnet/main.go
	source scripts/env.sh && go build -o $@ github.com/ava-labs/salticidae-go/test_msgnet
build/test_msgnet_tls: build/libsalticidae.a test_msgnet_tls/main.go
	source scripts/env.sh && go build -o $@ github.com/ava-labs/salticidae-go/test_msgnet_tls
build/test_p2p_stress: build/libsalticidae.a test_p2p_stress/main.go
	source scripts/env.sh && go build -o $@ github.com/ava-labs/salticidae-go/test_p2p_stress
build/bench_network: build/libsalticidae.a bench_network/main.go
	source scripts/env.sh && go build -o $@ github.com/ava-labs/salticidae-go/bench_network

clean:
	rm -rf build/
	scripts/clean.sh
