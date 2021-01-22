package eventsservicesmgmt

import goeh "github.com/hetacode/go-eh"

// ProjectServiceCreatedEvent new service has been created for given project
type ProjectServiceCreatedEvent struct {
	*goeh.EventData
	// should be unique
	ProjectName string `json:"project_name"`
	// should be unique per project
	ServiceName string `json:"service_name"`

	Version uint64 `json:"version"`
}

// GetType of event
func (e *ProjectServiceCreatedEvent) GetType() string {
	return "ProjectServiceCreatedEvent"
}

// GetVersion of event
func (e *ProjectServiceCreatedEvent) GetVersion() uint64 {
	return e.Version
}

// SetVersion of event
func (e *ProjectServiceCreatedEvent) SetVersion(v uint64) {
	e.Version = v
}
