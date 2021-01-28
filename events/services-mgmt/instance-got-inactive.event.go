package eventsservicesmgmt

import goeh "github.com/hetacode/go-eh"

// InstanceGotInactiveEvent service instance has been set as inactive
type InstanceGotInactiveEvent struct {
	*goeh.EventData
	// should be unique
	ProjectName string `json:"project_name"`
	// i.e.: hostname
	InstanceName string `json:"instance_name"`
	// should be unique per project
	ServiceName string `json:"service_name"`
	// time
	UpdateAt int64 `json:"update_at"`

	Version uint64 `json:"version"`
}

// GetType of event
func (e *InstanceGotInactiveEvent) GetType() string {
	return "InstanceGotInactiveEvent"
}

// GetVersion of event
func (e *InstanceGotInactiveEvent) GetVersion() uint64 {
	return e.Version
}

// SetVersion of event
func (e *InstanceGotInactiveEvent) SetVersion(v uint64) {
	e.Version = v
}
