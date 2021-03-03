//+build wireinject

package gtwcontainer

import "github.com/google/wire"

// Container struct keeping all of the required dependencies which are linked together
type Container struct {
}

// NewContainer constructor of Container
func NewContainer() *Container {
	wire.Build(
		wire.Struct(new(Container), "*"),
	)

	return nil
}
