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
}

// GetType of event
func (e *InstanceAddedToServiceEvent) GetType() string {
	return "InstanceAddedToServiceEvent"
}
