package eventsservicesmgmt

import goeh "github.com/hetacode/go-eh"

// UnregisterServiceCommand - represent command when given service went to stop/kill state
type UnregisterServiceCommand struct {
	*goeh.EventData
	// should be unique
	ProjectName string `json:"project_name"`
	// i.e.: hostname
	InstanceName string `json:"instance_name"`
	// should be unique per project
	ServiceName string `json:"service_name"`
}

// GetType of event
func (e *UnregisterServiceCommand) GetType() string {
	return "UnregisterServiceEvent"
}
