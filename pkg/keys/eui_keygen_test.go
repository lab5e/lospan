package keys

import (
	"fmt"
	"testing"

	"github.com/lab5e/lospan/pkg/protocol"
	"github.com/lab5e/lospan/pkg/storage"
)

func TestSimpleKeygen(t *testing.T) {
	// Make a MA-L for the keys
	ma, _ := protocol.NewMA([]byte{1, 2, 3})

	// NetID
	netID := uint32(0)

	storage := storage.NewMemoryStorage()

	keygen, err := NewEUIKeyGenerator(ma, netID, storage)
	if err != nil {
		t.Fatal("Couldn't create key generator: ", err)
	}

	generatedAppEUIs := make([]protocol.EUI, 0)
	generatedDeviceEUIs := make([]protocol.EUI, 0)
	generatedGenericKeys := make([]uint64, 0)

	const numKeys int = 2000
	for i := 0; i < numKeys; i++ {
		newAppEUI, err := keygen.NewAppEUI()
		if err != nil {
			t.Fatal("Got error generating app EUI: ", err)
		}
		generatedAppEUIs = append(generatedAppEUIs, newAppEUI)
		newDeviceEUI, err := keygen.NewDeviceEUI()
		if err != nil {
			t.Fatal("Got error generating device EUI: ", err)
		}
		generatedDeviceEUIs = append(generatedDeviceEUIs, newDeviceEUI)
		generatedGenericKeys = append(generatedGenericKeys, keygen.NewID("generic"))
	}
	// Ensure there's no collisions
	for i := 0; i < numKeys; i++ {
		for j := 0; j < numKeys; j++ {
			if j == i {
				continue
			}
			if generatedAppEUIs[i] == generatedAppEUIs[j] {
				t.Fatalf("Identical app EUI for i = %d, j = %d: %s", i, j, generatedAppEUIs[i])
			}
			if generatedDeviceEUIs[i] == generatedDeviceEUIs[j] {
				t.Fatalf("Identical device EUI for i = %d, j = %d: %s", i, j, generatedDeviceEUIs[i])
			}
			if generatedGenericKeys[i] == generatedGenericKeys[j] {
				t.Fatalf("Identical ID for i = %d, j = %d: %d", i, j, generatedGenericKeys[i])
			}
		}
	}
}

// Ensure you can't generate EUIs with MA-S and a NetID > 0x0F, MA-M with
// NetID > 0x0FFF or MA-L with NetID > 0xFFFF (we can't guarantee uniqueness)
func TestKeygenWithTooLargeNetID(t *testing.T) {
	storage := storage.NewMemoryStorage()

	// 25 bytes are used for the ID
	// MA-S is 36 bits, 64-36-25=3 bits for NetID
	maSmall, err := protocol.NewMA([]byte{1, 2, 3, 4, 5})
	if err != nil {
		t.Fatal("Could not create MA-S EUI: ", err)
	}

	_, err = NewEUIKeyGenerator(maSmall, protocol.MaxNetworkBitsMAS+1, storage)
	if err == nil {
		t.Error("Expected error when creating key generator with too big NetID and MA-S EUI")
	}
	// Then something that will fit exactly
	_, err = NewEUIKeyGenerator(maSmall, protocol.MaxNetworkBitsMAS, storage)
	if err != nil {
		t.Error("Couldn't create key generator with MA-S/3 bit NetID: ", err)
	}

	// MA-M is 28 bits, 64-28-25=11 bits for NetID
	maMedium, err := protocol.NewMA([]byte{1, 2, 3, 4})
	if err != nil {
		t.Fatal("Could not create MA-M EUI: ", err)
	}

	// Create something with too big NetID
	_, err = NewEUIKeyGenerator(maMedium, protocol.MaxNetworkBitsMAM+1, storage)
	if err == nil {
		t.Error("Expected error when creating key generator with too big NetID and MA-M EUI")
	}
	// ...then something that will fit exactly
	_, err = NewEUIKeyGenerator(maMedium, protocol.MaxNetworkBitsMAM, storage)
	if err != nil {
		t.Error("Expected MA-M EIU and 11 bit NetID to fit: ", err)
	}

	// MA-L is 24 bits, 64-24-25=15 bits for NetID
	maLarge, err := protocol.NewMA([]byte{1, 2, 3})
	if err != nil {
		t.Fatal("Could not create MA-L EUI: ", err)
	}

	_, err = NewEUIKeyGenerator(maLarge, protocol.MaxNetworkBitsMAL+1, storage)
	if err == nil {
		t.Error("Expected error when using NetID > 15 bits")
	}
	_, err = NewEUIKeyGenerator(maLarge, protocol.MaxNetworkBitsMAL, storage)
	if err != nil {
		t.Error("Expected MA-L and 15 bit NetID to fit: ", err)
	}
}

