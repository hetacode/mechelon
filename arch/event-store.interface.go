package arch

import goeh "github.com/hetacode/go-eh"

// EventStore is an interface to event store db implementation
type EventStore interface {
	GetSnapshot(key string, stateType interface{}) (state interface{}, lastEventPosition int64)
	GetEvents(key string, position int64) []goeh.Event
	PushNewEvents(key string, events []goeh.Event) error
}
