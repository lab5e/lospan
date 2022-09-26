.PHONY: all
.PHONY: vet
.PHONY: text
.PHONY: build 

all: vet test build

vet:
	go vet ./...

test:
	go test -timeout 10s ./...

build:
	go build 

