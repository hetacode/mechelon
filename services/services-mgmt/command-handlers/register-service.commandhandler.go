package smgcommandhandlers

import (
	"log"

	goeh "github.com/hetacode/go-eh"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
	smgcontainer "github.com/hetacode/mechelon/services/services-mgmt/container"
)

// RegisterServiceCommandHandler struct
type RegisterServiceCommandHandler struct {
	Container *smgcontainer.Container
}

// Handle RegisterServiceEvent
func (e *RegisterServiceCommandHandler) Handle(event goeh.Event) {
	log.Printf("RegisterServiceCommandHandler start")
	ev := event.(*eventsservicesmgmt.RegisterServiceCommand)
	aggr := e.Container.ServiceStateRepository.GetAggregator(ev.ProjectName, ev.ServiceName)
	aggr.RegisterNewService(ev.ProjectName, ev.ServiceName, ev.InstanceName)

	newEvents := aggr.GetPendingEvents()
	if err := e.Container.ServiceStateRepository.SaveEvents(ev.ProjectName, ev.ServiceName, newEvents); err != nil {
		log.Printf("RegisterServiceCommandHandler SaveEvents err: %s", err)
	}

	for _, nev := range newEvents {
		e.Container.EventsProducerBus.Publish(nev)

		// New service for project has been created - let's run worker job
		if nev.GetType() == new(eventsservicesmgmt.ProjectServiceCreatedEvent).GetType() {
			e.Container.WorkersManager.CreateWorker(ev.ProjectName, ev.ServiceName)
		}
	}
	log.Printf("RegisterServiceCommandHandler end")
}
