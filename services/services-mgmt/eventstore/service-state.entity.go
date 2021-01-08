package smgeventstore

// State basic interface
// Contains GetVersion func to get current version of state
type State interface {
	GetVersion() int64
}

// ServiceStateEntity represent a state of given service - keep status of all active instances
type ServiceStateEntity struct {
	Version     int64             `json:"version"`
	ProjectName string            `json:"project_name"`
	ServiceName string            `json:"service_name"`
	Instances   []ServiceInstance `json:"instances"`
}

// GetVersion of state
func (s *ServiceStateEntity) GetVersion() int64 {
	return s.Version
}

// ServiceInstance basic data like name or creation time
type ServiceInstance struct {
	Name      string       `json:"name"`
	CreatedAt int64        `json:"created_at"`
	State     ServiceState `json:"state"`
}

// ServiceState activity state of service
type ServiceState string

const (
	// Active service - health check is correct
	Active ServiceState = "active_state"
	// Idle - health check doesnt't give presence of service
	Idle ServiceState = "idle_state"
	// InActive service
	InActive ServiceState = "inactive_state"
)
