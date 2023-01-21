package server

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/lab5e/lospan/pkg/protocol"
)

var defaultMA protocol.MA

func init() {
	tmp := strings.Replace("00-09-09", "-", "", -1)
	tmp2, _ := hex.DecodeString(tmp)
	defaultMA, _ = protocol.NewMA(tmp2)
}

// Test custom parsing for all parameters. These cannot be tested multiple
// times or with custom parameters but I'm going to assume it works.
func TestCommandLineConfigDefaults(t *testing.T) {
	config := NewDefaultConfig()
	if err := config.Validate(); err != nil {
		t.Fatalf("Expected config to be valid: %v", err)
	}
}

func TestValidConfiguration(t *testing.T) {
	config := NewDefaultConfig()
	config.MA = "foof"
	if err := config.Validate(); err == nil {
		t.Fatal("Expected error from invalid MA string")
	}
	config.MA = "00-09-09-09-09-09-09"
	if err := config.Validate(); err == nil {
		t.Fatal("Expected error from too long MA string")
	}

	config.MA = "00-00-00"

	config.ConnectionString = ""
	if err := config.Validate(); err == nil {
		t.Fatalf("Expected error with no backend selected")
	}

}

func TestMAInvalidString(t *testing.T) {
	config := NewDefaultConfig()
	config.MA = "00-00-00"
	config.RootMA()

	config.MA = "foof"
	defer func() {
		if n := recover(); n == nil {
			t.Fatal("Should have gotten a panic")
		}
	}()

	config.RootMA() // should panic
	t.Fatal("I expected panic here")

}

func TestMAInvalidMA(t *testing.T) {
	config := NewDefaultConfig()
	config.MA = "01-02-03-04-05-06-07-08-09"
	defer func() {
		if n := recover(); n == nil {
			t.Fatal("Should have gotten a panic")
		}
	}()

	config.RootMA() // should panic
	t.Fatal("I expected panic here")
}
