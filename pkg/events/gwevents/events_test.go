package gwevents

import "testing"

func TestEventCreation(t *testing.T) {
	NewInactive()
	NewKeepAlive()
	NewTx("some data")
	NewRx("some data")
}
