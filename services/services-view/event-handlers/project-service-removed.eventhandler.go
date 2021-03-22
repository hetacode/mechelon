package svveventhandlers

import (
	"context"
	"log"

	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
	svvcontainer "github.com/hetacode/mechelon/services/services-view/container"
	"go.mongodb.org/mongo-driver/bson"

	goeh "github.com/hetacode/go-eh"
)

// ProjectServiceRemovedEventHandler struct
type ProjectServiceRemovedEventHandler struct {
	Container *svvcontainer.Container
}

// Handle ProjectServiceRemovedEvent
func (e *ProjectServiceRemovedEventHandler) Handle(event goeh.Event) {
	log.Printf("ProjectServiceRemovedEventHandler start")
	ev := event.(*eventsservicesmgmt.ProjectServiceRemovedEvent)

	// Check if project and service exists
	filter := bson.M{
		"project_name": ev.ProjectName,
		"service_name": ev.ServiceName,
	}

	res, err := e.Container.ProjectsMongoDBCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Printf("\033[31mProjectServiceRemovedEventHandlerfind service project err: %s\033[0m", err)
		return
	}

	if res.DeletedCount == 0 {
		log.Print("ProjectServiceRemovedEventHandler nothing was removed")
	}

	log.Printf("ProjectServiceRemovedEventHandler end")
}
