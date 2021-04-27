.PHONY: protoc
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		svc/volchestrator.proto \
		svc/volchestrator_admin.proto

build:
	go build

test:
	go test -race -v ./...

dev: build
	./volchestrator server
