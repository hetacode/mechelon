package main

import (
	"log"
	"os"

	gobus "github.com/hetacode/go-bus"
	goeh "github.com/hetacode/go-eh"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
	smgcommandhandlers "github.com/hetacode/mechelon/services/services-mgmt/command-handlers"
	smgcontainer "github.com/hetacode/mechelon/services/services-mgmt/container"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load("../../.env.dev")

	log.Println("services management svc is starting")
	waitCh := make(<-chan os.Signal)

	c := smgcontainer.NewContainer()
	eventsMgr := goeh.NewEventsHandlerManager()
	registerEventHandlers(eventsMgr, c)

	go initEventsConsumer(c.CommandsConsumerBus, eventsMgr)

	log.Println("services management is running")
	log.Printf("service bus is waiting for commands on the queue: \033[32m%s\033[0m", os.Getenv("SVC_SERVICES_MGMT_SB_COMMANDS_QUEUE"))
	log.Printf("events are produced into the service bus exchange: \033[32m%s\033[0m", os.Getenv("SVC_SERVICES_MGMT_SB_EVENTS_EXCHANGE"))
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

func registerEventHandlers(mgr *goeh.EventsHandlerManager, c *smgcontainer.Container) {
	mgr.Register(new(eventsservicesmgmt.RegisterServiceCommand), &smgcommandhandlers.RegisterServiceCommandHandler{Container: c})
	mgr.Register(new(eventsservicesmgmt.UnregisterServiceCommand), &smgcommandhandlers.UnregisterServiceCommandHandler{Container: c})
	mgr.Register(new(eventsservicesmgmt.RemoveServiceInstanceCommand), &smgcommandhandlers.RemoveServiceInstanceCommandHandler{Container: c})
	mgr.Register(new(eventsservicesmgmt.HealthCheckCommand), &smgcommandhandlers.HealthCheckCommandHandler{Container: c})
}
