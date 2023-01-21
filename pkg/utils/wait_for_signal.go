package utils

import (
	"os"
	"os/signal"

	"github.com/lab5e/lospan/pkg/lg"
)

// WaitForSignal waits for any signal (SIGHUP, SIGTERM...) before returning
func WaitForSignal() {
	terminator := make(chan bool)
	sigch := make(chan os.Signal, 2)
	signal.Notify(sigch, os.Interrupt)
	go func() {
		sig := <-sigch
		lg.Debug("Caught signal '%v'", sig)
		terminator <- true
	}()

	<-terminator
}
