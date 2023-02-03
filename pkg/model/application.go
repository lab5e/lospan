package model

import (
	"crypto/rand"

	"github.com/lab5e/lospan/pkg/protocol"
)

// Application represents a LoRa application instance.
type Application struct {
	AppEUI protocol.EUI // Application EUI
	Tag    string       // Tag data (for external ref)
}

// Equals returns true if the other application has identical fields. Just like
// ...Equals
func (a *Application) Equals(other Application) bool {
	return a.AppEUI == other.AppEUI
}

// NewApplication creates a new application instance
func NewApplication() Application {
	return Application{}
}

// GenerateAppNonce generates a new AppNonce, three random bytes that will be used
// to generate new devices.
func (a *Application) GenerateAppNonce() ([3]byte, error) {
	var nonce [3]byte
	_, err := rand.Read(nonce[:])
	return nonce, err
}
