package smgcommandhandlers

import (
	goeh "github.com/hetacode/go-eh"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
	smgcontainer "github.com/hetacode/mechelon/services/services-mgmt/container"
)

// RegisterServiceCommandHandler struct
type RegisterServiceCommandHandler struct {
	Container smgcontainer.Container
}

// Handle RegisterServiceEvent
func (e *RegisterServiceCommandHandler) Handle(event goeh.Event) {
	ev := event.(*eventsservicesmgmt.RegisterServiceCommand)
	panic(ev)
}
