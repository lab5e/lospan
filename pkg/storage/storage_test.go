package storage

import (
	"crypto/rand"
	"database/sql"
	"testing"

	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	assert := require.New(t)

	connectionString := ":memory:"
	db, err := sql.Open(driverName, connectionString)
	assert.Nil(err, "No error opening database")
	assert.Nil(createSchema(db), "Error creating db in %s", connectionString)
	db.Close()

	s, err := CreateStorage(connectionString)
	assert.Nil(err, "Did not expect error: %v", err)
	defer s.Close()

	testApplicationStorage(s, t)
	testDeviceStorage(s, t)
	testDataStorage(s, t)
	testGatewayStorage(s, t)

	testSimpleKeySequence(s, t)
	testMultipleSequences(s, t)
	testConcurrentSequences(s, t)
	testDownstreamStorage(s, t)
}

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
