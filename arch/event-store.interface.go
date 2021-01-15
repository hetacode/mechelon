package arch

import goeh "github.com/hetacode/go-eh"

// EventStore is an interface to event store db implementation
type EventStore interface {
	GetSnapshot(key string, stateType interface{}) (state interface{})
	SaveSnapshot(key string, state interface{}) error
	GetEvents(key string, position int64) []goeh.Event
	SaveNewEvents(key string, events []goeh.Event) (int64, error)
}
