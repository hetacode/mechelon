// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package gtwcontainer

// Injectors from container.go:

func NewContainer() *Container {
	container := &Container{}
	return container
}

// container.go:

// Container struct keeping all of the required dependencies which are linked together
type Container struct {
}
