package server

import (
	"log"
	"sync"
)

type route struct {
	id interface{}
	ch chan interface{}
}

// EventRouter is a channel event router. It will route events (or entities)
// based on the EUI. There may be multiple subscribers to the same EUI and each
// will receive a separate event. The channels are buffered and if the subscribers
// can't keep up with the events they will be dropped silently by the router.
type EventRouter struct {
	routes        []route
	mutex         *sync.Mutex
	channelLength int
}

// NewEventRouter creates a new event router
func NewEventRouter(channelLength int) EventRouter {
	return EventRouter{
		make([]route, 0),
		&sync.Mutex{},
		channelLength,
	}
}

// Subscribe subscribes to events for a particular gateway
func (e *EventRouter) Subscribe(identifier interface{}) <-chan interface{} {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	events := make(chan interface{}, e.channelLength)
	e.routes = append(e.routes, route{identifier, events})

	return events
}

// Unsubscribe from channel
func (e *EventRouter) Unsubscribe(ch <-chan interface{}) {
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
func (e *EventRouter) Publish(id interface{}, event interface{}) {
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
