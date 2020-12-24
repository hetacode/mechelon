package smgeventstore

import (
	"log"
	"time"

	goeh "github.com/hetacode/go-eh"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
)

// ServiceAggregator - an obejct to retrive state, generate events and save them to the event store
type ServiceAggregator struct {
	ID      string
	Version int64
	State   *ServiceStateEntity

	pendingEvents []goeh.Event
}

// NewServiceAggregator instance
func NewServiceAggregator() *ServiceAggregator {
	a := &ServiceAggregator{
		pendingEvents: make([]goeh.Event, 0),
	}
	return a
}

// Replay all needed events with given snapshot state
func (a *ServiceAggregator) Replay(state *ServiceStateEntity, events []goeh.Event) { // TODO: Maybe return new instance of Aggreagtor
	a.State = state
	if a.State == nil {
		a.State = &ServiceStateEntity{}
	}
	for _, ev := range events {
		switch ev.GetType() {
		case new(eventsservicesmgmt.RegisterServiceCommand).GetType():
			e := ev.(*eventsservicesmgmt.RegisterServiceCommand)
			// TODO: move to separate function
			if a.State.ServiceName == "" { // state doesn't exist
				a.State.ServiceName = e.ServiceName
				a.State.ProjectName = e.ProjectName
				a.State.Instances = []ServiceInstance{
					{
						Name:      e.InstanceName,
						CreatedAt: time.Now().Unix(),
						State:     Active,
					},
				}
			} else { // add instance to existing service instance state
				a.State.Instances = append(a.State.Instances, ServiceInstance{
					Name:      e.InstanceName,
					CreatedAt: time.Now().Unix(),
					State:     Active,
				})
			}
		case new(eventsservicesmgmt.UnregisterServiceCommand).GetType():
			e := ev.(*eventsservicesmgmt.UnregisterServiceCommand)
			if a.State.ServiceName == "" {
				log.Printf("err: cannot UnregisterServiceCommand because service '%s' for project '%s' doesn't exist", e.ServiceName, e.ProjectName)
				break
			} else {
				// TODO: check nullability of a.state.instances
				idx := -1
				for i, ins := range a.State.Instances {
					if ins.Name == e.InstanceName {
						idx = i
						break
					}
				}
				if idx < 0 {
					log.Printf("err: cannot UnregisterServiceCommand because instance '%s' for service '%s' from project '%s' doesn't exists", e.InstanceName, e.ServiceName, e.ProjectName)
					break
				}
				a.State.Instances[idx] = a.State.Instances[len(a.State.Instances)-1]
				// a.State.Instances[len(a.State.Instances)-1] = ServiceInstance{} - unnecessary?
				a.State.Instances = a.State.Instances[:len(a.State.Instances)-1]
			}
		}
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
// TODO:
