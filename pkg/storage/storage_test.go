package storage

import (
	"crypto/rand"

	"github.com/lab5e/lospan/pkg/protocol"
)

func makeRandomEUI() protocol.EUI {
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	ret := protocol.EUI{}
	copy(ret.Octets[:], randomBytes)
	return ret
}

func makeRandomData() []byte {
	randomBytes := make([]byte, 30)
	rand.Read(randomBytes)
	return randomBytes
}

func makeRandomKey() protocol.AESKey {
	var keyBytes [16]byte
	copy(keyBytes[:], makeRandomData())
	return protocol.AESKey{Key: keyBytes}
}
