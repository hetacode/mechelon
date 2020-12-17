package main

import (
	"os"

	gobus "github.com/hetacode/go-bus"
	goeh "github.com/hetacode/go-eh"
	"github.com/hetacode/mechelon/events"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
	smgeventhandlers "github.com/hetacode/mechelon/services/services-mgmt/event-handlers"
)

func main() {
	println("service management is starting")
	waitCh := make(<-chan os.Signal)

	eventsMgr := goeh.NewEventsHandlerManager()
	registerEventHandlers(eventsMgr)

	kind := gobus.RabbitMQServiceBusOptionsFanOutKind
	bus := gobus.NewRabbitMQServiceBus(events.NewEventsMapper(), &gobus.RabbitMQServiceBusOptions{
		Kind:      &kind,
		Exchanage: "test-ex", //TODO: take from env vars
		Queue:     "test-queue",
		Server:    "amqp://rabbit:5672",
	})
	go initEventsConsumer(bus, eventsMgr)

	println("service management is running")
	<-waitCh
}

func initEventsConsumer(bus gobus.ServiceBus, mgr *goeh.EventsHandlerManager) {
	msgCh, errCh := bus.Consume()
	for {
		select {
		case msg := <-msgCh:
			mgr.Execute(msg)
		case err := <-errCh:
			if err != nil {
				panic(err)
			}
		}
	}
}

func registerEventHandlers(mgr *goeh.EventsHandlerManager) {
	mgr.Register(new(eventsservicesmgmt.RegisterServiceEvent), &smgeventhandlers.RegisterServiceEventHandler{})
	mgr.Register(new(eventsservicesmgmt.UnregisterServiceEvent), &smgeventhandlers.UnregisterServiceEventHandler{})
}
