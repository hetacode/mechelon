package smgeventstore

import "github.com/EventStore/EventStore-Client-Go/client"

type ServiceStateRepository struct {
	// TODO: event store instance
	// snapshot store instance - for EventStoreDB will be the same as event store
	// message bus instance

	EventStoreClient *client.Client
}

func (r *ServiceStateRepository) GetAggregator(projectName, serviceName string) *ServiceAggregator {
	// TODO: get snapshot from EventStore
	return nil
}
