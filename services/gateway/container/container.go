//+build wireinject

package gtwcontainer

import (
	"os"

	"github.com/google/wire"
	gobus "github.com/hetacode/go-bus"
	goeh "github.com/hetacode/go-eh"
	"github.com/hetacode/mechelon/events"
)

// Container struct keeping all of the required dependencies which are linked together
type Container struct {
	CommandsBusProducer gobus.ServiceBus
}

// NewContainer constructor of Container
func NewContainer() *Container {
	wire.Build(
		events.NewEventsMapper,
		initCommandsBusProducerProvider,
		wire.Struct(new(Container), "*"),
	)

	return nil
}

func initCommandsBusProducerProvider(em *goeh.EventsMapper) gobus.ServiceBus {
	kind := gobus.RabbitMQServiceBusOptionsFanOutKind
	bus := gobus.NewRabbitMQServiceBus(em, &gobus.RabbitMQServiceBusOptions{
		Kind:      &kind,
		Exchanage: os.Getenv("SVC_SERVICES_MGMT_SB_COMMANDS_EXCHANGE"),
		Server:    os.Getenv("RABBITMQ_SERVER"),
	})

	return bus
}
