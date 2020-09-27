GOPATH:=$(shell go env GOPATH)

.PHONY: proto
proto:
	protoc -I. \
 		   -I/usr/local/include \
           -I${GOPATH}/src \
           -I${GOPATH}/src/github.com/googleapis/googleapis \
           -I${GOPATH}/src/github.com/pubgo/catdog/example/hello/proto \
           --proto_path=${GOPATH}/src:. \
           --micro_out=. \
           --go_out=${MODIFY}:. \
           proto/*

.PHONY: build
build: proto
	go build -o helloworld *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t helloworld:latest
