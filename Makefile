Project=github.com/pubgo/catdog
GOPath=$(shell go env GOPATH)
Version=$(shell git tag --sort=committerdate | tail -n 1)
GoROOT=$(shell go env GOROOT)
BuildTime=$(shell date "+%F %T")
CommitID=$(shell git rev-parse HEAD)
LDFLAGS=-ldflags " \
-X 'github.com/pubgo/catdog/version.GoROOT=${GoROOT}' \
-X 'github.com/pubgo/catdog/version.BuildTime=${BuildTime}' \
-X 'github.com/pubgo/catdog/version.GoPath=${GOPath}' \
-X 'github.com/pubgo/catdog/version.CommitID=${CommitID}' \
-X 'github.com/pubgo/catdog/version.Project=${Project}' \
-X 'github.com/pubgo/catdog/version.Version=${Version:-v0.0.1}' \
"

.PHONY: build
build:
	@go build ${LDFLAGS} -mod vendor -race -v -o main main.go

build_hello_test:
	go build ${LDFLAGS} -mod vendor -v -o main  example/hello/main.go

.PHONY: install
install:
	@go install ${LDFLAGS} .

.PHONY: release
release:
	@go build ${LDFLAGS} -race -v -o main main.go

.PHONY: test
test:
	@go test -race -v ./... -cover

.PHONY: tag_list
tag_list:
	@git tag -n --sort=committerdate | tee | tail -n 5

.PHONY: tag
tag:tag_list
	tg=$(shell read -p "Enter Tag Version: :" name;echo $$name)
	@git tag ${tg}
	@git push origin ${tg}
	@git tag -n --sort=committerdate | tee | tail -n 5

