package eventsservicesmgmt

import goeh "github.com/hetacode/go-eh"

// ProjectServiceRemovedEvent new service has been created for given project
type ProjectServiceRemovedEvent struct {
	*goeh.EventData
	// should be unique
	ProjectName string `json:"project_name"`
	// should be unique per project
	ServiceName string `json:"service_name"`

	Version uint64
}

// GetType of event
func (e *ProjectServiceRemovedEvent) GetType() string {
	return "ProjectServiceRemovedEvent"
}

// GetVersion of event
func (e *ProjectServiceRemovedEvent) GetVersion() uint64 {
	return e.Version
}

// SetVersion of event
func (e *ProjectServiceRemovedEvent) SetVersion(v uint64) {
	e.Version = v
}
