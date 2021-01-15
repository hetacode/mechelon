package smgcommandhandlers

import (
	goeh "github.com/hetacode/go-eh"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
	smgcontainer "github.com/hetacode/mechelon/services/services-mgmt/container"
)

// UnregisterServiceCommandHandler struct
type UnregisterServiceCommandHandler struct {
	Container *smgcontainer.Container
}

// Handle UnregisterServiceEvent
func (e *UnregisterServiceCommandHandler) Handle(event goeh.Event) {
	ev := event.(*eventsservicesmgmt.UnregisterServiceCommand)
	panic(ev)
}
