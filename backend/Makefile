
all: check

.PHONY: check
check:
	go vet ./...

.PHONY: deps
deps:
	( cd ../ ; env GO111MODULE=on go mod vendor)

.PHONY: run
run:
	go run main.go
