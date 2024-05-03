.PHONY: cmd all vet text build  tools priv

ifeq ($(GOPRIVATE),)
GOPRIVATE := github.com/lab5e/l5log
endif
all: vet priv build

priv:
	go env -w GOPRIVATE=$(GOPRIVATE)

vet:
	go vet ./...
	revive ./...

test:
	go test -timeout 10s ./...

build: cmd

cmd:
	cd cmd/lc && go build -o ../../bin/lc
	cd cmd/congress && go build -o ../../bin/congress
	cd cmd/datagenerator && go build -o ../../bin/datagenerator
	cd cmd/eagle-one && go build -o ../../bin/eagle-one

rpi: 
	cd cmd/lc && GOOS=linux GOARCH=arm go build -o ../../bin/lc.rpi
	cd cmd/congress && GOOS=linux GOARCH=arm go build -o ../../bin/congress.rpi

tools: 
	cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go get %

generate:
	buf mod update
	buf generate --path protobuf/lospan
