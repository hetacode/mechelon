package smgcommandhandlers

import (
	goeh "github.com/hetacode/go-eh"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
)

// UnregisterServiceCommandHandler struct
type UnregisterServiceCommandHandler struct {
}

// Handle UnregisterServiceEvent
func (e *UnregisterServiceCommandHandler) Handle(event goeh.Event) {
	ev := event.(*eventsservicesmgmt.UnregisterServiceCommand)
	panic(ev)
}
