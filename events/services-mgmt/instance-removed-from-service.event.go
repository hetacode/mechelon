package eventsservicesmgmt

import goeh "github.com/hetacode/go-eh"

// InstanceRemovedFromServiceEvent instance has been removed from service
type InstanceRemovedFromServiceEvent struct {
	*goeh.EventData
	// should be unique
	ProjectName string `json:"project_name"`
	// i.e.: hostname
	InstanceName string `json:"instance_name"`
	// should be unique per project
	ServiceName string `json:"service_name"`
	// time
	RemovedAt int64 `json:"removed_at"`
}

// GetType of event
func (e *InstanceRemovedFromServiceEvent) GetType() string {
	return "InstanceRemovedFromServiceEvent"
}
