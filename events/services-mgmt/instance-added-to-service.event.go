package eventsservicesmgmt

import goeh "github.com/hetacode/go-eh"

// InstanceAddedToServiceEvent service instance has been added to service for given project
type InstanceAddedToServiceEvent struct {
	*goeh.EventData
	// should be unique
	ProjectName string `json:"project_name"`
	// i.e.: hostname
	InstanceName string `json:"instance_name"`
	// should be unique per project
	ServiceName string `json:"service_name"`
	// time
	CreateAt int64 `json:"create_at"`

	Version uint64 `json:"version"`
}

// GetType of event
func (e *InstanceAddedToServiceEvent) GetType() string {
	return "InstanceAddedToServiceEvent"
}

// GetVersion of event
func (e *InstanceAddedToServiceEvent) GetVersion() uint64 {
	return e.Version
}

// SetVersion of event
func (e *InstanceAddedToServiceEvent) SetVersion(v uint64) {
	e.Version = v
}
