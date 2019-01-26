
.PHONY: setup
setup:
	bash .scripts/protoc-gen-go.sh

.PHONY: proto
proto:
	protoc -I proto service.proto --go_out=plugins=grpc:proto/go
	protoc -I proto service.proto --swift_out=proto/ios
