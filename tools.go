//go:build tools
// +build tools

package tools

import (
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/mgechev/revive"
	_ "golang.org/x/lint/golint"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
)
