package eventsservicesmgmt

import goeh "github.com/hetacode/go-eh"

// HealthCheckCommand - should be send periodic by each service instance
type HealthCheckCommand struct {
	*goeh.EventData
	// should be unique
	ProjectName string `json:"project_name"`
	// should be unique per project
	ServiceName string `json:"service_name"`
	// instance name
	InstanceName string `json:"instance_name"`
}

// GetType of event
func (e *HealthCheckCommand) GetType() string {
	return "HealthCheckCommand"
}
