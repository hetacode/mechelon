//+build wireinject

package svvcontainer

import "github.com/google/wire"

// Container struct keep of all dependencies
type Container struct {
}

// NewContainer instance
func NewContainer() *Container {
	wire.Build(
		wire.Struct(new(Container), "*"),
	)
	return nil
}
