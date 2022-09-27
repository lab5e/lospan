.PHONY: cmd all vet text build  tools 

all: vet  build

vet:
	go vet ./...

test:
	go test -timeout 10s ./...

build: cmd

cmd:
	cd cmd/congress && go build -o ../../bin/congress
	cd cmd/datagenerator && go build -o ../../bin/datagenerator
	cd cmd/eagle-one && go build -o ../../bin/eagle-one

tools: 
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go get %

generate:
	buf mod update
	buf generate --path protobuf/lospan