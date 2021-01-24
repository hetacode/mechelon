package eventsservicesmgmt

import goeh "github.com/hetacode/go-eh"

// RemoveServiceInstanceCommand - remove given instance from service
type RemoveServiceInstanceCommand struct {
	*goeh.EventData
	// should be unique
	ProjectName string `json:"project_name"`
	// should be unique per project
	ServiceName string `json:"service_name"`
	// i.e.: hostname
	InstanceName string `json:"instance_name"`
}

// GetType of event
func (e *RemoveServiceInstanceCommand) GetType() string {
	return "RemoveServiceInstanceCommand"
}
