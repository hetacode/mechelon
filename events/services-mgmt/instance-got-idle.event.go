package eventsservicesmgmt

import goeh "github.com/hetacode/go-eh"

// InstanceGotIdleEvent service instance has been set as idle state
type InstanceGotIdleEvent struct {
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
func (e *InstanceGotIdleEvent) GetType() string {
	return "InstanceGotIdleEvent"
}

// GetVersion of event
func (e *InstanceGotIdleEvent) GetVersion() uint64 {
	return e.Version
}

// SetVersion of event
func (e *InstanceGotIdleEvent) SetVersion(v uint64) {
	e.Version = v
}
