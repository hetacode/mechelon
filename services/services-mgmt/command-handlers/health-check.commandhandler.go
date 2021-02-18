package smgcommandhandlers

import (
	"log"

	goeh "github.com/hetacode/go-eh"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
	smgcontainer "github.com/hetacode/mechelon/services/services-mgmt/container"
)

// HealthCheckCommandHandler struct
type HealthCheckCommandHandler struct {
	Container *smgcontainer.Container
}

// Handle HealthCheckCommand
func (e *HealthCheckCommandHandler) Handle(event goeh.Event) {
	log.Printf("HealthCheckCommandHandler start")
	ev := event.(*eventsservicesmgmt.HealthCheckCommand)

	aggr := e.Container.ServiceStateRepository.GetAggregator(ev.ProjectName, ev.ServiceName)
	aggr.ServiceInstanceHealthCheck(ev.ProjectName, ev.ServiceName, ev.InstanceName)

	newEvents := aggr.GetPendingEvents()
	if err := e.Container.ServiceStateRepository.SaveEvents(ev.ProjectName, ev.ServiceName, newEvents); err != nil {
		log.Printf("HealthCheckCommandHandler SaveEvents err: %s", err)
	}

	for _, ev := range newEvents {
		e.Container.EventsProducerBus.Publish(ev)
	}
	log.Printf("HealthCheckCommandHandler end")
}
