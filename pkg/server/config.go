package server

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/lab5e/l5log/pkg/lg"
	"github.com/lab5e/lospan/pkg/protocol"
)

// Configuration holds the configuration for the system
type Configuration struct {
	GatewayPort          int
	HTTPServerPort       int
	NetworkID            uint   // The network ID that this instance handles. The default is 0
	MA                   string // String representation of MA
	DBConnectionString   string
	PrintSchema          bool
	DisableGatewayChecks bool
	UseSecureCookie      bool
	MemoryDB             bool
	OnlyLoopback         bool // use only loopback adapter - for testing
	Log                  lg.LogParameters
}

// This is the default configuration
const (
	DefaultGatewayPort = 8000
	DefaultHTTPPort    = 8080
	DefaultNetworkID   = 0
	DefaultMA          = "00-00-00"
)

// NewDefaultConfig returns the default configuration. Note that this configuration
// isn't valid right out of the box; a storage backend must be selected.
func NewDefaultConfig() *Configuration {
	return &Configuration{
		MA:             DefaultMA,
		HTTPServerPort: DefaultHTTPPort,
		NetworkID:      DefaultNetworkID,
	}
}

// RootMA returns the MA to use as the base MA for EUIs. The configuration
// is assumed to be valid at this point. If there's an error converting the
// MA it will panic.
func (cfg *Configuration) RootMA() protocol.MA {
	prefix, err := hex.DecodeString(strings.Replace(cfg.MA, "-", "", -1))
	if err != nil {
		panic("invalid format for MA string in configuration")
	}
	ret, err := protocol.NewMA(prefix)
	if err != nil {
		panic("unable to create MA")
	}
	return ret
}

// Validate checks the configuration for inconsistencies and errors. This
// function logs the warnings using the logger package as well.
func (cfg *Configuration) Validate() error {
	prefix, err := hex.DecodeString(strings.Replace(cfg.MA, "-", "", -1))
	if err != nil {
		return fmt.Errorf("invalid format for MA string: %v", err)
	}
	_, err = protocol.NewMA(prefix)
	if err != nil {
		return fmt.Errorf("unable to create MA: %v", err)
	}

	if cfg.DBConnectionString == "" && !cfg.MemoryDB {
		return errors.New("no backend storage selected. A connection string, embedded PostgreSQL or in-memory database must be selected")
	}

	return nil
}