// Ensure different NetIDs generate different EUIs for devices and applications
// even with the same sequence numbers
func TestKeygenWithDifferentNetID(t *testing.T) {
	storage := storage.NewMemoryStorage()
	// Use the same MA for both
	ma, _ := protocol.NewMA([]byte{0, 1, 2, 3, 4})
	keygen1, _ := NewEUIKeyGenerator(ma, 0, storage)
	keygen2, _ := NewEUIKeyGenerator(ma, 1, storage)

	eui1 := make([]protocol.EUI, 0)
	eui2 := make([]protocol.EUI, 0)

	keyCount := 1000
	for i := 0; i < keyCount; i++ {
		newEUI1, err := keygen1.NewDeviceEUI()
		if err != nil {
			t.Fatal("Got error generating EUI1: ", err)
		}
		eui1 = append(eui1, newEUI1)
		newEUI2, err := keygen2.NewDeviceEUI()
		if err != nil {
			t.Fatal("Got error generating EUI2: ", err)
		}
		eui2 = append(eui2, newEUI2)
	}

	for i := 0; i < keyCount; i++ {
		for j := 0; j < keyCount; j++ {
			if eui1[i] == eui2[j] {
				t.Errorf("Duplicate key for EUI1 (i=%d, j=%d). EUI=%s", i, j, eui1[i])
			}
		}
	}
}

// Ensure keys are allocated in a lazy fashion, ie they won't be allocated
// until someone retrieves a key.
func TestLazyInvocation(t *testing.T) {
	ksStorage := storage.NewMemoryStorage()

	ma, _ := protocol.NewMA([]byte{0, 1, 2, 3, 4})
	netID := uint32(0)
	identifier := "neteui"
	// Start by allocating keys directly
	name := fmt.Sprintf("%s/%04x/%s", ma.String(), netID, identifier)
	ch, err := ksStorage.AllocateKeys(name, 1, 1)
	if err != nil {
		t.Fatal("Got error allocating keys (first time)")
	}
	firstID := uint64(0)
	for v := range ch {
		firstID = v
	}
	kg1, _ := NewEUIKeyGenerator(ma, netID, ksStorage)
	kg2, _ := NewEUIKeyGenerator(ma, netID, ksStorage)

	// This ensures the keygen has started all of its goroutines
	kg1.NewID("foo1")
	kg2.NewID("foo2")

	// Repeat allocation. The ID should be the next in sequence
	ch, err = ksStorage.AllocateKeys(name, 1, 1)
	if err != nil {
		t.Fatal("Got error allocating keys (second time)")
	}
	lastID := uint64(0)
	for v := range ch {
		lastID = v
	}

	if lastID != (firstID + 1) {
		t.Fatalf("Expected next sequence to start at %d but it started at %d", (firstID + 1), lastID)
	}
}

// Simple benchmark for key generator. Grab 10 keys at a time
func BenchmarkKeygen(b *testing.B) {
	seqStorage := storage.NewMemoryStorage()
	ma, _ := protocol.NewMA([]byte{5, 4, 3, 2, 1})
	netID := uint32(1)
	memkeyGenerator, _ := NewEUIKeyGenerator(ma, netID, seqStorage)

	for i := 0; i < b.N; i++ {
		_, err := memkeyGenerator.NewAppEUI()
		if err != nil {
			b.Fatal("Got error retrieving key: ", err)
		}
	}
}
