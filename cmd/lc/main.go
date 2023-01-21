package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
)

func main() {

	var config params
	kctx := kong.Parse(&config, kong.Name("lc"))

	kong.Bind(&config)
	if err := kctx.Run(); err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		os.Exit(2)
	}

}
