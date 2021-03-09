package svveventhandlers

import (
	"context"
	"log"

	vssdb "github.com/hetacode/mechelon/services/services-view/db"

	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
	svvcontainer "github.com/hetacode/mechelon/services/services-view/container"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	goeh "github.com/hetacode/go-eh"
)

// ProjectServiceCreatedEventHandler struct
type ProjectServiceCreatedEventHandler struct {
	Container *svvcontainer.Container
}

// Handle ProjectServiceCreatedEvent
func (e *ProjectServiceCreatedEventHandler) Handle(event goeh.Event) {
	log.Printf("ProjectServiceCreatedEventHandler start")
	ev := event.(*eventsservicesmgmt.ProjectServiceCreatedEvent)

	// Check if project and service exists
	filter := bson.M{
		"name": ev.ProjectName,
	}

	var project *vssdb.ProjectEntity
	res := e.Container.ProjectsMongoDBCollection.FindOne(context.Background(), filter)
	if res.Err() == mongo.ErrNoDocuments {
		project = &vssdb.ProjectEntity{
			Name: ev.ProjectName,
			Services: []vssdb.ServiceEntity{
				{Name: ev.ServiceName},
			},
		}

		_, err := e.Container.ProjectsMongoDBCollection.InsertOne(context.Background(), project)

		if err != nil {
			log.Printf("\033[31mProjectServiceCreatedEventHandler insert project err: %s\033[0m", err)
			return
		}
		log.Printf("ProjectServiceCreatedEventHandler project '%s' with service '%s' has been added", ev.ProjectName, ev.ServiceName)
		return
	} else if res.Err() != nil {
		log.Printf("\033[31mProjectServiceCreatedEventHandler find project err: %s\033[0m", res.Err())
		return
	}

	err := res.Decode(&project)
	if err != nil {
		log.Printf("\033[31mProjectServiceCreatedEventHandler decode project entity err: %s\033[0m", err)
		return
	}
	project.Services = append(project.Services, vssdb.ServiceEntity{Name: ev.ServiceName})

	_, err = e.Container.ProjectsMongoDBCollection.UpdateOne(context.Background(), filter, bson.M{"$set": &project})
	if err != nil {
		log.Printf("\033[31mProjectServiceCreatedEventHandler update project entity err: %s\033[0m", err)
		return
	}

	log.Printf("ProjectServiceCreatedEventHandler end")
}
