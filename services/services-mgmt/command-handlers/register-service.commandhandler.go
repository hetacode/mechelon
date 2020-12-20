package smgcommandhandlers

import (
	goeh "github.com/hetacode/go-eh"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
)

// RegisterServiceCommandHandler struct
type RegisterServiceCommandHandler struct {
}

// Handle RegisterServiceEvent
func (e *RegisterServiceCommandHandler) Handle(event goeh.Event) {
	ev := event.(*eventsservicesmgmt.RegisterServiceCommand)
	panic(ev)
}
