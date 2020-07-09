TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
GOPATH?=$$(go env GOPATH)
GRPC_GATEWAY_PATH?=$$(go list -m -u -f '{{ .Dir }}' all | grep 'github.com/grpc-ecosystem/grpc-gateway@')

proto:
	protoc --proto_path=protos/:$(GRPC_GATEWAY_PATH)/third_party/googleapis/ \
		--go_out=plugins=grpc:. \
		--grpc-gateway_out=logtostderr=true:. \
		middleware.proto

.PHONY: proto