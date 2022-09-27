package storage

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	assert := require.New(t)

	dbFile := fmt.Sprintf("%s/lospan.db", os.TempDir())
	defer os.Remove(dbFile)

	connectionString := fmt.Sprintf("file:%s", dbFile)
	db, err := sql.Open(DriverName, connectionString)
	assert.Nil(err, "No error opening database")
	assert.Nil(createSchema(db), "Error creating db in %s", dbFile)
	db.Close()

	s, err := CreateStorage(connectionString, 10, 5, time.Minute)
	assert.Nil(err, "Did not expect error: %v", err)
	defer s.Close()

	doStorageTests(s, t)
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

// DoStorageTests tests all of the storage interfaces
func doStorageTests(store *Storage, t *testing.T) {

	testApplicationStorage(store, t)
	testDeviceStorage(store, t)
	testDataStorage(store, t)
	testGatewayStorage(store, t)

	testSimpleKeySequence(store, t)
	testMultipleSequences(store, t)
	testConcurrentSequences(store, t)
	testDownstreamStorage(store, t)

}
