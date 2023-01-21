package frequency

// Stub functions for frequency management.

import (
	"github.com/lab5e/lospan/pkg/protocol"
)

// TODO: Replace with proper frequency management. These values are the defaults
// for the EU868 band and won't work for other bands.

// GetDLSettingsOTAA returns the DLSettings value returned during OTAA
// join procedure. This returns the default values for now.
func GetDLSettingsOTAA() protocol.DLSettings {
	return protocol.DLSettings{
		RX1DRoffset: 0,
		RX2DataRate: 5,
	}
}

// GetRxDelayOTAA returns the RxDelay parameter for OTAA join procedures. This
// function always returns 1, the default value.
func GetRxDelayOTAA() uint8 {
	return 1
}

// GetCFListOTAA returns the CFList type used during OTAA. It always returns
// the default values.
func GetCFListOTAA() protocol.CFList {
	return protocol.CFList{}
}
