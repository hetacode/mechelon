//+build wireinject

package gtwcontainer

// Container struct keeping all of the required dependencies which are linked together
type Container struct {
}

// NewContainer constructor of Container
func NewContainer() (container *Container) {
	container = &Container{}

	return
}
