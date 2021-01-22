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

	Version uint64 `json:"version"`
}

// GetType of event
func (e *InstanceRemovedFromServiceEvent) GetType() string {
	return "InstanceRemovedFromServiceEvent"
}

// GetVersion of event
func (e *InstanceRemovedFromServiceEvent) GetVersion() uint64 {
	return e.Version
}

// SetVersion of event
func (e *InstanceRemovedFromServiceEvent) SetVersion(v uint64) {
	e.Version = v
}
