package svveventhandlers

import (
	"log"

	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
	svvcontainer "github.com/hetacode/mechelon/services/services-view/container"

	goeh "github.com/hetacode/go-eh"
)

// ProjectServiceCreatedEventHandler struct
type ProjectServiceCreatedEventHandler struct {
	Container *svvcontainer.Container
}

// Handle ProjectServiceCreatedEvent
func (e *ProjectServiceCreatedEventHandler) Handle(event goeh.Event) {
	log.Printf("ProjectServiceCreatedEventHandler start")
	ev := event.(*eventsservicesmgmt.ProjectServiceCreatedEvent)

	log.Printf("ProjectServiceCreatedEventHandler end")
	panic(ev)
}
