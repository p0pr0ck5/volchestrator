.PHONY: protoc
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		svc/volchestrator.proto \
		svc/volchestrator_admin.proto

.PHONY: build
build:
	go build

.PHONY: test
test:
	go test -race -v ./...

.PHONY: dev
dev: build
	./volchestrator server
