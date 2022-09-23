
build:
	go build -o bin/comfy cmd/comfy/sensibo_proxy.go cmd/comfy/tibber_proxy.go cmd/comfy/main.go

run:
	go run cmd/comfy/sensibo_proxy.go cmd/comfy/tibber_proxy.go cmd/comfy/main.go
