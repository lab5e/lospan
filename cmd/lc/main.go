package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/lab5e/l5log/pkg/lg"
)

func main() {
	var config params
	kctx := kong.Parse(&config, kong.Name("lc"))

	kong.Bind(&config)
	if err := kctx.Run(); err != nil {
		lg.Error("Error: %v", err)
		os.Exit(2)
	}

}
