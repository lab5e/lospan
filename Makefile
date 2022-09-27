.PHONY: all
.PHONY: vet
.PHONY: text
.PHONY: build 

all: vet  build

vet:
	go vet ./...

test:
	go test -timeout 10s ./...

build:
	go build -o bin/congress
	cd test-tools/datagenerator && go build -o ../../bin/datagenerator
	cd test-tools/eagle-one && go build -o ../../bin/eagle-one



