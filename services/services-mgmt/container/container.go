//+build wireinject

package smgcontainer

import (
	"github.com/google/wire"
	gobus "github.com/hetacode/go-bus"
	goeh "github.com/hetacode/go-eh"
	"github.com/hetacode/mechelon/arch"
	"github.com/hetacode/mechelon/events"
	smgeventstore "github.com/hetacode/mechelon/services/services-mgmt/eventstore"
	smgtypes "github.com/hetacode/mechelon/services/services-mgmt/types"
	smgworkers "github.com/hetacode/mechelon/services/services-mgmt/workers"
)

type Container struct {
	EventStore             arch.EventStore
	CommandsConsumerBus    smgtypes.CommandsConsumerBus
	EventsProducerBus      smgtypes.EventsProducerBus
	ServiceStateRepository *smgeventstore.ServiceStateRepository
	WorkersManager         *smgworkers.WorkersManager
}

func NewContainer() *Container {
	wire.Build(
		events.NewEventsMapper,
		initEventStoreProvider,
		initEventsProducerBusProvider,
		initCommandsConsumerBusProvider,
		initServiceStateRepositoryProvider,
		initWorkersManagerProvider,
		wire.Struct(new(Container), "*"),
	)
	return nil
}

func initWorkersManagerProvider(repository *smgeventstore.ServiceStateRepository) *smgworkers.WorkersManager {
	mgr := smgworkers.NewWorkersManager(repository)
	return mgr
}

func initEventsProducerBusProvider(em *goeh.EventsMapper) smgtypes.EventsProducerBus {
	kind := gobus.RabbitMQServiceBusOptionsFanOutKind
	bus := gobus.NewRabbitMQServiceBus(em, &gobus.RabbitMQServiceBusOptions{
		Kind:      &kind,
		Exchanage: "services-mgmt-events-ex",
		Server:    "amqp://localhost:5673",
	})

	return bus
}

func initServiceStateRepositoryProvider(es arch.EventStore) *smgeventstore.ServiceStateRepository {
	r := &smgeventstore.ServiceStateRepository{
		EventStore: es,
	}
	return r
}

func initCommandsConsumerBusProvider(em *goeh.EventsMapper) smgtypes.CommandsConsumerBus {
	kind := gobus.RabbitMQServiceBusOptionsFanOutKind
	bus := gobus.NewRabbitMQServiceBus(em, &gobus.RabbitMQServiceBusOptions{
		Kind:      &kind,
		Exchanage: "services-mgmt-ex",
		Queue:     "services-mgmt-commands-queue",
		Server:    "amqp://localhost:5673",
	})

	return bus
}

func initEventStoreProvider(em *goeh.EventsMapper) arch.EventStore {
	es := smgeventstore.NewServiceEventStore(em)
	return es
}
