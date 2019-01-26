# nwHacks 2019

## Protobuf

Install [protobuf](https://github.com/protocolbuffers/protobuf/releases) v3.6+,
then:

```sh
# for generating Go stubs
bash .scripts/protoc-go.sh

# for generating Swift stubs
brew install swift-protobuf 
```

Running `make proto` will generate all the stubs.

## Frontend

TODO

## Backend

```
cd backend
export GO111MODULE=on
go mod vendor
```
