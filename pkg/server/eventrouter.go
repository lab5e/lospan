package server

import (
	"log"
	"sync"
)

type route[I comparable, T any] struct {
	id I
	ch chan T
}

// EventRouter is a channel event router. It will route events (or entities)
// based on the EUI. There may be multiple subscribers to the same EUI and each
// will receive a separate event. The channels are buffered and if the subscribers
// can't keep up with the events they will be dropped silently by the router.
type EventRouter[I comparable, T any] struct {
	routes        []route[I, T]
	mutex         *sync.Mutex
	channelLength int
}

// NewEventRouter creates a new event router
func NewEventRouter[I comparable, T any](channelLength int) EventRouter[I, T] {
	return EventRouter[I, T]{
		make([]route[I, T], 0),
		&sync.Mutex{},
		channelLength,
	}
}

// Subscribe subscribes to events for a particular gateway
func (e *EventRouter[I, T]) Subscribe(identifier I) <-chan T {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	events := make(chan T, e.channelLength)
	e.routes = append(e.routes, route[I, T]{identifier, events})

	return events
}

// Unsubscribe from channel
func (e *EventRouter[I, T]) Unsubscribe(ch <-chan T) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	for i, r := range e.routes {
		if r.ch == ch {
			close(r.ch)
			e.routes = append(e.routes[:i], e.routes[i+1:]...)
			break
		}
	}
}

// Publish publishes a gateway event to subscribers. If there are no subscribers
// the event will be ignored. If the event subscribers can't keep up with the events
// the events will be silently dropped.
func (e *EventRouter[I, T]) Publish(id I, event T) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	for _, route := range e.routes {
		if route.id == id {
			select {
			case route.ch <- event:
				// This is OK
			default:
				log.Printf("Channel client isn't keeping up with reads. Skipping the event (%v) for ID %v", event, id)
				// Skip event
			}
		}
	}
}
