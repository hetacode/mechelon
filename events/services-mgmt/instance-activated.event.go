package eventsservicesmgmt

import goeh "github.com/hetacode/go-eh"

// InstanceActivatedEvent service instance has been activated
type InstanceActivatedEvent struct {
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
func (e *InstanceActivatedEvent) GetType() string {
	return "InstanceActivatedEvent"
}

// GetVersion of event
func (e *InstanceActivatedEvent) GetVersion() uint64 {
	return e.Version
}

// SetVersion of event
func (e *InstanceActivatedEvent) SetVersion(v uint64) {
	e.Version = v
}
