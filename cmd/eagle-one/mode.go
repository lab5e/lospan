package main

import (
	"github.com/lab5e/lospan/pkg/pb/lospan"
)

// E1Mode is the E1 modes
type E1Mode interface {
	Prepare(client lospan.LospanClient, app *lospan.Application, gw *lospan.Gateway) error
	Cleanup(client lospan.LospanClient, app *lospan.Application, gw *lospan.Gateway)
	Run(gatewayChannel chan string, publisher *EventRouter, app *lospan.Application, gw *lospan.Gateway)
	Failed() bool
}
