package smgcommandhandlers

import (
	"log"

	goeh "github.com/hetacode/go-eh"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
	smgcontainer "github.com/hetacode/mechelon/services/services-mgmt/container"
)

// RemoveServiceInstanceCommandHandler struct
type RemoveServiceInstanceCommandHandler struct {
	Container *smgcontainer.Container
}

// Handle UnregisterServiceEvent
func (e *RemoveServiceInstanceCommandHandler) Handle(event goeh.Event) {
	log.Printf("RemoveServiceInstanceCommandHandler start")
	ev := event.(*eventsservicesmgmt.RemoveServiceInstanceCommand)

	aggr := e.Container.ServiceStateRepository.GetAggregator(ev.ProjectName, ev.ServiceName)
	aggr.RemoveInstanceFromService(ev.ProjectName, ev.ServiceName, ev.InstanceName)

	newEvents := aggr.GetPendingEvents()
	if err := e.Container.ServiceStateRepository.SaveEvents(ev.ProjectName, ev.ServiceName, newEvents); err != nil {
		log.Printf("RemoveServiceInstanceCommandHandler SaveEvents err: %s", err)
	}

	log.Printf("RemoveServiceInstanceCommandHandler end")
}
