//+build wireinject

package svvcontainer

import (
	"os"

	"github.com/google/wire"
	gobus "github.com/hetacode/go-bus"
	goeh "github.com/hetacode/go-eh"
	"github.com/hetacode/mechelon/events"
)

// Container struct keep of all dependencies
type Container struct {
	EventsConsumerBus gobus.ServiceBus
}

// NewContainer instance
func NewContainer() *Container {
	wire.Build(
		events.NewEventsMapper,
		initEventsConsumerBusProvider,
		wire.Struct(new(Container), "*"),
	)
	return nil
}

func initEventsConsumerBusProvider(em *goeh.EventsMapper) gobus.ServiceBus {
	kind := gobus.RabbitMQServiceBusOptionsFanOutKind
	bus := gobus.NewRabbitMQServiceBus(em, &gobus.RabbitMQServiceBusOptions{
		Kind:      &kind,
		Exchanage: os.Getenv("SVC_SERVICES_MGMT_SB_EVENTS_EXCHANGE"),
		Queue:     os.Getenv("SVC_SERVICES_VIEW_SB_EVENTS_QUEUE"),
		Server:    os.Getenv("RABBITMQ_SERVER"),
	})

	return bus
}
