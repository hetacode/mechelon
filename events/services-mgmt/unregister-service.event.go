package eventsservicesmgmt

import goeh "github.com/hetacode/go-eh"

// UnregisterServiceEvent - represent command when given service went to stop/kill state
type UnregisterServiceEvent struct {
	*goeh.EventData
	// should be unique
	ProjectName string `json:"project_name"`
	// should be unique per project
	ServiceName string `json:"service_name"`
}

// GetType of event
func (e *UnregisterServiceEvent) GetType() string {
	return "UnregisterServiceEvent"
}
