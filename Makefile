
src_files = cmd/comfy/sensibo_proxy.go cmd/comfy/tibber_proxy.go  cmd/comfy/price_cache.go
test_files = cmd/comfy/price_cache_test.go

build:
	go build -o bin/comfy cmd/comfy/main.go ${src_files}

run:
	go run cmd/comfy/main.go ${src_files}

test:
	go test ${src_files} ${test_files}
