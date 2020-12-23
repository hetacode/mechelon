package smgeventstore

import (
	goeh "github.com/hetacode/go-eh"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
)

// ServiceAggregator - an obejct to retrive state, generate events and save them to the event store
type ServiceAggregator struct {
	ID      string
	Version int64
	State   *ServiceStateEntity

	pendingEvents []goeh.Event
}

// NewServiceAggregator instance
func NewServiceAggregator() *ServiceAggregator {
	a := &ServiceAggregator{
		pendingEvents: make([]goeh.Event, 0),
	}
	return a
}

// Replay all needed events with given snapshot state
func (a *ServiceAggregator) Replay(state *ServiceStateEntity, events []goeh.Event) { // TODO: Maybe return new instance of Aggreagtor
	a.State = state
	for _, ev := range events {
		switch ev.GetType() {
		case new(eventsservicesmgmt.RegisterServiceCommand).GetType():
			e := ev.(*eventsservicesmgmt.RegisterServiceCommand)
			panic(e)
			// TODO: create needded modyfications
		case new(eventsservicesmgmt.UnregisterServiceCommand).GetType():
			e := ev.(*eventsservicesmgmt.RegisterServiceCommand)
			panic(e)
			// create needded modyfications
		}
	}
}

func (a *ServiceAggregator) GetPendingEvents() []goeh.Event {
	return a.pendingEvents
}
