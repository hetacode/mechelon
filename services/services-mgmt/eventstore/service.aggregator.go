package smgeventstore

import (
	"fmt"
	"log"
	"time"

	goeh "github.com/hetacode/go-eh"
	"github.com/hetacode/mechelon/arch"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
)

// ServiceAggregator - an obejct to retrive state, generate events and save them to the event store
type ServiceAggregator struct {
	ID    string
	State *ServiceStateEntity

	pendingEvents []goeh.Event
}

// NewServiceAggregator instance
func NewServiceAggregator() *ServiceAggregator {
	a := &ServiceAggregator{
		pendingEvents: make([]goeh.Event, 0),
	}
	return a
}

// GetVersion of state
func (a *ServiceAggregator) GetVersion() int64 {
	if a.State == nil {
		return 0
	}
	return a.State.Version
}

// Replay all needed events with given snapshot state
func (a *ServiceAggregator) Replay(state *ServiceStateEntity, events []goeh.Event) { // TODO: Maybe return new instance of Aggreagtor
	a.State = state
	if a.State == nil {
		a.State = &ServiceStateEntity{}
	}
	for _, ev := range events {
		switch ev.GetType() {
		case new(eventsservicesmgmt.ProjectServiceCreatedEvent).GetType():
			e := ev.(*eventsservicesmgmt.ProjectServiceCreatedEvent)
			if a.State.ServiceName == e.ServiceName {
				log.Printf("err ProjectServiceCreatedEvent service '%s' exist for project '%s'", e.ServiceName, e.ProjectName)
				return
			}
			a.State.ServiceName = e.ServiceName
			a.State.ProjectName = e.ProjectName
			a.State.Instances = make([]ServiceInstance, 0)

		case new(eventsservicesmgmt.InstanceAddedToServiceEvent).GetType():
			e := ev.(*eventsservicesmgmt.InstanceAddedToServiceEvent)
			// TODO: move to separate function
			if a.State.ServiceName == "" { // state doesn't exist
				log.Printf("err InstanceAddedToServiceEvent service '%s' doesn't exist for project '%s'", e.ServiceName, e.ProjectName)
				return
			} else { // add instance to existing service instance state
				a.State.Instances = append(a.State.Instances, ServiceInstance{
					Name:      e.InstanceName,
					CreatedAt: e.CreateAt,
					State:     Active,
				})
			}

		case new(eventsservicesmgmt.ProjectServiceRemovedEvent).GetType():
			e := ev.(*eventsservicesmgmt.ProjectServiceRemovedEvent)
			a.State.Instances = make([]ServiceInstance, 0)
			log.Printf("service '%s' has been removed from '%s' project", e.ServiceName, e.ProjectName)

		case new(eventsservicesmgmt.InstanceRemovedFromServiceEvent).GetType():
			e := ev.(*eventsservicesmgmt.InstanceRemovedFromServiceEvent)
			if a.State.ServiceName != "" {
				log.Printf("err: InstanceRemovedFromServiceEvent service '%s' for project '%s' doesn't exist", e.ServiceName, e.ProjectName)
				break
			} else {
				if a.State.Instances == nil {
					log.Printf("err: InstanceRemovedFromServiceEvent state instances slice is nil for service '%s' in '%s' project", e.ServiceName, e.ProjectName)
					return
				}
				idx := -1
				for i, ins := range a.State.Instances {
					if ins.Name == e.InstanceName {
						idx = i
						break
					}
				}
				if idx < 0 {
					log.Printf("err: cannot InstanceRemovedFromServiceEvent because instance '%s' for service '%s' from project '%s' doesn't exists", e.InstanceName, e.ServiceName, e.ProjectName)
					break
				}
				a.State.Instances[idx] = a.State.Instances[len(a.State.Instances)-1]
				// a.State.Instances[len(a.State.Instances)-1] = ServiceInstance{} - unnecessary?
				a.State.Instances = a.State.Instances[:len(a.State.Instances)-1]
			}
		}
		ee := ev.(arch.ExtendedEvent)
		a.State.Version = int64(ee.GetVersion())
	}
}

// GetPendingEvents - new events
func (a *ServiceAggregator) GetPendingEvents() []goeh.Event {
	return a.pendingEvents
}

// Clear new events
func (a *ServiceAggregator) Clear() {
	a.pendingEvents = make([]goeh.Event, 0)
}

// Modfications of aggregate

// RegisterNewService - generate an event for new service registered in system
func (a *ServiceAggregator) RegisterNewService(projectName, serviceName, instanceName string) {
	id := fmt.Sprintf("%s-%s", projectName, serviceName)
	if a.State.ServiceName == "" {
		createdServiceEv := &eventsservicesmgmt.ProjectServiceCreatedEvent{
			EventData:   &goeh.EventData{ID: id},
			ProjectName: projectName,
			ServiceName: serviceName,
		}
		a.pendingEvents = append(a.pendingEvents, createdServiceEv)
	}

	createInstanceEv := &eventsservicesmgmt.InstanceAddedToServiceEvent{
		EventData:    &goeh.EventData{ID: id},
		ProjectName:  projectName,
		ServiceName:  serviceName,
		InstanceName: instanceName,
		CreateAt:     time.Now().Unix(),
	}
	a.pendingEvents = append(a.pendingEvents, createInstanceEv)
}

// RemoveInstanceFromService - just generate an event about removed instance from service
func (a *ServiceAggregator) RemoveInstanceFromService(projectName, serviceName, instanceName string) {
	id := fmt.Sprintf("%s-%s", projectName, serviceName)

	found := false
	for _, inst := range a.State.Instances {
		if inst.Name == instanceName {
			found = true
		}
	}

	if found {
		ev := &eventsservicesmgmt.InstanceRemovedFromServiceEvent{
			EventData:    &goeh.EventData{ID: id},
			ProjectName:  projectName,
			ServiceName:  serviceName,
			InstanceName: instanceName,
			RemovedAt:    time.Now().Unix(),
		}
		a.pendingEvents = append(a.pendingEvents, ev)
	} else {
		log.Printf("RemoveInstanceFromService cannot find '%s' instance in '%s' service for '%s' project", instanceName, serviceName, projectName)
	}
}

// RemoveService from specific project
// If any instance exists generate events for each one
func (a *ServiceAggregator) RemoveService(projectName, serviceName string) {
	id := fmt.Sprintf("%s-%s", projectName, serviceName)

	instancesEvents := make([]goeh.Event, 0, 1)
	// Get all instances
	for _, inst := range a.State.Instances {
		ev := &eventsservicesmgmt.InstanceRemovedFromServiceEvent{
			EventData:    &goeh.EventData{ID: id},
			ProjectName:  projectName,
			ServiceName:  serviceName,
			InstanceName: inst.Name,
			RemovedAt:    time.Now().Unix(),
		}
		instancesEvents = append(instancesEvents, ev)
	}

	removedServiceEv := &eventsservicesmgmt.ProjectServiceRemovedEvent{
		EventData:   &goeh.EventData{ID: id},
		ProjectName: projectName,
		ServiceName: serviceName,
	}

	a.pendingEvents = append(a.pendingEvents, removedServiceEv)
	a.pendingEvents = append(a.pendingEvents, instancesEvents...)
}
