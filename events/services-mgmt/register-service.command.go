package eventsservicesmgmt

import goeh "github.com/hetacode/go-eh"

// RegisterServiceCommand - register new service and/or service instace
type RegisterServiceCommand struct {
	*goeh.EventData
	// should be unique
	ProjectName string `json:"project_name"`
	// i.e.: hostname
	InstanceName string `json:"instance_name"`
	// should be unique per project
	ServiceName string `json:"service_name"`
}

// GetType of event
func (e *RegisterServiceCommand) GetType() string {
	return "RegisterServiceCommand"
}
