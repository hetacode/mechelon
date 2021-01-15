package smgeventstore

import (
	"fmt"

	goeh "github.com/hetacode/go-eh"
	"github.com/hetacode/mechelon/arch"
)

// ServiceStateRepository struct
type ServiceStateRepository struct {
	EventStore arch.EventStore
}

// GetAggregator - create aggregator instance with rebuilt state
func (r *ServiceStateRepository) GetAggregator(projectName, serviceName string) *ServiceAggregator {
	key := fmt.Sprintf("%s-%s", projectName, serviceName)
	var state *ServiceStateEntity
	stateSnap := r.EventStore.GetSnapshot(key, new(ServiceStateEntity))
	if stateSnap != nil {
		state = stateSnap.(*ServiceStateEntity)
	}
	position := int64(0)
	if state != nil {
		position = state.GetVersion()
	}

	events := r.EventStore.GetEvents(key, position)
	aggr := NewServiceAggregator()
	aggr.Replay(state, events)
	if len(events) >= 20 {
		if err := r.EventStore.SaveSnapshot(key, aggr.State); err != nil {
			panic(err)
		}
	}

	return aggr
}

// SaveEvents to event store
func (r *ServiceStateRepository) SaveEvents(projectName, serviceName string, events []goeh.Event) error {
	key := fmt.Sprintf("%s-%s", projectName, serviceName)
	if err := r.EventStore.SaveNewEvents(key, events); err != nil {
		return err
	}

	return nil
}
