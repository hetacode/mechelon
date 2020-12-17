package eventsservicesmgmt

import goeh "github.com/hetacode/go-eh"

// RegisterServiceEvent - represent command when given service is running state - so just show being
type RegisterServiceEvent struct {
	*goeh.EventData
	// should be unique
	ProjectName string `json:"project_name"`
	// should be unique per project
	ServiceName string `json:"service_name"`
}

// GetType of event
func (e *RegisterServiceEvent) GetType() string {
	return "RegisterServiceEvent"
}
