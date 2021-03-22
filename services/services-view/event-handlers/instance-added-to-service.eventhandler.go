package svveventhandlers

import (
	"context"
	"log"
	"time"

	goeh "github.com/hetacode/go-eh"
	eventsservicesmgmt "github.com/hetacode/mechelon/events/services-mgmt"
	smgeventstore "github.com/hetacode/mechelon/services/services-mgmt/eventstore"
	svvcontainer "github.com/hetacode/mechelon/services/services-view/container"
	vssdb "github.com/hetacode/mechelon/services/services-view/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// InstanceAddedToServiceEventHandler struct
type InstanceAddedToServiceEventHandler struct {
	Container *svvcontainer.Container
}

// Handle ProjectServiceCreatedEvent
func (e *InstanceAddedToServiceEventHandler) Handle(event goeh.Event) {
	log.Printf("InstanceAddedToServiceEventHandler start")
	ev := event.(*eventsservicesmgmt.InstanceAddedToServiceEvent)
	// Check if project and service exists
	filter := bson.M{
		"project_name": ev.ProjectName,
		"service_name": ev.ServiceName,
	}

	var project *vssdb.ServiceEntity
	res := e.Container.ProjectsMongoDBCollection.FindOne(context.Background(), filter)
	if res.Err() == mongo.ErrNoDocuments {
		log.Printf("\033[31mInstanceAddedToServiceEventHandler cannot find project '%s' with service '%s'\033[0m", ev.ProjectName, ev.ServiceName)
		return
	} else if res.Err() != nil {
		log.Printf("\033[31mInstanceAddedToServiceEventHandler find project '%s' with service '%s' err: %s\033[0m", ev.ProjectName, ev.ServiceName, res.Err())
		return
	}

	err := res.Decode(&project)
	if err != nil {
		log.Printf("\033[31mInstanceAddedToServiceEventHandler decode project entity err: %s\033[0m", err)
		return
	}

	for _, i := range project.Instances {
		if i.Name == ev.InstanceName {
			log.Printf("InstanceAddedToServiceEventHandler instance '%s' of service '%s' for project '%s' exists", ev.InstanceName, ev.ServiceName, ev.ProjectName)
			return
		}
	}

	now := time.Now().Unix()
	project.Instances = append(project.Instances, vssdb.InstanceEntity{
		Name:      ev.InstanceName,
		CreatedAt: now,
		UpdatedAt: now,
		Status:    string(smgeventstore.Active),
	})

	_, err = e.Container.ProjectsMongoDBCollection.UpdateOne(context.Background(), filter, bson.M{"$set": &project})
	if err != nil {
		log.Printf("\033[31mProjectServiceCreatedEventHandler instance '%s' for service '%s' in project '%s' update err: %s\033[0m", ev.InstanceName, ev.ServiceName, ev.ProjectName, err)
		return
	}

	log.Printf("InstanceAddedToServiceEventHandler end")
}
