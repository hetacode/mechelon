package arch

// ExtendedEvent with GetVersion of event
type ExtendedEvent interface {
	GetVersion() uint64
	SetVersion(v uint64)
}
