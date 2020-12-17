package smgeventhandlers

import (
	goeh "github.com/hetacode/go-eh"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
)

// UnregisterServiceEventHandler struct
type UnregisterServiceEventHandler struct {
}

// Handle UnregisterServiceEvent
func (e *UnregisterServiceEventHandler) Handle(event goeh.Event) {
	ev := event.(*eventsservicesmgmt.UnregisterServiceEvent)
	panic(ev)
}
