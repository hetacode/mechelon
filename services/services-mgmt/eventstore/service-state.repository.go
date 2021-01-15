package smgeventstore

import (
	"fmt"

	"github.com/hetacode/mechelon/arch"
)

type ServiceStateRepository struct {
	// TODO:
	// message bus instance

	EventStore arch.EventStore
}

// GetAggregator - create aggregator instance with rebuilt state
func (r *ServiceStateRepository) GetAggregator(projectName, serviceName string) *ServiceAggregator {
	key := fmt.Sprintf("%s-%s", projectName, serviceName)
	state := r.EventStore.GetSnapshot(key, new(ServiceEventStore)).(*ServiceStateEntity)
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
