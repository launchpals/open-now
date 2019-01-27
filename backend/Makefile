
all: check

.PHONY: check
check:
	go vet ./...

.PHONY: deps
deps:
	( cd ../ ; env GO111MODULE=on go mod vendor ; env GO111MODULE=on go mod tidy )

.PHONY: run
run:
	go run main.go

.PHONY: proto
proto:
	(cd ../ ; make proto)
