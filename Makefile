
.PHONY: setup
setup:
	bash .scripts/protoc-go.sh

.PHONY: proto
proto:
	protoc -I proto service.proto --go_out=plugins=grpc:proto/go
	protoc -I proto service.proto --swift_out=proto/swift --swiftgrpc_out=Client=true,Server=false:proto/swift
