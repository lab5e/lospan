package gwevents

type gwEventType string

// GwEvent types are OOB events for the gateway. They will be sent
// as a debugging aid for gateways. The gateway interface(s) forwards all
// gateway events to a buffered channel which will distribute the events
// to listeners.
type GwEvent struct {
	Type gwEventType `json:"event"`          // EventType holds the event type (see constants)
	Data string      `json:"data,omitempty"` // The data sent or received from the gateway (if applicable)
}

// NewInactive creates a new inactive event
func NewInactive() GwEvent {
	return GwEvent{gwEventType("Inactive"), ""}
}

// NewKeepAlive creates a new keepalive event
func NewKeepAlive() GwEvent {
	return GwEvent{gwEventType("KeepAlive"), ""}
}

// NewRx creates a new Rx event for the gateway
func NewRx(data string) GwEvent {
	return GwEvent{gwEventType("Rx"), data}
}

// NewTx creates a new Tx event for the gateway
func NewTx(data string) GwEvent {
	return GwEvent{gwEventType("Tx"), data}
}
