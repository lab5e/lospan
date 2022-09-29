package server

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// Simple one-shot route test
func TestEventRouter(t *testing.T) {

	router := NewEventRouter[int, string](2)

	ch := router.Subscribe(0)

	router.Publish(0, "inactive")

	select {
	case <-ch:
		// This is ok
	case <-time.After(10 * time.Millisecond):
		t.Fatal("Didn't get an event on the channel")
	}

	router.Unsubscribe(ch)
}

// Test with multiple routes (and channels)
func TestEventRouterMultipleRoutes(t *testing.T) {
	const numEvents = 4
	router := NewEventRouter[uint64, string](numEvents)
	wg := sync.WaitGroup{}

	const routes = 10
	ids := make([]uint64, routes)
	for i := 0; i < routes; i++ {
		ids[i] = uint64(i)
	}

	chans := make([]<-chan string, routes)

	for i := 0; i < routes; i++ {
		chans[i] = router.Subscribe(ids[i])
	}

	wg.Add(routes)
	for _, ch := range chans {
		ch := ch
		go func() {
			received := 0
			for {
				select {
				case <-ch:
					received++
					if received == numEvents {
						wg.Done()
						return
					}
				case <-time.After(100 * time.Millisecond):
					panic(fmt.Sprintf("Didn't receive data! Got just %d events, expected 5", received))
				}
			}
		}()
	}

	publish := func() {
		for i := 0; i < routes; i++ {
			router.Publish(ids[i], "keepalive")
			router.Publish(ids[i], "keepalive")
			router.Publish(ids[i], "some data")
			router.Publish(ids[i], "some data")
		}
	}

	publish()

	wg.Wait()

	for i := routes - 1; i >= 0; i-- {
		router.Unsubscribe(chans[i])
	}

	publish()
}

// Create multiple copies of the same subscription and size up and down. The
// output isn't *that* interesting; the test just ensures edge cases aren't missed.
func TestResize(t *testing.T) {
	const routeCount = 100
	router := NewEventRouter[uint64, string](2)

	var subs []<-chan string

	id := uint64(12)
	for i := 0; i < routeCount; i++ {
		ch := router.Subscribe(id)
		subs = append(subs, ch)
	}

	// Publish one
	router.Publish(id, "inactive")

	for i := 0; i < routeCount/2; i++ {
		router.Unsubscribe(subs[rand.Int()%routeCount])
	}

	router.Publish(id, "keepalive")

	for i := 0; i < routeCount; i++ {
		router.Unsubscribe(subs[i])
	}
}

// ratio of miss/hit
const ratio = 50

func setupBenchmark(count int) (*EventRouter[int, int], []<-chan int) {
	e := NewEventRouter[int, int](1)
	rand.Seed(time.Now().UnixNano())
	chs := make([]<-chan int, count)
	for i := 0; i < count; i++ {
		if rand.Intn(100) < ratio {
			chs[i] = e.Subscribe(i)
		}
	}
	return &e, chs
}

func runTest(e *EventRouter[int, int], chans []<-chan int, count int) {

	for i := 0; i < count; i++ {
		idx := rand.Intn(len(chans))
		e.Publish(idx, i)
		select {
		case <-chans[idx]:
		default:
		}
	}
}

func BenchmarkEventRouter100(b *testing.B) {
	e, chs := setupBenchmark(100)
	b.ResetTimer()
	runTest(e, chs, b.N)
}

func BenchmarkEventRouter1000(b *testing.B) {
	e, chs := setupBenchmark(1000)
	b.ResetTimer()
	runTest(e, chs, b.N)
}

func BenchmarkEventRouter10000(b *testing.B) {
	e, chs := setupBenchmark(10000)
	b.ResetTimer()
	runTest(e, chs, b.N)
}

func BenchmarkEventRouter100000(b *testing.B) {
	e, chs := setupBenchmark(100000)
	b.ResetTimer()
	runTest(e, chs, b.N)
}
