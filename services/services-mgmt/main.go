package main

import (
	"os"

	gobus "github.com/hetacode/go-bus"
	goeh "github.com/hetacode/go-eh"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
	smgcommandhandlers "github.com/hetacode/mechelon/services/services-mgmt/command-handlers"
	smgcontainer "github.com/hetacode/mechelon/services/services-mgmt/container"
)

func main() {
	println("service management is starting")
	waitCh := make(<-chan os.Signal)

	c := smgcontainer.NewContainer()
	eventsMgr := goeh.NewEventsHandlerManager()
	registerEventHandlers(eventsMgr, c)

	go initEventsConsumer(c.CommandsConsumerBus, eventsMgr)

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

func registerEventHandlers(mgr *goeh.EventsHandlerManager, c *smgcontainer.Container) {
	mgr.Register(new(eventsservicesmgmt.RegisterServiceCommand), &smgcommandhandlers.RegisterServiceCommandHandler{Container: c})
	mgr.Register(new(eventsservicesmgmt.UnregisterServiceCommand), &smgcommandhandlers.UnregisterServiceCommandHandler{Container: c})
	mgr.Register(new(eventsservicesmgmt.RemoveServiceInstanceCommand), &smgcommandhandlers.RemoveServiceInstanceCommandHandler{Container: c})
	mgr.Register(new(eventsservicesmgmt.HealthCheckCommand), &smgcommandhandlers.HealthCheckCommandHandler{Container: c})
}
