package events

import (
	goeh "github.com/hetacode/go-eh"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
)

// NewEventsMapper create new instance of events mapper with registered events
func NewEventsMapper() *goeh.EventsMapper {
	m := new(goeh.EventsMapper)
	m.Register(new(eventsservicesmgmt.RegisterServiceCommand))
	m.Register(new(eventsservicesmgmt.UnregisterServiceCommand))
	m.Register(new(eventsservicesmgmt.RemoveServiceInstanceCommand))
	m.Register(new(eventsservicesmgmt.HealthCheckCommand))

	m.Register(new(eventsservicesmgmt.InstanceAddedToServiceEvent))
	m.Register(new(eventsservicesmgmt.InstanceRemovedFromServiceEvent))
	m.Register(new(eventsservicesmgmt.ProjectServiceCreatedEvent))
	m.Register(new(eventsservicesmgmt.ProjectServiceRemovedEvent))
	m.Register(new(eventsservicesmgmt.InstanceActivatedEvent))
	m.Register(new(eventsservicesmgmt.InstanceGotIdleEvent))
	m.Register(new(eventsservicesmgmt.InstanceGotInactiveEvent))

	return m
}
