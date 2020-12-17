package events

import (
	goeh "github.com/hetacode/go-eh"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
)

// NewEventsMapper create new instance of events mapper with registered events
func NewEventsMapper() *goeh.EventsMapper {
	m := new(goeh.EventsMapper)
	m.Register(new(eventsservicesmgmt.RegisterServiceEvent))
	m.Register(new(eventsservicesmgmt.UnregisterServiceEvent))

	return m
}
