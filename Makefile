.PHONY: protoc
proto:
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		svc/*.proto

.PHONY: test
test: clean
	go test -race -v ./... && echo '\e[1;32mAll good!\e[0m' || echo '\e[1;31mNope!\e[0m'

.PHONE: clean
clean:
	go clean -cache