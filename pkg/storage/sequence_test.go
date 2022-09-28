package storage

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// SimpleKeySequence tests a simple sequence
func testSimpleKeySequence(seq *Storage, t *testing.T) {
	assert := require.New(t)

	const numKeys = 10
	received := 0
	ids, err := seq.AllocateKeys("something", numKeys+1, 0)
	assert.NoError(err)

	previous, ok := <-ids
	assert.True(ok)

	for {
		val, ok := <-ids
		if !ok {
			break
		}
		assert.Equal(val, previous+1, "Numbers should be in sequence. Old = %d, new = %d", previous, val)

		previous = val
		received++
	}

	assert.Equal(numKeys, received, "Expected %d keys", numKeys)

}

// MultipleSequences tests two sequences in parallel
func testMultipleSequences(seq *Storage, t *testing.T) {
	const num1 = 10
	const num2 = 25
	some, err := seq.AllocateKeys("something", num1, 1)
	if err != nil {
		t.Fatal("Could not allocate keys: ", err)
	}
	other, err := seq.AllocateKeys("other", num2, 1)
	if err != nil {
		t.Fatal("Could not allocate keys: ", err)
	}
	received := 0
	prevsome := uint64(0)
	prevother := uint64(0)
	for i := 0; i < (num1+num2)*2; i++ {
		if someval, ok := <-some; ok {
			if someval < prevsome {
				t.Errorf("Got equal to or lesser val in some iteration %d", i)
			}
			prevsome = someval
			received++
		}
		if otherval, ok := <-other; ok {
			if otherval < prevother {
				t.Errorf("Got equal to or lesser val in other iteration %d", i)
			}
			received++
		}
	}
	if received < (num1 + num2) {
		t.Fatalf("Got %d keys expected %d", received, (num1 + num2))
	}
}

// ConcurrentSequences tests concurrent retrieval from sequences. Each number
// should be bigger than the old
func testConcurrentSequences(seq *Storage, t *testing.T) {

	const interval = 100

	test := func(wg *sync.WaitGroup) {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			some, err := seq.AllocateKeys("concurrent", interval+1, 0)
			if err != nil {
				t.Error("Got error creating sequence: ", err)
				return
			}
			previous, ok := <-some
			if !ok {
				t.Error("Sequence is closed. ")
				return
			}
			for j := 0; j < interval; j++ {
				val, ok := <-some
				if !ok {
					t.Error("Sequence is closed. Didn't expect it to close now.")
				}
				if val <= previous {
					t.Errorf("Got same or smaller ID. Last was %d, current is %d", previous, val)
					return
				}
			}
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(3)
	go test(&wg)
	go test(&wg)
	go test(&wg)

	wg.Wait()
}
