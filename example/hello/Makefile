GOPATH:=$(shell go env GOPATH)

.PHONY: proto clear
proto: clear
	protoc -I. \
 		   -I/usr/local/include \
           -I${GOPATH}/src \
           -I${GOPATH}/src/github.com/googleapis/googleapis \
           -I${GOPATH}/src/github.com/pubgo/catdog/example/hello/proto \
           --proto_path=${GOPATH}/src:. \
           --catdog_out=. \
           --go_out=. \
           proto/*

.PHONY: run build
run: proto
	go run cmd/main.go

build:
	go build -o main cmd/main.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t helloworld:latest

clear:
	rm -rf proto/*.go
	rm -rf proto/**/*.go