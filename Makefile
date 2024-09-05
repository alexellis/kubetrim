Version := $(shell git describe --tags --dirty)
GitCommit := $(shell git rev-parse HEAD)
LDFLAGS := "-s -w -X github.com/alexellis/kubetrim/pkg.Version=$(Version) -X github.com/alexellis/kubetrim/pkg.GitCommit=$(GitCommit)"
SOURCE_DIRS = cmd pkg main.go
export GO111MODULE=on

.PHONY: all
all: gofmt test build dist compress hash

.PHONY: build
build:
	go build

.PHONY: gofmt
gofmt:
	@test -z $(shell gofmt -l -s $(SOURCE_DIRS) ./ |grep -v vendor/| tee /dev/stderr) || (echo "[WARN] Fix formatting issues with 'make gofmt'" && exit 1)

.PHONY: test
test:
	CGO_ENABLED=0 go test $(shell go list ./... | grep -v /vendor/|xargs echo) -cover

.PHONY: dist
dist:

	mkdir -p bin/
	mkdir -p uploads/
	rm -rf bin/kubetrim*
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -o bin/kubetrim
	CGO_ENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS) -o bin/kubetrim-darwin
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -a -ldflags $(LDFLAGS) -o bin/kubetrim-darwin-arm64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags $(LDFLAGS) -o bin/kubetrim-arm64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags $(LDFLAGS) -o bin/kubetrim.exe

.PHONY: hash
hash:
	rm -rf uploads/*.sha256 && ./hack/hashgen.sh

.PHONY: compress
compress:
	./hack/compress.sh