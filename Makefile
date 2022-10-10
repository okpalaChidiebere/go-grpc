PROTO_DIR := ./protos

build:
	protoc --proto_path=$(PROTO_DIR) --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    protos/*.proto

# NOTE: you you have the /pb folder already created
clean:
	rm -rf pb/*.pb.*

.PHONY: clean build