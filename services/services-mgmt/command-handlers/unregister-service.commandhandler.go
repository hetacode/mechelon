package smgcommandhandlers

import (
	"log"

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
	log.Printf("UnregisterServiceCommandHandler start")
	ev := event.(*eventsservicesmgmt.UnregisterServiceCommand)

	aggr := e.Container.ServiceStateRepository.GetAggregator(ev.ProjectName, ev.ServiceName)
	aggr.RemoveService(ev.ProjectName, ev.ServiceName)

	newEvents := aggr.GetPendingEvents()
	if err := e.Container.ServiceStateRepository.SaveEvents(ev.ProjectName, ev.ServiceName, newEvents); err != nil {
		log.Printf("RegisterServiceCommandHandler SaveEvents err: %s", err)
	}

	for _, ev := range newEvents {
		e.Container.EventsProducerBus.Publish(ev)
	}
	log.Printf("RegisterServiceCommandHandler end")
}
