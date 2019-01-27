# Base build  wit hdependencies
FROM golang:1.11-alpine AS deps
RUN apk add git gcc g++ 
WORKDIR /project
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download
 
# This image builds the weavaite server
FROM deps AS build
COPY . .
RUN go install ./backend

# Copy into minimal final image
FROM alpine AS open_now
RUN apk add ca-certificates
COPY --from=build /go/bin/backend /bin/open_now_server
EXPOSE 8081
ENTRYPOINT ["/bin/open_now_server"]
