package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/lab5e/l5log/pkg/lg"
)

func main() {
	lg.InitLogs("lc", lg.LogParameters{
		Type:     "plain",
		Level:    "debug",
		LiveLogs: false,
	})
	var config params
	kctx := kong.Parse(&config, kong.Name("lc"))

	kong.Bind(&config)
	if err := kctx.Run(); err != nil {
		lg.Error(err.Error())
		os.Exit(2)
	}

}