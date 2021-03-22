package main

import (
	"log"
	"os"

	svveventhandlers "github.com/hetacode/mechelon/services/services-view/event-handlers"

	gobus "github.com/hetacode/go-bus"
	goeh "github.com/hetacode/go-eh"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
	svvcontainer "github.com/hetacode/mechelon/services/services-view/container"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../../.env.dev")

	log.Println("services view svc is starting")
	waitCh := make(<-chan os.Signal)

	container := svvcontainer.NewContainer()
	eventsMgr := goeh.NewEventsHandlerManager()
	registerEventHandlers(container, eventsMgr)
	go initEventsConsumer(container.EventsConsumerBus, eventsMgr)

	log.Println("services view is running")
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
				log.Fatalf("Events consumer err: %s", err)
			}
		}
	}
}

func registerEventHandlers(c *svvcontainer.Container, mgr *goeh.EventsHandlerManager) {
	mgr.Register(new(eventsservicesmgmt.ProjectServiceCreatedEvent), &svveventhandlers.ProjectServiceCreatedEventHandler{Container: c})
}
