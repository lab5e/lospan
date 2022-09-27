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
	cd cmd/congress && go build -o ../../bin/congress
	cd cmd/datagenerator && go build -o ../../bin/datagenerator
	cd cmd/eagle-one && go build -o ../../bin/eagle-one



