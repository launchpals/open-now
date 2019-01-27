# Open Now

![nwHacks 2019](https://img.shields.io/badge/nwhacks-2019-06C1C0.svg) [![Deployed with Inertia](https://img.shields.io/badge/deploying%20with-inertia-blue.svg)](https://github.com/ubclaunchpad/inertia)

There are many cases where innocent people commuting home are attacked shortly after getting off the bus. Last week, one of our members felt threatened that they were going to be followed home after highly uncomfortable interactions with a stranger on the bus. In this state of panic, Google Maps was found to be lacking in design for users who are in fight or flight mode, when all they want to do is to find the closest safe place to get off the bus.

*Open Now* presents a simple, easy-to-parse interface that immediately presents you with a number of options for reaching the nearest safe haven, as recommended based on contextual data.

## Features

* Quick, at-a-glance overview of possible routes nearby safe havens and destinations
* Intelligent suggestions based on contextual data such as your current and predicted trajectory, public transit mode, walking pace, and weather
* Detailed, turn-by-turn directions and destination details just a tap away
* Dark theme optimized for night-time environments

## How we built it

*Open Now* is an iOS app built in Swift, backed by a server written in Golang that powers our intelligent point-of-interest recommendations. The app and the server communicates using protocol buffers to serialize data transfer over Google’s remote procedure call framework, gRPC.

The server communicates with open-source public transit databases as well as the Google Maps Platform to generate recommendations, and is hosted using [Inertia](https://github.com/ubclaunchpad/inertia) — a continuous deployment tool that we previously built — to handle automated updates on our cloud instance.

## Development

### Protobuf

Install [protobuf](https://github.com/protocolbuffers/protobuf/releases) v3.6+,
then:

```sh
# for generating Go stubs
bash .scripts/protoc-go.sh

# for generating Swift stubs
brew install swift-protobuf 
```

Running `make proto` will generate all the stubs.

### Frontend

TODO

## Backend

```
cd backend
export GO111MODULE=on
go mod vendor
```
