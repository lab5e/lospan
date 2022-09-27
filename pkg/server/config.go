package server

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

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
	Syslog               bool
	DisableGatewayChecks bool
	UseSecureCookie      bool
	LogLevel             uint
	PlainLog             bool // Fancy stderr logs with emojis and colors
	MemoryDB             bool
	OnlyLoopback         bool // use only loopback adapter - for testing
	DebugPort            int  // Debug port - 0 for random, default 8081
	DBMaxConnections     int
	DBIdleConnections    int
	DBConnLifetime       time.Duration
}

// This is the default configuration
const (
	DefaultGatewayPort     = 8000
	DefaultHTTPPort        = 8080
	DefaultDebugPort       = 8081
	DefaultNetworkID       = 0
	DefaultMA              = "00-09-09"
	DefaultConnectHost     = "connect.staging.telenordigital.com"
	DefaultConnectClientID = "telenordigital-connectexample-web"
	DefaultLogLevel        = 0
	DefaultMaxConns        = 200
	DefaultIdleConns       = 100
	DefaultConnLifetime    = 10 * time.Minute
)

// NewDefaultConfig returns the default configuration. Note that this configuration
// isn't valid right out of the box; a storage backend must be selected.
func NewDefaultConfig() *Configuration {
	return &Configuration{
		MA:                DefaultMA,
		HTTPServerPort:    DefaultHTTPPort,
		NetworkID:         DefaultNetworkID,
		LogLevel:          DefaultLogLevel,
		DebugPort:         DefaultDebugPort,
		DBMaxConnections:  DefaultMaxConns,
		DBConnLifetime:    DefaultConnLifetime,
		DBIdleConnections: DefaultIdleConns,
	}
}

// NewMemoryNoAuthConfig returns a configuration with no authentication and
// memory-backed storage. This is a valid configuration.
func NewMemoryNoAuthConfig() *Configuration {
	ret := NewDefaultConfig()
	ret.MemoryDB = true
	return ret
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
