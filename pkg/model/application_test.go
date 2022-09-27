package model

import "testing"

// This is borderline gaming the system but the code should at least run through once
func TestAppNonceGenerator(t *testing.T) {
	app := NewApplication()
	app.GenerateAppNonce()
}
