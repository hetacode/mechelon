//+build wireinject

package svvcontainer

import (
	"context"
	"os"
	"time"

	"github.com/google/wire"
	gobus "github.com/hetacode/go-bus"
	goeh "github.com/hetacode/go-eh"
	"github.com/hetacode/mechelon/events"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Container struct keep of all dependencies
type Container struct {
	EventsConsumerBus         gobus.ServiceBus
	ProjectsMongoDBCollection *mongo.Collection
}

// NewContainer instance
func NewContainer() *Container {
	wire.Build(
		events.NewEventsMapper,
		initEventsConsumerBusProvider,
		initMongoDBClientProvider,
		wire.Struct(new(Container), "*"),
	)
	return nil
}

func initMongoDBClientProvider() *mongo.Collection {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*15)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_SERVER")))
	if err != nil {
		panic(err)
	}
	collection := client.Database(os.Getenv("SVC_SERVICES_VIEW_DB_NAME")).Collection(os.Getenv("SVC_SERVICES_VIEW_DB_PROJECTS_COLLECTION"))
	return collection
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
