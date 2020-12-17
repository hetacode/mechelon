package smgeventhandlers

import (
	goeh "github.com/hetacode/go-eh"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
)

// RegisterServiceEventHandler struct
type RegisterServiceEventHandler struct {
}

// Handle RegisterServiceEvent
func (e *RegisterServiceEventHandler) Handle(event goeh.Event) {
	ev := event.(*eventsservicesmgmt.RegisterServiceEvent)
	panic(ev)
}
