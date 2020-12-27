package eventsservicesmgmt

import goeh "github.com/hetacode/go-eh"

// ProjectServiceRemovedEvent new service has been created for given project
type ProjectServiceRemovedEvent struct {
	*goeh.EventData
	// should be unique
	ProjectName string `json:"project_name"`
	// should be unique per project
	ServiceName string `json:"service_name"`
}

// GetType of event
func (e *ProjectServiceRemovedEvent) GetType() string {
	return "ProjectServiceRemovedEvent"
}
