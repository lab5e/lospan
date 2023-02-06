package server

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/lab5e/lospan/pkg/protocol"
)

// Parameters holds the configuration for the system
type Parameters struct {
	GRPCEndpoint         string `kong:"help='gRPC endpoint for API',default=':5150'"`
	GatewayPort          int    `kong:"help='Port for gateway interface',default='8000'"`
	NetworkID            uint   `kong:"help='Network ID for server',default='0'"`
	MA                   string `kong:"help='MA for key generator',default='00-00-00'"`
	ConnectionString     string `kong:"help='SQLite connection string',default=':memory:'"`
	DisableGatewayChecks bool   `kong:"help='Disable gateway IP address checking'"`
	DisableNonceCheck    bool   `kong:"help='Disable nonce check for devices',default='false'"`
}

// NewDefaultConfig returns the default configuration. Note that this configuration
// isn't valid right out of the box; a storage backend must be selected.
func NewDefaultConfig() *Parameters {
	return &Parameters{
		MA:               "00-00-00",
		NetworkID:        0,
		ConnectionString: ":memory:",
		GatewayPort:      8000,
	}
}

// RootMA returns the MA to use as the base MA for EUIs. The configuration
// is assumed to be valid at this point. If there's an error converting the
// MA it will panic.
func (cfg *Parameters) RootMA() protocol.MA {
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
func (cfg *Parameters) Validate() error {
	prefix, err := hex.DecodeString(strings.Replace(cfg.MA, "-", "", -1))
	if err != nil {
		return fmt.Errorf("invalid format for MA string: %v", err)
	}
	_, err = protocol.NewMA(prefix)
	if err != nil {
		return fmt.Errorf("unable to create MA: %v", err)
	}

	if cfg.ConnectionString == "" {
		return errors.New("connection string is blank")
	}

	return nil
}
