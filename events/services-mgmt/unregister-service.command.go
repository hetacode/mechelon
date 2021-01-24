package eventsservicesmgmt

import goeh "github.com/hetacode/go-eh"

// UnregisterServiceCommand - unregister whole service, it means even with dependend instances
type UnregisterServiceCommand struct {
	*goeh.EventData
	// should be unique
	ProjectName string `json:"project_name"`
	// should be unique per project
	ServiceName string `json:"service_name"`
}

// GetType of event
func (e *UnregisterServiceCommand) GetType() string {
	return "UnregisterServiceEvent"
}
