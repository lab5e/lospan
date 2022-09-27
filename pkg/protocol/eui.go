package protocol

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// EUI represents an IEEE EUI-64 identifier. The identifier is described at
// http://standards.ieee.org/develop/regauth/tut/eui64.pdf
type EUI struct {
	Octets [8]byte
}

// String returns a string representation of the EUI (XX-XX-XX-XX...)
func (eui EUI) String() string {
	return fmt.Sprintf("%02x-%02x-%02x-%02x-%02x-%02x-%02x-%02x",
		eui.Octets[0], eui.Octets[1], eui.Octets[2], eui.Octets[3],
		eui.Octets[4], eui.Octets[5], eui.Octets[6], eui.Octets[7])
}

// EUIFromString converts a string on the format "xx-xx-xx..." in hex to an
// internal representation
func EUIFromString(euiStr string) (EUI, error) {
	tmpBuf, err := hex.DecodeString(strings.TrimSpace(strings.Replace(euiStr, "-", "", -1)))
	if err != nil {
		return EUI{}, err
	}
	if len(tmpBuf) != 8 {
		return EUI{}, ErrInvalidParameterFormat
	}
	ret := EUI{}
	copy(ret.Octets[:], tmpBuf)
	return ret, nil
}

// EUIFromUint64 converts an uint64 value to an EUI.
func EUIFromInt64(val int64) EUI {
	ret := EUI{}
	for i := 7; i >= 0; i-- {
		ret.Octets[i] = byte(val & 0xFF)
		val >>= 8
	}
	return ret
}

// ToUint64 returns the EUI as a uin64 integer
func (eui *EUI) ToInt64() int64 {
	ret := int64(0)
	for i := 0; i < 8; i++ {
		ret <<= 8
		ret += int64(eui.Octets[i])
	}
	return ret
}
